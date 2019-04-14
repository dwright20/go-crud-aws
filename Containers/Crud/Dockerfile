FROM golang:1.8

#copy local import
COPY ./game /go/src/web_project/ecs_crud/game

WORKDIR /go/src/app
COPY . .

RUN go get -d github.com/gorilla/mux
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]