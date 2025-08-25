FROM golang:1.21-alpine 
RUN apk add --no-cache git 
WORKDIR /app 
COPY . . 
RUN go mod tidy 
RUN go build -o main main.go 
EXPOSE 8080 
CMD ["./main"] 
