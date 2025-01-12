FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /build

COPY ./go.mod .
RUN go mod download && go mod verify

COPY . .
RUN go build -o /bin/server cmd/server/main.go

FROM alpine:3.21

WORKDIR /app
COPY --from=builder /bin/server /bin/server

EXPOSE 6000

CMD [ "/bin/server" ]
