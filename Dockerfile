FROM golang:1.9

WORKDIR /go/src/webcrawler

COPY . .
COPY github.com $GOPATH/src/github.com
COPY golang.org $GOPATH/src/golang.org

RUN echo "Asia/shanghai" > /etc/timezone;
RUN  go build main.go 

CMD ["./main"]

