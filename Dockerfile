FROM --platform=$BUILDPLATFORM golang:1.23.2-alpine AS builder
LABEL authors="gucooing"
RUN apk add --no-cache bash protoc protobuf-dev curl python3 py3-pip \
    gcc musl-dev libffi-dev openssl-dev

RUN python3 -m pip install --upgrade pip && \
    pip install awscli \
    
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
WORKDIR /app
COPY . .
ARG R2_ENDPOINT
ARG R2_BUCKET
RUN --mount=type=secret,id=R2_ACCESS_KEY_ID \
    --mount=type=secret,id=R2_SECRET_ACCESS_KEY \
    export R2_ACCESS_KEY_ID=$(cat /run/secrets/R2_ACCESS_KEY_ID) && \
    export R2_SECRET_ACCESS_KEY=$(cat /run/secrets/R2_SECRET_ACCESS_KEY) && \
    AWS_ACCESS_KEY_ID=$R2_ACCESS_KEY_ID \
    AWS_SECRET_ACCESS_KEY=$R2_SECRET_ACCESS_KEY \
    aws s3 cp --endpoint-url "$R2_ENDPOINT" "s3://$R2_BUCKET" ./pkg/mx/excel.go
RUN cd ./common/server_only && \
    protoc --proto_path=. --go_out=. --go_opt=paths=source_relative *.proto && \
    cd ../../
RUN go build -ldflags="-s -w" -o /app/BaPs ./cmd/BaPs/BaPs.go

# 最终镜像
FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /usr/ba
COPY --from=builder /app/BaPs .
COPY --from=builder /app/data/ ./data/
RUN chmod +x BaPs
EXPOSE 5000/tcp
ENTRYPOINT ["./BaPs"]