#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

declare -a tests=(
   "governance_rest"
   "holder_rest"
   "issuer_rest"
   "verifier_rest"
   "vc_rest"
   "w3c_workflow"
   "comparator_rest"
   "confidential-storage-hub"
   "vault_server"
)

TAGS=governance_rest

PWD=`pwd`
cd test/bdd

echo "---
Important: for these test to run correctly, you must ensure that your hosts file has the following entries:
127.0.0.1 testnet.orb.local
127.0.0.1 vault.kms.example.com
127.0.0.1 edv.rest.example.com
---
"
echo "Running adapter integration tests with tag=$TAGS"


go test -count=1 -v -cover . -p 1 -timeout=40m $TAGS

cd $PWD
