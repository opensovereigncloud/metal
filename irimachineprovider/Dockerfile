FROM golang:1.21.6 as builder
ARG GOARCH

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
COPY hack/ hack/
RUN --mount=type=ssh --mount=type=secret,id=github_pat GITHUB_PAT_PATH=/run/secrets/github_pat ./hack/setup-git-redirect.sh \
  && mkdir -p -m 0600 ~/.ssh \
  && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts \
  && go mod download

COPY apis/ apis/
COPY applyconfiguration/ applyconfiguration/
COPY irimachineprovider/ irimachineprovider/
COPY pkg/ pkg/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -a -o /metal-provider irimachineprovider/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
USER 65532:65532
ENTRYPOINT ["/metal-provider"]

COPY --from=builder /metal-provider .
