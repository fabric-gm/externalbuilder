FROM golang:1.15.2-alpine3.12 as builder
WORKDIR /workspace
ENV GOPROXY https://goproxy.cn
COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd/ cmd/
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o launcher ./cmd/launcher


FROM hyperledger/fabric-peer:2.5.0
RUN mkdir -p /builders/golang/bin

COPY --from=builder /workspace/launcher  /builders/golang/bin/externalcc
RUN ln -s /builders/golang/bin/externalcc /builders/golang/bin/detect
RUN ln -s /builders/golang/bin/externalcc /builders/golang/bin/build
RUN ln -s /builders/golang/bin/externalcc /builders/golang/bin/release
RUN ln -s /builders/golang/bin/externalcc /builders/golang/bin/run

COPY images/fabric-peer/2.5.0/k8scc.yaml /builders/golang/bin/k8scc.yaml
