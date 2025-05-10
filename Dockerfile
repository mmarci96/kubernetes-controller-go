FROM node:20-slim AS app-build

WORKDIR /app

COPY apps/client/ ./

RUN npm install

RUN npm run build

FROM golang:1.24.3-alpine AS build-proxy
WORKDIR /build
RUN apk add --no-cache git
COPY ./proxy/go.mod ./
RUN go mod download
COPY ./proxy/ .
COPY --from=app-build /app/dist ./static/game
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM alpine:3.20 AS proxy
WORKDIR /app
RUN apk add --no-cache dumb-init
COPY --from=build-proxy /build/main /app/main
COPY --from=build-proxy /build/static /app/static
EXPOSE 80
ENTRYPOINT ["/app/main"]


