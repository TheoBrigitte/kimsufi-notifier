#!/bin/bash

set -eu

a=$SMS_USER
a=$SMS_PASS
echo "> checking $HARDWARE in $COUNTRY"
curl -Ss "https://eu.api.ovh.com/1.0/dedicated/server/availabilities?country=$COUNTRY&hardware=$HARDWARE" | jq -e '.[].datacenters[] | select(.availability != "unavailable")'
curl -iXPOST https://smsapi.free-mobile.fr/sendmsg -d'{"msg":"'"$HARDWARE is available"'","user":"'"$SMS_USER"'","pass":"'"$SMS_PASS"'"}' -H'Content-Type: application/json'
