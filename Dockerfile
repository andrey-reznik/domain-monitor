FROM oven/bun:alpine as frontend
WORKDIR /src
COPY . .
RUN bun install
RUN bun run build
RUN ls

FROM docker.io/golang:1.24.1-alpine as go
WORKDIR /src
RUN mkdir app
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .
RUN templ generate
RUN go build -o /app/main ./cmd/main.go
RUN ls

FROM docker.io/alpine:latest
COPY --from=frontend /src/assets /app/assets
COPY --from=go /app/main /app/main

WORKDIR /app

# Setup the data directory
RUN mkdir /app/data
VOLUME /app/data

# Command to run the executable
CMD ["./main","--data-dir","/app/data"]
