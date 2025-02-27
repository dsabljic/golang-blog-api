# golang-blog-api

## Getting started

### Install dependencies

```bash
go mod tidy
```

### Prepare the environment (db and redis)

```bash
docker compose up --build -d
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
