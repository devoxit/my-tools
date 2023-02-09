FROM golang:1.19
WORKDIR /usr/src/app
COPY go.mod ./
# COPY go.sum ./
# RUN go mod download
RUN go version
RUN pwd
RUN pwd
RUN go build
EXPOSE 6655
CMD ["/usr/src/app/rserver", "tcp", "6655" ]
