FROM easzlab.io.local:5000/golang:1.15.2-alpine3.12 as builder
WORKDIR /workspace
ENV GOPROXY https://goproxy.cn
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY cmd/ cmd/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o fileserver ./cmd/fileserver/fileserver.go


FROM easzlab.io.local:5000/ubuntu:20.04

COPY --from=builder /workspace/fileserver .

ENTRYPOINT ["/fileserver"]

EXPOSE 8080
