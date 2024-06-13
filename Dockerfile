# Multi-stage build to avoid having build/install artifacts in the final runtime Docker image.

# ---- Frontend build ----
FROM node:20 as frontend-build

WORKDIR /app
COPY frontend/ .
RUN npm install && npm run build

# ---- Build ----
FROM golang:1.22 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
COPY --from=frontend-build /app/dist /app/frontend/dist

RUN CGO_ENABLED=1 go build -o react-app

# ---- Runtime ----
FROM golang:1.22 as runtime

WORKDIR /app

COPY --from=build /app/react-app /react-app

# expose default port by default
EXPOSE 8080

CMD ["/react-app"]
