
FROM golang:latest
COPY . /app
WORKDIR /app
RUN go build -o main .
EXPOSE 3030
CMD ["./main"]
