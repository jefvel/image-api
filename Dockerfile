FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -o /image-api ./cmd/main.go

EXPOSE 8000

CMD ["/image-api"]