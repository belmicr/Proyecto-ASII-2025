# syntax=docker/dockerfile:1

### Etapa de build
FROM golang:1.22 AS build
WORKDIR /app

# Dependencias primero (cache)
COPY go.mod go.sum ./
RUN go mod download

# Código
COPY . .

# Binario estático
RUN CGO_ENABLED=0 GOOS=linux go build -o hotels-api .

### Imagen final mínima
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/hotels-api /app/hotels-api
EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/app/hotels-api"]
