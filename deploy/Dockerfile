FROM golang:1.23
WORKDIR /app


COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o -ldflags="-s -w" server page-analyzer

CMD ["./server"]