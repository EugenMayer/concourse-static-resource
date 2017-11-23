#!/bin/sh
set -e
set -x
echo "---------running out tests---------------"
/opt/resource/out /tests/assets < /tests/assets/payload_out_version
/opt/resource/out /tests/assets < /tests/assets/payload_out_static
