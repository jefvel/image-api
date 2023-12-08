FROM golang:1.18

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -o /image-api ./cmd/image-api/main.go

EXPOSE 8000

CMD ["/image-api"]