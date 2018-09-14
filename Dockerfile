FROM golang:1.10 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go get -u github.com/golang/dep/...

ARG PACKAGE_NAME=github.com/shiroyagicorp/prometheus-json-exporter

COPY . /go/src/$PACKAGE_NAME
RUN cd /go/src/$PACKAGE_NAME && dep ensure -vendor-only
RUN go install $PACKAGE_NAME

FROM alpine:latest  
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/prometheus-json-exporter .
CMD ["./prometheus-json-exporter"]
EXPOSE 9116
