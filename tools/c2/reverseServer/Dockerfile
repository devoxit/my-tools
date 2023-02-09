FROM golang:1.19
WORKDIR /usr/src/app
COPY go.mod ./
# COPY go.sum ./
# RUN go mod download
RUN go build -o rserver
EXPOSE 6655
CMD [ "/rserver", "tcp", "6655" ]
