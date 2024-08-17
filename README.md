# kimsufi-notifier

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/TheoBrigitte/kimsufi-notifier/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/TheoBrigitte/kimsufi-notifier/tree/main)
[![GitHub release](https://img.shields.io/github/release/TheoBrigitte/kimsufi-notifier.svg)](https://github.com/TheoBrigitte/kimsufi-notifier/releases)

Collection of bash scripts used to check for kimsufi server availability.

It supports sending notifications to OpsGenie and Telegram when a server is available.

### Usage

#### List available servers

List all available servers from OVH Eco catalog in a specific country.

```
$ COUNTRY=FR bin/list.sh
> fetching servers in FR
> fetched  servers
PlanCode                Category      Name                 Price (EUR)
22sk010                 kimsufi       KS-1                 4.99
22sk011                 kimsufi       KS-2                 8.99
22sk012                 kimsufi       KS-3                 13.99
22sk020                 kimsufi       KS-4                 15.99
22sk030                 kimsufi       KS-5                 19.99
...
```

#### Check availability

Check availability of a specific server in one or multiple datacenters.

```
$ PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg bin/check.sh
> checking 22sk010 availability in fr,gra,rbx,sbg
> checked  22sk010 available    in fr,gra,rbx,sbg
```

Notification(s) can be sent whenever a server is available. Either one or multiple notification services can be used.

Supported notification services:
- [OpsGenie](https://www.atlassian.com/software/opsgenie) via [Alerts API](https://docs.opsgenie.com/docs/alert-api)
- [Telegram](https://telegram.org/) via [Bots API#sendMessage](https://core.telegram.org/bots/api#sendmessage)

Example with OpsGenie:
```
$ PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg OPSGENIE_API_KEY=******** bin/check.sh
> checking 22sk010 availability in fr,gra,rbx,sbg
> checked  22sk010 available    in fr,gra,rbx,sbg
> sending OpsGenie notification
> sent    OpsGenie notification
```

Example with Telegram:
```
$ PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg TELEGRAM_BOT_TOKEN=******** TELEGRAM_CHAT_ID=******** bin/check.sh
> checking 22sk010 availability in fr,gra,rbx,sbg
> checked  22sk010 available    in fr,gra,rbx,sbg
> sending Telegram notification
> sent    Telegram notification
```
