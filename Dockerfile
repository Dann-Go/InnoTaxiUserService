FROM golang:alpine AS builder

WORKDIR /app

ADD . ./
RUN go mod download
RUN go mod tidy
RUN go mod vendor

RUN go build -o /app/userservice ./cmd


FROM alpine:latest
RUN addgroup -g 1000 app
RUN adduser -u 1000 -G app -h /home/goapp -D goapp
USER goapp
WORKDIR /app
COPY --from=builder /app/userservice  /app/
COPY --from=builder /app/internal/migrations/  /app/internal/migrations/

CMD ["./userservice"]