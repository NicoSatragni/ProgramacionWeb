FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/tp2 .

FROM alpine:3.24
RUN adduser -D -u 1000 app
USER app

COPY --from=builder /app/tp2 /usr/local/bin/tp2
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["tp2"]


