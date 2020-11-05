#!/bin/bash
#
# Check Kimsufi server availability
#
# Usages:
# 	COUNTRY=fr HARDWARE=1801sk12 SMS_USER=******** SMS_PASS=************** check.sh
# 	COUNTRY=fr HARDWARE=1801sk12 NO_SMS=true check.sh

set -eu

## required environement variables
_=$HARDWARE
_=$COUNTRY


OVH_URL=${OVH_URL:-https://eu.api.ovh.com/1.0}

# check availability from api
echo "> checking $HARDWARE in $COUNTRY"
if ! curl -Ss "${OVH_URL}/dedicated/server/availabilities?country=${COUNTRY}&hardware=${HARDWARE}" | jq -e '.[].datacenters[] | select(.availability != "unavailable")'; then
	echo "> $HARDWARE not available in $COUNTRY"
	exit 1
fi

echo "> $HARDWARE available in $COUNTRY"

# stop when NO_SMS variable is set
test ! -v NO_SMS

_=$SMS_USER
_=$SMS_PASS

FREEMOBILE_URL=${FREEMOBILE_URL:-https://smsapi.free-mobile.fr/sendmsg}

# send sms
# receiver phone number is the account holder
echo "> sending message"
curl -iXPOST "$FREEMOBILE_URL" -d'{"msg":"'"$HARDWARE is available"'","user":"'"$SMS_USER"'","pass":"'"$SMS_PASS"'"}' -H'Content-Type: application/json'
echo "> message sent"
