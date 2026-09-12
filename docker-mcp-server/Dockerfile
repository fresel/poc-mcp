FROM golang:1.27-alpine AS builder

WORKDIR /app/go-mcp-server

COPY go-mcp-server/go.mod go-mcp-server/go.sum ./
RUN go mod download

COPY go-mcp-server/ ./

RUN mkdir -p /app/test \
    && CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /app/bin/go-mcp-server .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app/go-mcp-server

COPY --from=builder --chown=65532:65532 /app/bin/go-mcp-server /app/bin/go-mcp-server
COPY --from=builder --chown=65532:65532 /app/test /app/test

ENTRYPOINT ["/app/bin/go-mcp-server"]