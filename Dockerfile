FROM golang:1.7 as builder

WORKDIR /go/src/github.com/mustafaakin/lambda-performance-test
RUN go get golang.org/x/crypto/bcrypt

ADD . /go/src/github.com/mustafaakin/lambda-performance-test
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /tmp/mystaticexe github.com/mustafaakin/lambda-performance-test/container

FROM scratch
COPY --from=builder /tmp/mystaticexe /mystaticexe
ENTRYPOINT ["/mystaticexe"]
