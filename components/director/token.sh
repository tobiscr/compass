#!/usr/bin/env bash

set -e

SCOPES="integration_system:read integration_system:write application:read application:write application_template:write application_template:read"
TOKEN_PAYLOAD='{"scopes": "'${SCOPES}'","tenant":"380da7fb-767e-45cf-8fcc-829f97655d1b"}'
ENCODED_TOKEN_PAYLOAD=$(echo -e ${TOKEN_PAYLOAD} | base64 | tr -d \\n)
INTERNAL_TOKEN="eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.${ENCODED_TOKEN_PAYLOAD//=}."

echo $INTERNAL_TOKEN
