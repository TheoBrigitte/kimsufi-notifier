<p align="center">
    <img src="assets/kimsufi-logo.webp" alt="Kimsufi logo" height="100px">
</p>

<p align="center">
  <a href="https://github.com/TheoBrigitte/kimsufi-notifier/releases"><img src="https://img.shields.io/github/release/TheoBrigitte/kimsufi-notifier.svg" alt="Github release"></a>
  <a href="https://dl.circleci.com/status-badge/redirect/gh/TheoBrigitte/kimsufi-notifier/tree/main"><img src="https://dl.circleci.com/status-badge/img/gh/TheoBrigitte/kimsufi-notifier/tree/main.svg?style=svg" alt="CircleCI"></a>
</p>

## About

[OVH Eco dedicated servers](https://eco.ovhcloud.com) are known for their low prices and high demand. As a result, they are often out of stock. This collection of bash scripts is used to check for server availability and send notifications when a server is available.

## Features

- List available servers from OVH Eco catalog in a specific country
- Check availability of a specific server in one or multiple datacenters
- Send notifications to OpsGenie and/or Telegram when a server is available

## Quickstart

```
git clone git@github.com:TheoBrigitte/kimsufi-notifier.git
cd kimsufi-notifier
cp config.env.example config.env
bin/check.sh
```

## Configuration

Configuration is done through environment variables. The following variables are available:

- `COUNTRY`: country code to list servers from (e.g. `FR`)
- `PLAN_CODE`: plan code to check availability for (e.g. `22sk010`)
- `DATACENTERS`: comma-separated list of datacenters to check availability in (e.g. `fr,gra,rbx,sbg`)
- `OPSGENIE_API_KEY`: API key to use OpsGenie notification service
- `TELEGRAM_CHAT_ID`: chat ID to use Telegram notification service
- `TELEGRAM_BOT_TOKEN`: bot token to use Telegram notification service

More details can be found in the [config.env.example](config.env.example) file.

## Usage

```
$ bin/list.sh -h
Usage: list.sh

List servers from OVH Eco (including Kimsufi) catalog

Arguments
  -c, --country    Country code (required)
                     Allowed values for ovh-eu: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
                     Allowed values for ovh-ca: ASIA, AU, CA, IN, QC, SG, WE, WS
                     Allowed values for ovh-us: US
  --category       Server category (default all)
                     Allowed values: kimsufi, soyoustart, rise, uncategorized
  -e, --endpoint   OVH API endpoint (default: ovh-eu)
                     Allowed values: ovh-eu, ovh-ca, ovh-us
  -d, --debug      Enable debug mode (default: false)
  -h, --help       Display this help message

  Arguments can also be set as environment variables see config.env.example
  Command line arguments take precedence over environment variables

Example:
    list.sh --country FR
    list.sh --country FR --category kimsufi
```

```
$ bin/check.sh -h
Usage: check.sh

Check OVH Eco (including Kimsufi) server availability

Arguments
  -p, --plan-code        Plan code to check (e.g. 24ska01)
  --datacenters          Comma-separated list of datacenters to check availability for (default all)
                           Allowed values: bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw (non exhaustive list)
  -e, --endpoint         OVH API endpoint (default: ovh-eu)
                           Allowed values: ovh-eu, ovh-ca, ovh-us
  -d, --debug            Enable debug mode (default: false)
  -h, --help             Display this help message

  Arguments can also be set as environment variables see config.env.example
  Command line arguments take precedence over environment variables

Environment variables
    OPSGENIE_API_KEY      API key for OpsGenie to receive notifications
    TELEGRAM_BOT_TOKEN    Bot token for Telegram to receive notifications
    TELEGRAM_CHAT_ID      Chat ID for Telegram to receive notifications
    HEALTHCHECKS_IO_UUID  UUID for healthchecks.io to ping after successful run

Example:
  check.sh --plan-code 24ska01
  check.sh --plan-code 24ska01 --datacenters fr,gra,rbx,sbg
```

### Examples

#### List available servers

List servers from OVH Eco catalog in a specific country and from a specific category.

```
$ bin/list.sh --country FR --category kimsufi
> fetching servers in FR
> fetched  servers
PlanCode              Category    Name                            Price (EUR)
24ska01               kimsufi     KS-A | Intel i7-6700k           4.99
24sk10                kimsufi     KS-1 | Intel Xeon-D 1520        16.99
24sk40                kimsufi     KS-4 | Intel Xeon-E3 1230 v6    16.99
24sk20                kimsufi     KS-2 | Intel Xeon-D 1540        18.99
24sk30                kimsufi     KS-3 | Intel Xeon-E3 1245 v5    18.99
24sk50                kimsufi     KS-5 | Intel Xeon-E3 1270 v6    25.99
24sk60                kimsufi     KS-6 | AMD Epyc 7351P           38.99
24skstor01            kimsufi     KS-STOR | Intel Xeon-D 1521     49.99
22skgameapac01-sgp    kimsufi     KS-GAME-APAC-1-1                53.99
22skgameapac01-syd    kimsufi     KS-GAME-APAC-1-1                53.99
24ska01-syd           kimsufi     KS-A | Intel i7-6700k           53.99
24sk30-sgp            kimsufi     KS-3 | Intel Xeon-E3 1245 v5    54.99
24sk30-syd            kimsufi     KS-3 | Intel Xeon-E3 1245 v5    54.99
24sk40-syd            kimsufi     KS-4 | Intel Xeon-E3 1230 v6    54.99
24sk70                kimsufi     KS-7 | AMD Epyc 7451            76.99
```

#### Check availability

Check availability of a specific server in one all datacenters.

```
$ bin/check.sh --plan-code 24ska01
> checking 24ska01 availability in all datacenters
> checked  24ska01 unavailable  in fr,gra,rbx,sbg
```

## Notifications

Notification(s) can be sent whenever a server is available. Either one or multiple notification services can be used.

Supported notification services:
- [OpsGenie](https://www.atlassian.com/software/opsgenie) via [Alerts API](https://docs.opsgenie.com/docs/alert-api)
- [Telegram](https://telegram.org/) via [Bots API#sendMessage](https://core.telegram.org/bots/api#sendmessage)

In order to receive notifications the appropriate environment variables must be set:

- `OPSGENIE_API_KEY`: API key to use OpsGenie notification service
- `TELEGRAM_CHAT_ID`: chat ID to use Telegram notification service
- `TELEGRAM_BOT_TOKEN`: bot token to use Telegram notification service

More details can be found in the [config.env.example](config.env.example) file.

### Examples

Example with OpsGenie:
```
$ export OPSGENIE_API_KEY=********
$ bin/check.sh --plan-code 24ska01
> checking 24ska01 availability in all datacenters
> checked  24ska01 available    in fr,gra,rbx,sbg
> sending OpsGenie notification
> sent    OpsGenie notification
```

Example with Telegram:
```
$ export TELEGRAM_BOT_TOKEN=********
$ export TELEGRAM_CHAT_ID=********
$ bin/check.sh --plan-code 24ska01
> checking 24ska01 availability in all datacenters
> checked  24ska01 available    in fr,gra,rbx,sbg
$ PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg TELEGRAM_BOT_TOKEN=******** TELEGRAM_CHAT_ID=******** bin/check.sh
> checking 22sk010 availability in fr,gra,rbx,sbg
> checked  22sk010 available    in fr,gra,rbx,sbg
> sending Telegram notification
> sent    Telegram notification
```
