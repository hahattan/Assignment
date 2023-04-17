ARG BASE=golang:1.20-alpine3.17
FROM ${BASE} AS builder

WORKDIR /

RUN apk add --update --no-cache make git

COPY . .

RUN make client

FROM alpine:3.17

RUN apk add --update --no-cache

WORKDIR /

COPY --from=builder /client/client /client

ENTRYPOINT ["/client"]
