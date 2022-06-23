FROM golang:1.14.7 as builder

ENV GO111MODULE=on
RUN mkdir /build
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN GOOS=linux go build -o xm .
FROM scratch
COPY --from=builder /build/xm /app/
WORKDIR /app

CMD ["./xm"]
