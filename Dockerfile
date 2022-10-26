# docker build --build-arg MULTICAM_RELEASE=6.18.4.4961 --build-arg EXAMPLE=blink -t multicam:latest .
FROM ubuntu:20.04 AS multicam-base

ARG MULTICAM_RELEASE=6.18.4.4961
ENV MULTICAM_SDK_VERSION=$MULTICAM_RELEASE

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
RUN go get golang.org/x/tools/cmd/stringer && go generate .
RUN go test -v
RUN mkdir -p /build && \
    go build -o /build/basic ./examples/basic/ && \
    go build -o /build/blink ./examples/blink/ && \
    go build -o /build/capture ./examples/capture/ && \
    go build  -o /build/metadata ./examples/metadata/ && \
    go build -o /build/multi ./examples/multi/ && \
    go build  -o /build/multisignal ./examples/multisignal/

CMD ["/build/${EXAMPLE}"]
