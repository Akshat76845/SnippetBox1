FROM golang:1.22
WORKDIR /src/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o snippetbox ./cmd/web
EXPOSE 4000

CMD ["./snippetbox"]

