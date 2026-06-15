FROM golang:1.25.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker ./cmd/worker/main.go

FROM scratch AS api
COPY --from=builder /api /api
EXPOSE 8080
CMD ["/api"]

FROM scratch AS worker
COPY --from=builder /worker /worker
EXPOSE 8081
CMD ["/worker"]
