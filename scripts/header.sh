#!/bin/bash

ORG_ID="$1"
ACCOUNT_NUMBER="$2"


function print_out_usage {
	cat <<EOF
Usage: ./scripts/header.sh <org_id> [account_number]
EOF
}

function error {
	local err=$?
	print_out_usage >&2
	printf "error: %s\n" "$*" >&2
	exit $err
}

[ "${ORG_ID}" != "" ] || error "ORG_ID is required and cannot be empty"

if [ "$( uname -s )" == "Darwin" ]; then
BASE64ENC="base64"
else
BASE64ENC="base64 -w0"
fi
export BASE64ENC

ENC="$(echo "{\"identity\":{\"type\":\"Associate\",\"account_number\":\"$2\",\"internal\":{\"org_id\":\"$1\"}}}" | $BASE64ENC )"
echo "x-rh-identity: $ENC"
