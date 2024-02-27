# stage build
FROM golang:1.22.0 AS build

WORKDIR /shylockgo

COPY . /shylockgo

RUN go get github.com/gin-contrib/cors

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o shylockgo main.go

# final image
FROM alpine:latest

WORKDIR /shylockgo

COPY --from=build /shylockgo ./

CMD [ "./shylockgo" ]
