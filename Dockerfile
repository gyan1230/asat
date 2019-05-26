FROM alpine:latest

ARG mongo_addr 

ENV mongo_uri=$mongo_addr

ADD asat  /asat

RUN echo $mongo_uri

CMD ["./asat"]