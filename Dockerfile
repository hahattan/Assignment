FROM golang:1.18-alpine AS builder

WORKDIR /assignment

RUN apk add --update --no-cache make git

COPY . .

RUN make build

FROM alpine:3.14

WORKDIR /
COPY --from=builder /assignment/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]