FROM golang:1.19
WORKDIR /usr/src/app
COPY go.mod ./
# COPY go.sum ./
# RUN go mod download
RUN go version
# RUN go build
EXPOSE 6655
CMD [ "go", "run","./rserver.go", "tcp", "6655" ]
