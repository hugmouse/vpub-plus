# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app
ADD . /app
RUN go build -o bin/vpub main.go

# Simulate running this thing as a user,
# not really needed though!
FROM alpine:latest

RUN apk --no-cache add ca-certificates
COPY --from=builder /app/sbin/vpub /usr/local/bin/vpub
RUN adduser -D -g '' vpub
USER vpub
EXPOSE 8080
ENTRYPOINT ["vpub"]
