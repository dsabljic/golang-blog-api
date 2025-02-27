package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostgresPostStore struct {
	db *sql.DB
}

func (s *PostgresPostStore) Create(ctx context.Context, post *Post) error {
	query := `
		insert into posts (content, title, user_id, tags)
		values ($1, $2, $3, $4) returning id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresPostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	query := `
		select id, user_id, title, content, created_at, updated_at, tags, version
		from posts
		where id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostgresPostStore) Delete(ctx context.Context, id int64) error {
	query := `
		delete from posts where id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresPostStore) Update(ctx context.Context, post *Post) error {
	query := `
		update posts set title = $1, content = $2, version = version + 1 
		where id = $3 and version = $4
		returning version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostgresPostStore) GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE 
			f.user_id = $1 	AND
			(p.title ILIKE '%' || $4 || '%' OR p.content LIKE '%' || $4 || '%') AND
			(p.tags @> $5 OR $5 = '{}') AND
			(CAST($6 AS TIMESTAMPTZ) IS NULL OR p.created_at >= $6) AND 
			(CAST($7 AS TIMESTAMPTZ) IS NULL OR p.created_at >= $7)
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var sinceParam interface{}
	if !fq.Since.IsZero() {
		sinceParam = fq.Since
	} else {
		sinceParam = nil
	}

	var untilParam interface{}
	if !fq.Until.IsZero() {
		untilParam = fq.Until
	} else {
		untilParam = nil
	}

	rows, err := s.db.QueryContext(ctx, query, userId, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags), sinceParam, untilParam)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
}
