FROM golang:1.9.1-alpine as builder

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg

RUN mkdir -p /go/src/github.com/obieq/docker-hub-webhook
ADD . /go/src/github.com/obieq/docker-hub-webhook/

WORKDIR /go/src/github.com/obieq/docker-hub-webhook

RUN go build -o docker-hub-webhook .

# EXPOSE 4000
# CMD /go/src/github.com/obieq/docker-hub-webhook/docker-hub-webhook

FROM docker
WORKDIR /go/src/github.com/obieq/docker-hub-webhook
COPY --from=builder /go/src/github.com/obieq/docker-hub-webhook/ .

EXPOSE 4000
CMD ./docker-hub-webhook