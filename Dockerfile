FROM golang:1.26rc2-alpine3.23 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go generate
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /vpub main.go

FROM alpine:20260127

COPY --from=builder /vpub /vpub
EXPOSE 8080
ENTRYPOINT ["/vpub"]