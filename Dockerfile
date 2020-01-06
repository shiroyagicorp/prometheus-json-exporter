FROM golang:1.13 as builder

COPY . /src

WORKDIR /src

RUN go build -ldflags "-linkmode external -extldflags -static"


FROM alpine:3.11

EXPOSE 9116

ENV USER prometheus

RUN addgroup ${USER} && adduser -D -G ${USER} -h /${USER} ${USER} && \
    apk upgrade --no-cache && \
    apk add --no-cache ca-certificates

COPY --from=builder /src/prometheus-json-exporter /bin

USER ${USER}
CMD ["/bin/prometheus-json-exporter"]
