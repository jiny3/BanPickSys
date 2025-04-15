# Use golang LTS version as builder
FROM golang:latest AS build

WORKDIR /app

COPY . .

RUN go env -w GOPROXY='https://goproxy.io,direct'
RUN go mod tidy
RUN go build -o bin/main .

FROM ubuntu:22.04

WORKDIR /app

COPY --from=build /app/bin ./bin
COPY --from=build /app/config ./config
COPY --from=build /app/static ./static

# Expose application port (change if needed)
EXPOSE 10088

# Use proper CMD array syntax
ENV GIN_MODE=release
CMD ["bin/main"]
