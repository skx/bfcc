#!/bin/bash

# Install tools to test our code-quality.
go get -u golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
go get -u honnef.co/go/tools/cmd/staticcheck

# At this point failures cause aborts
set -e

# Run the static-check tool
t=$(mktemp)
staticcheck -checks all ./... > $t
if [ -s $t ]; then
    echo "Found errors via 'staticcheck'"
    cat $t
    rm $t
    exit 1
fi
rm $t

# Run the shadow-checker
echo "Launching shadowed-variable check .."
go vet -vettool=$(which shadow) ./...
echo "Completed shadowed-variable check .."

# Run golang tests
go test ./...

# Run the actual test-cases
make test
