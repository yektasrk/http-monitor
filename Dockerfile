FROM golang:1.17 as builder
ENV GO111MODULE=on

WORKDIR /opt
COPY Makefile go.mod go.sum ./
RUN make dep
COPY . .
RUN make build

FROM ubuntu:18.04
ENV TZ=Asia/Tehran
COPY --from=builder /opt/monitord /http-monitor/monitord
ADD ./configs/configs.yaml /configs/
ENTRYPOINT ["/http-monitor/monitord", "--config", "./configs/configs.yaml"]
