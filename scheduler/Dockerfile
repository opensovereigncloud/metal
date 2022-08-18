# /*
# Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# Build the manager binary
FROM golang:1.19 as builder

ARG GOARCH
ARG GOPRIVATE

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

COPY hack/ hack/
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN --mount=type=ssh --mount=type=secret,id=github_pat GITHUB_PAT_PATH=/run/secrets/github_pat ./hack/setup-git-redirect.sh \
  && mkdir -p -m 0600 ~/.ssh \
  && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts \
  && go mod download

# Copy the go source
COPY scheduler/application/scheduler.go main.go
COPY common common
COPY apis apis/
COPY pkg pkg/
COPY internal internal/
COPY scheduler scheduler/

# Build
RUN GOMAXPROCS=1 CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -a -o metal-api-scheduler main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/metal-api-scheduler .
USER 65532:65532

ENTRYPOINT ["/metal-api-scheduler"]
