#!/bin/sh
set -eu

VGOPATH="$(mktemp -d)"
MODELSSCHEMA="$(mktemp)"
trap 'rm -rf "$VGOPATH" "$MODELSSCHEMA"' EXIT
go mod download && bin/vgopath -o "$VGOPATH"
GOROOT="${GOROOT:-"$(go env GOROOT)"}"
export GOROOT
GOPATH="$VGOPATH"
export GOPATH
GO111MODULE=off
export GO111MODULE

APIS_APPLYCONFIGURATION='github.com/ironcore-dev/metal/apis/metal/v1alpha4'
APIS_OPENAPI="k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource,k8s.io/api/core/v1,$APIS_APPLYCONFIGURATION"

bin/openapi-gen \
  --output-base "$GOPATH/src" \
  --go-header-file hack/boilerplate.go.txt \
  --input-dirs "$APIS_OPENAPI" \
  --output-package "github.com/ironcore-dev/metal/openapi" \
  -O zz_generated.openapi \
  --report-filename "openapi/api_violations.report"

bin/models-schema --openapi-package "github.com/ironcore-dev/metal/openapi" --openapi-title "metal-api" > "$MODELSSCHEMA"
bin/applyconfiguration-gen \
  --output-base "$GOPATH/src" \
  --go-header-file hack/boilerplate.go.txt \
  --input-dirs "$APIS_APPLYCONFIGURATION" \
  --openapi-schema "$MODELSSCHEMA" \
  --output-package "github.com/ironcore-dev/metal/applyconfiguration"
