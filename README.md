# Golang+React+Vite Single Binary app

This is an example/demo/boilplate of a Vite + React app served in a single binary using Golang.

The Vite app is in [frontend/](./frontend) directory. You can change it to Next.js or anything **that builds into static files**.  
The Golang code is in [main.go](./main.go), it uses [chi](https://github.com/go-chi/chi) and you can modify it to provide an API for the frontend.

## Development

Run the local dev env with:

```bash
docker compose up
```

This will start a dev server for the React app and the Golang server. The Golang app will proxy requests to the React app.

Open [http://localhost:8080](http://localhost:8080).

### Without Docker

If you don't want to use Docker, you can open two terminals and run the backend & frontend separately:

```bash
cd frontend && npm run dev
```

```bash
ENV=dev VITE_URL=http://localhost:5137 go run ./main.go
```

### Env vars

These are the env vars used by this

```
# Running environment for the Golang backend
ENV=dev

# Host for binding the Golang HTTP server
HOST=

# Port for binding the Golang HTTP server
PORT=8080

# URL to the Vite dev server (used only in dev env, not when building the "prod" binary)
VITE_URL=http://localhost:5137
```

## Understand

This is just a quick demo of an idea to serve a React app from a single binary using Golang. The React app is built separately and the static files are embedded into the binary using the `go:embed` directive when the Golang app binary is built. In dev env, the Golang code in `main.go`-file proxies all requests to the Vite dev server served by `VITE_URL`.

All the important backend/server bits are in the [main.go](./main.go) file.

Frontend is completely isolated into the [frontend/](./frontend) directory.

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

In build, we use `go:embed` to embed the frontend into the binary.

### Build Docker image

Alternativel, you can build the binary into a Docker image:

```bash
docker build -t react-go .
```

And run it with:

```bash
docker run -p 8080:8080 react-go
```
