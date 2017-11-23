#!/bin/sh
set -e
set -x

echo "---------running in tests---------------"
/opt/resource/in /tmp/subfolder < /tests/assets/payload_in

if [ -s /tmp/subfolder/docker-sync-0.5.0.gem ]; then
  echo "in test succeeded, file has a size"
else
  echo "in test fail, file has 0 size"
  exit 1
fi

/opt/resource/in /tmp/subfolder < /tests/assets/payload_in_auth
echo "payload with auth succeeded"