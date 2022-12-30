FROM golang:1.19.3-alpine
RUN mkdir customer
COPY . /customer
WORKDIR /customer
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 3000