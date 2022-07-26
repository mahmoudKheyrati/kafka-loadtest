FROM golang:1.15-buster AS build

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod tidy && go mod vendor

COPY ./ ./

RUN go build -o kafka-publish-load-test
RUN ls

#FROM gcr.io/distroless/base-debian10

#WORKDIR /

#RUN apt-get -y install librdkafka-dev
#COPY --from=build /app/kafka-publish-load-test /app/kafka-publish-load-test

ENTRYPOINT ["/app/kafka-publish-load-test", "-server-name", "kafkaLoadTestKuber","-bootstrap-servers", "192.168.12.122:9093, 192.168.12.123:9093, 192.168.12.124:9093", "-topic", "kafka-loadtest", "-batch-count", "100", "-batch-size", "100000"]
