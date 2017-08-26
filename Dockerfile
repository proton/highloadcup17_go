FROM golang:1.9.0-stretch

RUN go get "github.com/valyala/fasthttp"
RUN go get "github.com/pquerna/ffjson"
RUN go get "github.com/buger/jsonparser"
ENV GOMAXPROCS=4
ADD src/*.go ./
RUN go build -o app *.go
EXPOSE 80
CMD ./app
