FROM golang

RUN go get "github.com/valyala/fasthttp"
RUN go get "github.com/pquerna/ffjson"
ENV GOMAXPROCS=4
ADD src/*.go ./
RUN go build -o app *.go
EXPOSE 80
CMD ./app
