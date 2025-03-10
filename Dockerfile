FROM --platform=$BUILDPLATFORM golang:1.23.2-alpine AS builder
LABEL authors="gucooing"
RUN apk add --no-cache bash protoc protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
WORKDIR /app
COPY . .
RUN cd ./common/server_only && \
    protoc --proto_path=. --go_out=. --go_opt=paths=source_relative *.proto && \
    cd ../../
RUN go build -ldflags="-s -w" -tags "rel" -o /app/BaPs ./cmd/BaPs/BaPs.go

# 最终镜像
FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /usr/ba
COPY --from=builder /app/BaPs .
COPY --from=builder /app/data/ ./data/
RUN chmod +x BaPs
EXPOSE 5000/tcp
ENTRYPOINT ["./BaPs"]