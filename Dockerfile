FROM iron/go:dev

RUN adduser -D -h /go go
USER go

ENV GOPATH=/go
ENV PATH=$PATH:/go/bin

RUN go get -u github.com/gophergala2016/gocompatible
RUN rm -rf /go/src/github.com/gophergala2016/gocompatible

CMD gocompatible
