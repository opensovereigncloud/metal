#!/usr/bin/env bash

#
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
#

set -o errexit
set -o pipefail

# Turn colors in this script off by setting the NO_COLOR variable in your
# environment to any value:
#
# $ NO_COLOR=1 test.sh
NO_COLOR=${NO_COLOR:-""}
if [ -z "$NO_COLOR" ]; then
  header=$'\e[1;33m'
  reset=$'\e[0m'
else
  header=''
  reset=''
fi

function header_text {
  echo "$header$*$reset"
}

function setup_envtest_env {
  header_text "setting up env vars"

  # Setup env vars
  KUBEBUILDER_ASSETS=${KUBEBUILDER_ASSETS:-""}
  if [[ -z "${KUBEBUILDER_ASSETS}" ]]; then
    export KUBEBUILDER_ASSETS=$(pwd)/testbin/k8s/${VERSION}-$(go env GOOS)-$(go env GOARCH)
    header_text "KUBEBUILDER_ASSETS=${KUBEBUILDER_ASSETS}"
  fi
}

function fetch_envtest_tools {
    header_text "fetching kubebuilder asset"
    
    local basepath="$(go env GOPATH)"
    if [[ ! -f $basepath/bin/setup-envtest ]]; then
        go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
    fi

    $basepath/bin/setup-envtest use --bin-dir $(pwd)/testbin $VERSION
}