FROM gocv/opencv:latest AS deps
WORKDIR /go/src/github.com/Fukkatsuso/sudoku-solver-app
RUN apt update \
    && apt install -y \
    tesseract-ocr \
    libtesseract-dev
COPY go.* ./
RUN go mod download
ENV PORT 8080
EXPOSE 8080

FROM deps AS dev
RUN go install github.com/pilu/fresh@latest
CMD [ "fresh" ]

FROM deps AS release
COPY . .
RUN go build -o server main.go
CMD [ "./server" ]
