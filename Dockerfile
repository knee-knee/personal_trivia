FROM golang:latest
RUN mkdir /app
ADD . /app/
ENV AWS_ACCESS_KEY_ID=*************************
ENV AWS_SECRET_ACCESS_KEY=*****************************
WORKDIR /app
EXPOSE 8080
RUN go build -o main .
CMD ["/app/main"]