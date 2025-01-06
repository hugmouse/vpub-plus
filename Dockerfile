FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go generate
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /vpub main.go

FROM scratch

COPY --from=builder /vpub /vpub
EXPOSE 8080
ENTRYPOINT ["/vpub"]