FROM golang:1.17-alpine

WORKDIR /home
RUN apk add ffmpeg

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o crunch

CMD ["./crunch"]
