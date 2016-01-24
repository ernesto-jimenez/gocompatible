FROM iron/go:dev

RUN adduser -D -h /go go

ENV GOPATH=/go
ENV PATH=$PATH:/go/bin

ADD . /go/src/github.com/gophergala2016/gocompatible
RUN go install github.com/gophergala2016/gocompatible
RUN rm -rf /go/src

USER go
CMD gocompatible
