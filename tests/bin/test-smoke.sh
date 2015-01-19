#!/bin/bash
#
# Preps a test environment and runs `make test-smoke`
# against artifacts produced from the current source tree
#

# fail on any command exiting non-zero
set -eo pipefail

# absolute path to current directory
export THIS_DIR=$(cd $(dirname $0); pwd)

# use the built client binaries
export PATH=$DEIS_ROOT/deisctl:$DEIS_ROOT/client/dist:$PATH

export DEIS_NUM_INSTANCES=3
make discovery-url
vagrant up --provider virtualbox

until deisctl list >/dev/null 2>&1; do
    sleep 1
done

make dev-release

# configure platform settings
deisctl config platform set domain=$DEIS_TEST_DOMAIN
deisctl config platform set sshPrivateKey=$DEIS_TEST_SSH_KEY

time deisctl install platform
time deisctl start platform

time make test-smoke
