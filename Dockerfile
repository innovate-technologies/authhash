FROM golang:1.11 as build

COPY ./ /go/src/github.com/innovate-technologies/authhash
WORKDIR /go/src/github.com/innovate-technologies/authhash

RUN go build ./

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/innovate-technologies/authhash/authhash /opt/authhash/authhash
COPY --from=build /go/src/github.com/innovate-technologies/authhash/web /opt/authhash/web

WORKDIR /opt/authhash/

CMD ["/opt/authhash"]