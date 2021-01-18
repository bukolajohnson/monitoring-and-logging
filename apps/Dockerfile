FROM golang:1.12-stretch AS builder
WORKDIR $GOPATH/pinger
COPY . .
RUN go build -o app .

FROM ubuntu:bionic

COPY --from=builder /go/pinger/app /
ENV PORT=8000
EXPOSE $PORT

CMD [ "/app" ]
