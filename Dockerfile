FROM golang:1.10

RUN go get -u github.com/golang/dep/...

ARG PACKAGE_NAME=github.com/shiroyagicorp/prometheus-json-exporter

COPY . /go/src/$PACKAGE_NAME
RUN cd /go/src/$PACKAGE_NAME && dep ensure -vendor-only
RUN go install $PACKAGE_NAME

FROM alpine:latest  
COPY --from=0 /go/bin/prometheus-json-exporter .
CMD ["./prometheus-json-exporter"]
EXPOSE 9116
