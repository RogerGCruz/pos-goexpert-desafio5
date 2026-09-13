ARG TARGETOS=linux
ARG TARGETARCH=amd64
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY cmd/ ./cmd
COPY src/ ./src
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /loadtest ./cmd

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /loadtest /loadtest
RUN chmod +x /loadtest
ENTRYPOINT ["/loadtest"]
