FROM golang:1.11

WORKDIR /go/src/app

RUN go get -d -v github.com/gin-gonic/gin github.com/jinzhu/gorm github.com/jinzhu/gorm/dialects/sqlite gopkg.in/mailgun/mailgun-go.v1 golang.org/x/crypto/bcrypt github.com/satori/go.uuid

COPY . .

# RUN go install -v ./

EXPOSE 8080

RUN go build server.go adminFunctions.go userFunctions.go models.go generalFunctions.go

CMD ["./server"]
