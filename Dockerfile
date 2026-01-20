FROM golang:alpine

WORKDIR /

COPY node-order-scheduler  /usr/local/bin

CMD ["node-order-scheduler"]