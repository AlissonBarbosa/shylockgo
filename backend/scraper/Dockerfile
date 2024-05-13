# stage build
FROM golang:1.22.2 AS build

WORKDIR /shylockgo

COPY . /shylockgo

RUN go get github.com/gin-contrib/cors

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o shylockgo main.go

# final image
FROM alpine:latest

RUN apk add --no-cache tzdata
ENV TZ=America/Recife

WORKDIR /shylockgo

COPY --from=build /shylockgo ./

CMD [ "./shylockgo" ]
