# docker build --build-arg MULTICAM_RELEASE=6.18.4.4961 --build-arg EXAMPLE=blink -t multicam:latest .
FROM ubuntu:20.04 AS multicam-base

ARG MULTICAM_RELEASE=6.18.4.4961

RUN apt-get update && apt-get install -y apt-file file make gcc linux-headers-5.15.0-46-generic wget
COPY multicam-linux/multicam-linux-x86_64-${MULTICAM_RELEASE}.tar.gz /
RUN tar xf multicam-linux-x86_64-${MULTICAM_RELEASE}.tar.gz
WORKDIR /multicam-linux-x86_64-${MULTICAM_RELEASE}
RUN ./install.sh -m ./drivers/lib64/libMultiCam.so


FROM multicam-base AS multicam-golang

ARG GO_RELEASE=1.18.3

RUN wget https://dl.google.com/go/go${GO_RELEASE}.linux-amd64.tar.gz && \
    tar xf go${GO_RELEASE}.linux-amd64.tar.gz -C /usr/local && \
    rm go${GO_RELEASE}.linux-amd64.tar.gz
ENV PATH=${PATH}:/usr/local/go/bin
RUN go version


FROM multicam-golang AS multicam-app

ARG EXAMPLE=basic

COPY . /src/github.com/northvolt/go-multicam
WORKDIR /src/github.com/northvolt/go-multicam
RUN go build -o /mcexample ./examples/${EXAMPLE}/

ENTRYPOINT ["/mcexample"]
