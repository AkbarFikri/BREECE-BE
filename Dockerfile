FROM golang:alphine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .
RUN go build -o build

FROM alphine:latest AS release

COPY --from=builder /app/build /

ENTRYPOINT [ "/build" ]