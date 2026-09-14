ARG TARGETOS=linux
ARG TARGETARCH=amd64
FROM golang:1.20-alpine AS builder
WORKDIR /app/src
# copy module file from src (go.sum may be absent) and the source
COPY src/go.mod ./
COPY src/ ./

# build from the src module
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /loadtest ./cmd

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /loadtest /loadtest
RUN chmod +x /loadtest
ENTRYPOINT ["/loadtest"]
