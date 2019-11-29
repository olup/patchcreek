FROM golang AS builder

WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
RUN go get -d -v

RUN go build -o /go/bin/pubsub

FROM alpine

WORKDIR /go/bin

COPY --from=builder /go/bin/pubsub ./pubsub

CMD ["./pubsub"]