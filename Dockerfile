FROM golang:1.18 as builder
COPY ./ /usr/local/go/src/github.com/Rbsn-joses/kafka-replicator
WORKDIR /usr/local/go/src/github.com/Rbsn-joses/kafka-replicator
ENV GO111MODULE=on

RUN go mod verify  && go build -o /entrypoint main.go 

FROM debian:latest
ENV TZ=America/Sao_Paulo
RUN apt-get update && apt-get -y install -y curl iputils-ping net-tools nano telnet
ARG USERNAME=kafka-replicator
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME  && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

USER $USERNAME


ENV WORKER_COUNT 8
WORKDIR /tmp/workspace/clusters
COPY --from=builder /entrypoint /
ENTRYPOINT [ "/entrypoint" ]
