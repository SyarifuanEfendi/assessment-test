FROM golang:1.18.4

RUN apt-get-update && apt-get install -y

ENV PKG_NAME=assessment-test/
ENV PKG_PATH=$GOPATH/src/$PKG_NAME
WORKDIR $PKG_PATH

COPY . $PKG_PATH/

RUN echo $PWD
RUN go mod vendor

WORKDIR $PKG_PATH/
RUN echo $PWD

RUN go build server.go
EXPOSE 1323
CMD [ "./server" ]