FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/api ./cmd/api && \
    CGO_ENABLED=0 go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/worker ./cmd/worker

FROM alpine:3
RUN addgroup -S app && adduser -S -G app app && mkdir -p /app/.data/uploads && chown -R app:app /app
WORKDIR /app
COPY --from=build /out/api /out/worker /app/
USER app
ENV SERVER_ADDR=:8081 GIN_MODE=release APP_TIMEZONE=Asia/Ho_Chi_Minh
EXPOSE 8081
CMD ["/app/api"]
