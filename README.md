# golang-blog-api

## Getting started

### Install dependencies

```bash
go mod tidy
```

### Configure Environment Variables

Create a .env file in the root directory and populate it with the following:

```env
ADDR=:8080
DB_ADDR=postgres://admin:adminpassword@localhost:54332/social?sslmode=disable
ENV=development
SENDGRID_API_KEY=...
FROM_EMAIL=[email address from which the confirmation email is sent]
```

### Start services (pg db and redis)

```bash
docker compose up --build -d
```

### Run db migrations

```bash
make migrate-up
```

### Run the server

```bash
air
```

or you can start the server manually by running:

```bash
go build -o ./bin/main ./cmd/api
./bin/main
```

### Setup and run the frontend (just a confirmation page for now)

```bash
cd web
npm install
npm run dev
```
