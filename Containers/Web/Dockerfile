FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go get -d github.com/gorilla/mux
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]