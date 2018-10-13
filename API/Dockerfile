FROM golang:1.11

WORKDIR /go/src/app

COPY . .

RUN go get -d -v github.com/gin-gonic/gin
RUN go get -d -v github.com/jinzhu/gorm
RUN go get -d -v github.com/jinzhu/gorm/dialects/postgres
RUN go get -d -v github.com/jinzhu/gorm/dialects/sqlite

RUN go install -v ./

EXPOSE 8080

CMD ["go", "run", "server.go"]
