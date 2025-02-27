package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int       `json:"limit" validaeta:"gte=1 lte=50"`
	Offset int       `json:"offset" validate:"gte=0"`
	Sort   string    `json:"sort" validate:"oneof=asc desc"`
	Tags   []string  `json:"tags" validate:"max=5"`
	Search string    `json:"search" validate:"max=100"`
	Since  time.Time `json:"since"`
	Until  time.Time `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = l
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	} else {
		fq.Tags = []string{}
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		parsedSince, err := parseTime(since)
		if err == nil {
			fq.Since = *parsedSince
		}
	}

	until := qs.Get("since")
	if until != "" {
		parsedUntil, err := parseTime(until)
		if err == nil {
			fq.Until = *parsedUntil
		}
	}

	return fq, nil
}

func parseTime(s string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, s) // Ensure it's RFC3339
	if err != nil {
		return nil, err
	}
	return &t, nil
}
