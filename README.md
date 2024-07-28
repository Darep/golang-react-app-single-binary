# Single binary React app using Golang

This is an example/demo of a Vite + React app served in a single binary using Golang. The Vite app is in [frontend/](./frontend) directory. You can change it to Next.js or anything **that builds into static files** and use the Golang code provide an API & serve the static files.

## Development

Run the local dev env with:

```bash
docker compose up
```

This will start a dev server for the React app and the Golang server. The Golang app will proxy requests to the React app.

Open [http://localhost:8080](http://localhost:8080).

## Build

To build to a single binary locally:

```bash
cd frontend && npm run build && cd ..
go build -o react-app main.go
```

And then start it:

```bash
PORT=8080 ./react-app
```

Open [http://localhost:8080](http://localhost:8080). You should see the Vite example app.

If you wish to build multiple binaries for different platforms, you can do that with `GOARCH` and `GOOS` env variables. For example, to build for ARM-based Androids:

```bash
GOARCH=arm64 GOOS=linux go build -o react-app-linux-arm64 main.go
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
