FROM golang:1.14

WORKDIR /src/
COPY .  .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o operator-test-framework


FROM bitnami/kubectl
WORKDIR /target
# 拷贝后端二进制文件
COPY --from=0 /src/operator-test-framework .
ENTRYPOINT [ "/target/operator-test-framework","run","--configPath","/conf"]