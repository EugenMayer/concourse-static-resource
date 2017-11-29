#!/bin/sh
set -e
#set -x

echo "---------running in tests---------------"

echo "-------------------  dynamic url with extract"
/opt/resource/in /tmp/subfolder < /tests/assets/payload_in_dynamic_extract | jq .

if [ -s /tmp/subfolder/mysql-connector-java-5.1.44/mysql-connector-java-5.1.44-bin.jar ]; then
  echo "in test succeeded, file has a size"
else
  echo "in test fail, file has 0 size"
  exit 1
fi

echo "-------------------  static url"
/opt/resource/in /tmp/subfolder < /tests/assets/payload_in_static_uri | jq .

if [ -s /tmp/subfolder/docker-sync-0.5.0.gem ]; then
  echo "in test succeeded, file has a size"
else
  echo "in test fail, file has 0 size"
  exit 1
fi

echo "-------------------  out to in handover.."
/opt/resource/in /tmp/subfolder < /tests/assets/payload_in_out_to_in_handover | jq .

if [ -s /tmp/subfolder/docker-sync-0.5.0.gem ]; then
  echo "in test succeeded, file has a size"
else
  echo "in test fail, file has 0 size"
  exit 1
fi

echo "------------------- check to in handover.."
/opt/resource/in /tmp/subfolder < /tests/assets/payload_in_check_to_in_handover | jq .

if [ -s /tmp/subfolder/docker-sync-0.5.0.gem ]; then
  echo "in test succeeded, file has a size"
else
  echo "in test fail, file has 0 size"
  exit 1
fi

echo "------------------- authen connection with static url.."
/opt/resource/in /tmp/subfolder < /tests/assets/payload_in_static_uri_auth | jq .
echo "payload with auth succeeded"