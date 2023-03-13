FROM golang:1.16 as builder
COPY ./ /usr/local/go/src/github.com/Rbsn-joses/kafka-replicator
WORKDIR /usr/local/go/src/github.com/Rbsn-joses/kafka-replicator
ENV GO111MODULE=on

RUN go mod verify && go mod tidy && go build -o /entrypoint main.go 

FROM debian:latest
ENV TZ=America/Sao_Paulo
RUN apt-get update && apt-get -y install -y curl iputils-ping net-tools nano telnet
ARG USERNAME=kafka-replicator
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Create the user
RUN groupadd --gid $USER_GID $USERNAME  && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

# ********************************************************
# * Anything else you want to do like clean up goes here *
# ********************************************************

# [Optional] Set the default user. Omit if you want to keep the default as root.
USER $USERNAME


ENV WORKER_COUNT 8
WORKDIR /tmp/workspace/clusters
COPY --from=builder /entrypoint /
ENTRYPOINT [ "/entrypoint" ]
