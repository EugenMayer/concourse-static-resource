#!/bin/sh
set -e
echo "---------running out tests---------------"

echo "-------------- payload_out_version"
/opt/resource/out /tests/assets < /tests/assets/payload_out_version | jq .
echo "payload with version succeeded"

echo "-------------- payload_out_static"
/opt/resource/out /tests/assets < /tests/assets/payload_out_static | jq .
echo "payload without placeholder succeeded"

echo "-------------- payload_out_glob"
/opt/resource/out /tests/assets < /tests/assets/payload_out_glob | jq .
echo "payload without source file global succeeded"
