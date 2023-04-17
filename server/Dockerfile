ARG BASE=golang:1.20-alpine3.17
FROM ${BASE} AS builder

WORKDIR /

RUN apk add --update --no-cache make git

COPY . .

RUN make server

FROM alpine:3.17

RUN apk add --update --no-cache dumb-init

WORKDIR /

COPY --from=builder /server/server /server

EXPOSE 50051

ENTRYPOINT ["/server"]