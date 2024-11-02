FROM docker.io/golang:latest AS build-env
COPY . /build
WORKDIR /build
RUN \
  go mod init github.com/rwxrob/ssh &&\
  go work init &&\
  go work use . &&\
  go mod tidy &&\
  go build ./cmd/runonany

FROM ubuntu:latest
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=America/New_York
RUN apt-get update && apt-get install -y openssh-client openssh-server vim
COPY --from=build-env /build/runonany /usr/bin/runonany
ADD testdata/runonany.yaml /etc
RUN groupadd -r user && useradd -r -g user user
RUN mkdir -p /home/user/.ssh
COPY testdata/keys/* /home/user/.ssh
RUN chmod 600 /home/user/.ssh/user
RUN chown -R user:user /home/user
COPY entrypoint /
ENTRYPOINT ["/entrypoint"]

