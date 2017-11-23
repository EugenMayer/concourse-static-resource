#!/bin/sh
set -e
set -x
echo "---------running out tests---------------"
/opt/resource/out /tests/assets < /tests/assets/payload_out_version
echo "payload with version succeeded"
/opt/resource/out /tests/assets < /tests/assets/payload_out_static
echo "payload without placeholder succeeded"
/opt/resource/out /tests/assets < /tests/assets/payload_out_glob
echo "payload without source file global succeeded"
