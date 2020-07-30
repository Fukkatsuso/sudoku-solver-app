# base
FROM ubuntu:18.04 AS base

RUN apt update && apt -y install \
        gnupg \
        sudo \
        curl \
        git \
        make

ENV GO_VERSION 1.14.6
ENV GOPATH /go 
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "${GOPATH}/src" "${GOPATH}/bin" && chmod -R 777 "${GOPATH}"
    
RUN curl -Lso go.tar.gz "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

RUN go get -u -v -d gocv.io/x/gocv && \
    cd $GOPATH/src/gocv.io/x/gocv && \
    make install

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv 8529b1e0f8bf7f65c12fabb0a4bcbd87cef9e52d && \
    echo "deb http://ppa.launchpad.net/alex-p/tesseract-ocr/ubuntu bionic main" > /etc/apt/sources.list.d/tesseract-ocr.list && \
    echo "deb-src http://ppa.launchpad.net/alex-p/tesseract-ocr/ubuntu bionic main" >> /etc/apt/sources.list.d/tesseract-ocr.list && \
    apt update && apt -y install \
        tesseract-ocr \
        libtesseract-dev \
        libleptonica-dev && \
    rm -rf /var/lib/apt/lists/*

RUN go get -v -t github.com/otiai10/gosseract

RUN go get -v -u github.com/Fukkatsuso/sudoku

ENV PORT 8080
EXPOSE 8080

WORKDIR $GOPATH/src/github.com/Fukkatsuso/sudoku-solver-app

# dev
FROM base AS dev

RUN go get -v github.com/oxequa/realize

# release
FROM base AS release

COPY . .

RUN go build -o server main.go

CMD [ "./server" ]