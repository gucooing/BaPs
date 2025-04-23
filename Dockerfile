FROM --platform=$BUILDPLATFORM golang:1.23.2-alpine AS builder
LABEL authors="gucooing"
RUN apk add --no-cache bash protoc protobuf-dev curl

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
WORKDIR /app
COPY . .

RUN --mount=type=secret,id=excel_url,env=EXCEL_URL \
    wget "$EXCEL_URL" --quiet -O ./pkg/mx/excel.go

RUN --mount=type=secret,id=gdconf_dev,env=GDCONF_DEV \
    wget "$GDCONF_DEV" --quiet -O ./gdconf/game.config.dev.go

RUN cd ./common/server_only && \
    protoc --proto_path=. --go_out=. --go_opt=paths=source_relative *.proto && \
    cd ../../
RUN go build -ldflags="-s -w -X github.com/gucooing/BaPs/pkg/mx.Docker=1" -o /app/BaPs ./cmd/BaPs/BaPs.go

# 最终镜像
FROM alpine:latest
RUN apk add --no-cache bash tzdata
WORKDIR /usr/ba
COPY --from=builder /app/BaPs .
COPY --from=builder /app/data/ ./data/
RUN chmod +x BaPs
EXPOSE 5000/tcp
ENTRYPOINT ["./BaPs -c ./config/config.json"]