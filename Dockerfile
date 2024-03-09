# Dockerfile

# ======= build stage =======
FROM golang:alpine as stage-build

RUN apk update && apk add --no-cache git

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .
RUN go build -o executable ./cmd/app

# ======= release stage =======
FROM alpine:latest as stage-release

COPY --from=stage-build /executable /
COPY --from=stage-build /.env /
COPY --from=stage-build /api /api

ENTRYPOINT ["./executable"]