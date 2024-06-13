# Single binary React app using Golang

React app served in a single binary using Golang.

## Getting Started

Run the local dev env with:

```bash
docker compose up
```

## Build

To build the binary locally:

```bash
cd frontend && npm run build && cd ..
go build -o react-app main.go
```

And then start it:

```bash
PORT=8080 ./react-app
```

## Build Docker image

To build the Docker image:

```bash
docker build -t react-go .
```

And run it with:

```bash
docker run -p 8080:8080 react-go
```

## Understand

This is just a quick demo of an idea to serve a React app from a single binary using Golang. The React app is built separately and the static files are embedded into the binary using the `go:embed` directive when the Golang app binary is built.

All the important bits are in the [main.go](./main.go) file.
