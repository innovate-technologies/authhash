FROM golang:1.11 as build

COPY ./ /go/src/github.com/innovate-technologies/authhash
WORKDIR /go/src/github.com/innovate-technologies/authhash

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/innovate-technologies/authhash/authhash /opt/authhash/authhash
COPY --from=build /go/src/github.com/innovate-technologies/authhash/web /opt/authhash/web

RUN chmod +x /opt/authhash/authhash

WORKDIR /opt/authhash/

CMD ["/opt/authhash/authhash"]