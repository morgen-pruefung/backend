# Morgen-Pruefung.de Backend

## Tech Stack

- Go

## Build and run

```bash
$ go build -o out/backend ./...
$ ./out/backend
```

## Docker

Build the image:

```bash
docker build -t mp-backend .
```

Run the container:

```bash
docker run -p 4242:4242 mp-backend
```