FROM golang:1.11.4

ARG mongo_addr 

ENV mongo_uri=$mongo_addr

WORKDIR /

COPY . .

RUN go build -o asat main.go


RUN echo $mongo_uri

CMD ["./asat"]
