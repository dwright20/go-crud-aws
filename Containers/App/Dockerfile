FROM golang:1.8

#copy local import
COPY ./hiddenCreds /go/src/web_project/ecs_app/hiddenCreds

WORKDIR /go/src/app
COPY . .

RUN go get -d github.com/gorilla/mux
RUN go get github.com/lib/pq
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]