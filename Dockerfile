FROM ubuntu:16.04

#Begin: install prerequisites
RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential \
        curl \
        git \
        libcurl3-dev \
        libfreetype6-dev \
        libpng12-dev \
        libzmq3-dev \
        locate \
        pkg-config \
        rsync \
        software-properties-common \
        sudo \
        unzip \
        vim \
        wget \
        zip \
        zlib1g-dev \
        && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
#End: install prerequisites

#Begin: install golang
ENV GOLANG_VERSION 1.11
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_SHA256_CHECKSUM b3fcf280ff86558e0559e185b601c9eade0fd24c900b4c63cd14d1d38613e499
ENV GOPATH $HOME/go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz && \
    echo "$GOLANG_SHA256_CHECKSUM golang.tar.gz" | sha256sum -c - && \
    sudo tar -C /usr/local -xzf golang.tar.gz && \
    rm golang.tar.gz && \
    mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
#End: install golang

#Begin: install delve
RUN go get github.com/derekparker/delve/cmd/dlv
#End: install delve

#Begin: install dep
RUN apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
    wget -O dep.zip https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64.zip && \
    echo '287b08291e14f1fae8ba44374b26a2b12eb941af3497ed0ca649253e21ba2f83 dep.zip' | sha256sum -c - && \
    unzip -d /usr/bin dep.zip && rm dep.zip
#End: install dep

#Begin: install github-teamwork
RUN go get -u github.com/ctava/github-teamwork
ADD .env /go/src/github.com/ctava/github-teamwork
WORKDIR /go/src/github.com/ctava/github-teamwork
RUN dep ensure
#End: install github-teamwork