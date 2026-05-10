#!/usr/bin/env bash

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOPATH_MOD=$(go env GOPATH)/pkg/mod

GO_MODULE=$(grep "^module" "$PROJECT_ROOT/go.mod" | awk '{print $2}')

dependencies=(
)

packages=(
"logger"
)

PROTOC_OPTS=()
for dep in "${dependencies[@]}"; do
  module="${dep%:*}"
  subdir="${dep#*:}"

  version=$(cd "$PROJECT_ROOT" && go list -m "$module" | awk '{print $2}')

  if [ -n "$subdir" ]; then
      PROTOC_OPTS+=("-I" "$GOPATH_MOD/$module@$version/$subdir")
  else
      PROTOC_OPTS+=("-I" "$GOPATH_MOD/$module@$version")
  fi
done

for package in "${packages[@]}"; do

mkdir -p ../protobuf/$package

protoc -I .\
 "${PROTOC_OPTS[@]}" \
 --go_out="../protobuf/$package" --go-grpc_out="../protobuf/$package" \
 --go_opt=module="${GO_MODULE}/protobuf/$package" \
 --go-grpc_opt=module="${GO_MODULE}/protobuf/$package" \
 ./${GO_MODULE}/$package/*.proto

done

