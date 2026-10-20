FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api

FROM alpine:3.21
RUN adduser -D -H app
USER app
COPY --from=build /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
