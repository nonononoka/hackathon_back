FROM golang:1.22.2 as build
WORKDIR /back
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main /back/main.go
CMD /back/main