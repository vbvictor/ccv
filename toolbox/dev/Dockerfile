FROM debian:bullseye-slim

ENV DEBIAN_FRONTEND=noninteractive
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

RUN apt-get update && apt-get install -y \
    gcc \
    git \
    curl \
    wget \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://ftp.gnu.org/gnu/make/make-4.4.1.tar.gz && \
    tar xf make-4.4.1.tar.gz && \
    cd make-4.4.1 && \
    ./configure && \
    make && \
    make install && \
    cd .. && \
    rm -rf make-4.4.1 make-4.4.1.tar.gz

RUN curl -fsSL https://go.dev/dl/go1.23.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0

WORKDIR /app

VOLUME ["/app", "/go/pkg/mod"]

CMD ["bash"]
