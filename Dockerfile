FROM golang:1.24.5-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o gogit main.go

CMD ["./gogit", "init"]