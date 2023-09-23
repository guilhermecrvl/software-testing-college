
FROM golang:latest
COPY . /app
WORKDIR /app
RUN go build -o main .
EXPOSE 8080
# Comando para iniciar a aplicação quando o contêiner for executado
CMD ["./main"]