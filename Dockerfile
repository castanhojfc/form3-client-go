FROM golang:1.20

RUN mkdir /form3-client-go

ADD . /form3-client-go

WORKDIR /form3-client-go
