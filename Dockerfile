FROM golang:1.20

RUN mkdir /client

ADD . /client

WORKDIR /client
