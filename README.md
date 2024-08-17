<p align="center">
    <img src="assets/kimsufi-logo.webp" alt="Kimsufi logo" height="100px">
</p>

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/TheoBrigitte/kimsufi-notifier/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/TheoBrigitte/kimsufi-notifier/tree/main)
[![GitHub release](https://img.shields.io/github/release/TheoBrigitte/kimsufi-notifier.svg)](https://github.com/TheoBrigitte/kimsufi-notifier/releases)

## About

[OVH Eco dedicated servers](https://eco.ovhcloud.com) are known for their low prices and high demand. As a result, they are often out of stock. This collection of bash scripts is used to check for server availability and send notifications when a server is available.

## Features

- List available servers from OVH Eco catalog in a specific country
- Check availability of a specific server in one or multiple datacenters
- Send notifications to OpsGenie and/or Telegram when a server is available

### Quickstart

```
git clone git@github.com:TheoBrigitte/kimsufi-notifier.git
cd kimsufi-notifier
cp config.env.example config.env
bin/check.sh
```

### Configuration

Configuration is done through environment variables. The following variables are available:

- `COUNTRY`: country code to list servers from (e.g. `FR`)
- `PLAN_CODE`: plan code to check availability for (e.g. `22sk010`)
- `DATACENTERS`: comma-separated list of datacenters to check availability in (e.g. `fr,gra,rbx,sbg`)
- `OPSGENIE_API_KEY`: API key to use OpsGenie notification service
- `TELEGRAM_CHAT_ID`: chat ID to use Telegram notification service
- `TELEGRAM_BOT_TOKEN`: bot token to use Telegram notification service

More details can be found in the [config.env.example](config.env.example) file.

### Usage

#### List available servers

List servers from OVH Eco catalog in a specific country.

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
