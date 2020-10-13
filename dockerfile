FROM golang:latest
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
WORKDIR /app
COPY . .
RUN go build
EXPOSE 9090
CMD ["./test"]
