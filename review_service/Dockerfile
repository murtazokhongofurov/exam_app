FROM golang:1.19.3-alpine
RUN mkdir review
COPY . /review
WORKDIR /review
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 5000