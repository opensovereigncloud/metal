#!/bin/sh
set -eu

VGOPATH="$(mktemp -d)"
MODELSSCHEMA="$(mktemp)"
trap 'rm -rf "$VGOPATH" "$MODELSSCHEMA"' EXIT
go mod download && go run github.com/ironcore-dev/vgopath -o "$VGOPATH"
GOROOT="${GOROOT:-"$(go env GOROOT)"}"
export GOROOT
GOPATH="$VGOPATH"
export GOPATH
GO111MODULE=off
export GO111MODULE

APIS_APPLYCONFIGURATION='github.com/ironcore-dev/metal/api/v1alpha1'
APIS_OPENAPI="k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource,k8s.io/api/core/v1,$APIS_APPLYCONFIGURATION"

go run k8s.io/code-generator/cmd/openapi-gen \
  --output-base "$GOPATH/src" \
  --go-header-file hack/boilerplate.go.txt \
  --input-dirs "$APIS_OPENAPI" \
  --output-package "github.com/ironcore-dev/metal/client/openapi" \
  -O zz_generated.openapi \
  --report-filename "/dev/null"

go run github.com/ironcore-dev/metal/internal/tools/models-schema > "$MODELSSCHEMA"
go run k8s.io/code-generator/cmd/applyconfiguration-gen \
  --output-base "$GOPATH/src" \
  --go-header-file hack/boilerplate.go.txt \
  --input-dirs "$APIS_APPLYCONFIGURATION" \
  --openapi-schema "$MODELSSCHEMA" \
  --output-package "github.com/ironcore-dev/metal/client/applyconfiguration"
