FROM golang:1.22.2 as build
WORKDIR /back
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY /Users/murakinonoka/downloads/term5-nonoka-muraki-04f45199993c.json /back/serviceAccountKey.prod.json
RUN go build -o main /back/main.go
CMD /back/main