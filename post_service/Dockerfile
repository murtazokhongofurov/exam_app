FROM golang:1.19.3-alpine
RUN mkdir post
COPY . /post
WORKDIR /post
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 7000