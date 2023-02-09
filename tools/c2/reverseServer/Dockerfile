FROM golang:1.19
WORKDIR /usr/src/app
COPY go.mod ./
# COPY go.sum ./
# RUN go mod download
COPY . .
RUN go version
RUN ls
RUN pwd
RUN go build
EXPOSE 6655
CMD ["/usr/src/app/rserver", "tcp", "6655" ]
