FROM golang:1.24-alpine AS build
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache tzdata
ENV TZ="UTC"
COPY --from=build /app/server /app/
EXPOSE 8080
RUN mkdir -p /app/storage
VOLUME /app/storage
ENV DB_PATH="/app/storage/database/recordgo.db"
ENV STORAGE_ROOT="/app/storage"
ENV PORT="8080"
ENV GIN_MODE="release"
ENV PARENT_COINS_SYNC="true"
ENTRYPOINT ["./server"]