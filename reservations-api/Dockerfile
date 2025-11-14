# syntax=docker/dockerfile:1

FROM golang:1.22 AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/server .

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=build /app/bin/server /app/server
COPY --from=build /app/db /app/db
EXPOSE 8086
ENV GIN_MODE=release
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
