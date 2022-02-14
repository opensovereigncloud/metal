# Build the manager binary
FROM core.harbor.onmetal.de/proxy/library/golang:1.17 as builder

ARG GOARCH=''
ARG GOPRIVATE
ARG GIT_USER
ARG GIT_PASSWORD

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN git config --global url."https://${GIT_USER}:${GIT_PASSWORD}@github.com".insteadOf "https://github.com"

RUN go mod download

# Copy the go source
COPY main.go main.go
COPY apis/ apis/
COPY controllers/ controllers/
COPY internal/ internal/

# Build
RUN GOMAXPROCS=1 CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH:-$(go env GOARCH)} go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
