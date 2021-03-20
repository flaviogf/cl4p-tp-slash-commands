FROM golang:1.16.2
WORKDIR /go/src/app
COPY . .
RUN go build .
ENTRYPOINT [ "./cl4p-tp-slash-commands" ]
