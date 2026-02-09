# syntax=docker/dockerfile:1

FROM golang:1.24

WORKDIR /forum
COPY . .
RUN go mod download

RUN go build -o app ./cmd/web

EXPOSE 8080

CMD [ "./app" ]
