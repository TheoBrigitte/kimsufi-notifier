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
	exit 0
fi

echo "> $HARDWARE available in $COUNTRY"

# stop when NO_SMS variable is set
test ! -v NO_SMS

_=$SMS_USER
_=$SMS_PASS
_=$SMS_DELAY
_=$SMS_STATE_FILE

# only send 1 sms every $SMS_DELAY
if test -e $SMS_STATE_FILE && test $(($(date +%s) - $(stat -c %Y $SMS_STATE_FILE))) -lt $SMS_DELAY; then
	echo "> sms already sent at $(stat -c %y $SMS_STATE_FILE) (delay: "$(TZ=UTC0 printf '%(%Hh%Mm%Ss)T\n' "$SMS_DELAY")")"
	exit 0
fi

touch $SMS_STATE_FILE

FREEMOBILE_URL=${FREEMOBILE_URL:-https://smsapi.free-mobile.fr/sendmsg}

# send sms
# receiver phone number is the account holder
message="$HARDWARE is available\nhttps://www.kimsufi.com/fr/commande/kimsufi.xml?reference=$HARDWARE ."
echo "> sending message"
curl -iXPOST "$FREEMOBILE_URL" -d'{"msg":"'"$message"'","user":"'"$SMS_USER"'","pass":"'"$SMS_PASS"'"}' -H'Content-Type: application/json'
echo "> message sent"
