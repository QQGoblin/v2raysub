# Create a minimal container to run a Golang static binary
FROM registry.lqingcloud.cn/library/alpine:3.16.2

COPY  ./v2raysub /v2raysub

ENTRYPOINT ["/v2raysub"]
EXPOSE 80