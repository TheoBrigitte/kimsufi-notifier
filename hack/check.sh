#!/bin/bash

set -eu

## Required env variables
_=$HARDWARE
_=COUNTRY
_=SMS_USER
_=SMS_PASS

echo "> checking $HARDWARE in $COUNTRY"
curl -Ss "https://eu.api.ovh.com/1.0/dedicated/server/availabilities?country=$COUNTRY&hardware=$HARDWARE" | jq -e '.[].datacenters[] | select(.availability != "unavailable")'

echo "> $HARDWARE available in $COUNTRY"

echo "> sending message"
curl -iXPOST https://smsapi.free-mobile.fr/sendmsg -d'{"msg":"'"$HARDWARE is available"'","user":"'"$SMS_USER"'","pass":"'"$SMS_PASS"'"}' -H'Content-Type: application/json'
echo "> message sent"
