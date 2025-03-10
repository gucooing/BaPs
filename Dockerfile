FROM --platform=$BUILDPLATFORM golang:1.23.2-alpine AS builder
LABEL authors="gucooing"
RUN apk add --no-cache bash protoc protobuf-dev curl

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
WORKDIR /app
COPY . .

#ARG EXCEL_GO_SECRET
#RUN mkdir -p ./pkg/mx && \
#    printf "%s" "$EXCEL_GO_SECRET" | base64 -d > ./pkg/mx/excel.go

RUN --mount=type=secret,id=excel_go,env=EXCEL_GO_SECRET \
    mkdir -p ./pkg/mx && \
    cat /run/secrets/excel_go > ./pkg/mx/excel.go

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