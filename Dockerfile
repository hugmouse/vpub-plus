FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . . 
RUN go generate && go build -o /vpub main.go

FROM scratch

COPY --from=builder /vpub /vpub
EXPOSE 8080
ENTRYPOINT ["/vpub"]