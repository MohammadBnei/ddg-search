FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine AS certimage
RUN apk add --no-cache ca-certificates

FROM scratch

COPY --from=certimage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certimage /usr/share/ca-certificates/mozilla/isrgrootx1.pem /etc/ssl/certs/

WORKDIR /app

COPY --from=builder /app/server .

ENV PORT=8080

EXPOSE 8080

CMD ["./server"]
