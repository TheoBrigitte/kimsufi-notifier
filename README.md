<p align="center">
    <img src="assets/kimsufi-logo.webp" alt="Kimsufi logo" height="100px">
</p>

<p align="center">
  <a href="https://github.com/TheoBrigitte/kimsufi-notifier/releases"><img src="https://img.shields.io/github/release/TheoBrigitte/kimsufi-notifier.svg" alt="Github release"></a>
  <a href="https://dl.circleci.com/status-badge/redirect/gh/TheoBrigitte/kimsufi-notifier/tree/main"><img src="https://dl.circleci.com/status-badge/img/gh/TheoBrigitte/kimsufi-notifier/tree/main.svg?style=svg" alt="CircleCI"></a>
</p>

## About <img src="./assets/info.svg" width="24">

[OVH Eco dedicated servers](https://eco.ovhcloud.com) are known for their low prices and high demand. As a result, they are often out of stock. This collection of bash scripts is used to check for server availability and send notifications when a server is available.

This is my playground to learn how the OVH API works, it is then ported into a Telegram Bot at [https://t.me/KimsufiNotifierBot](https://t.me/KimsufiNotifierBot) which is more user-friendly and provides more features.

## Features <img src="./assets/star.svg" width="24">

- [List available servers](#list-available-servers) from OVH Eco catalog
- [Check availability](#check-availability) of a specific server in one or multiple datacenters
- Send [notifications](#notifications-) to OpsGenie and/or Telegram when a server is available
- [Order a server](#order-a-server) directly from the command line

## Quickstart <img src="./assets/rocket.svg" width="24">

```
git clone git@github.com:TheoBrigitte/kimsufi-notifier.git
cd kimsufi-notifier
cp config.env.example config.env
bin/check.sh
```

## Run from CI &nbsp;<img src="./assets/rotate.svg" width="24">

See [RUN_IN_CI.md](RUN_IN_CI.md) for more information on how to run the check script using different CI services.

## Configuration <img src="./assets/configuration.svg" width="24">

Configuration is done through environment variables. The following variables are available:

- `COUNTRY`: country code to list servers from (e.g. `FR`)
- `CATEGORY`: server category to list servers from (e.g. `kimsufi`)
- `PLAN_CODE`: plan code to check availability for (e.g. `24ska01`)
- `DATACENTERS`: comma-separated list of datacenters to check availability in (e.g. `fr,gra,rbx,sbg`)
- `ENDPOINT`: OVH API endpoint to use (e.g. `ovh-eu`)
- `HEALTHCHECKS_IO_UUID`: UUID for healthchecks.io to ping after successful run

More details can be found in the [config.env.example](config.env.example) file.

## Notifications <img src="./assets/notifications.svg" width="24">

Notification(s) can be sent whenever a server is available. Either one or multiple notification services can be used.

Supported notification services:
- [Discord](https://discord.com/) via [Webhook](https://discord.com/developers/docs/resources/webhook)
- [Gotify](https://gotify.net/)
- [OpsGenie](https://www.atlassian.com/software/opsgenie) via [Alerts API](https://docs.opsgenie.com/docs/alert-api)
- [Telegram](https://telegram.org/) via [Bots API#sendMessage](https://core.telegram.org/bots/api#sendmessage)

In order to use a notification service, it is recommended to set its environment variables in the config file, see [config.env.example](config.env.example).

### Discord <img src="./assets/discord.svg" width="24">

In order to receive notifications for Discord, the appropriate environment variable must be set:

- `DISCORD_WEBHOOK`: Webhook URL to use for Discord notification service

See [Intro to Webhook](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks) for how to make and use a Discord Webhook.

### Gotify <img src="./assets/gotify.png" width="24">

In order to recieve notifications for Gotify, the appropriate environment variables must be set:

- `GOTIFY_URL`: URL to use for Gotify notification service
- `GOTIFY_TOKEN`: token to use for Gotify notification service
- `GOTIFY_PRIORITY`: prority for Gotify notification service

See [Gotify Push messages](https://gotify.net/docs/pushmsg) documentation for more information.

### OpsGenie <img src="./assets/opsgenie.svg" width="24">

In order to recieve notifications for OpsGenie, the appropriate environment variables must be set:

- `OPSGENIE_API_KEY`: API key to use OpsGenie notification service

See [OpsGenie API key](https://support.atlassian.com/opsgenie/docs/api-key-management/) or [OpsGenie integration](https://support.atlassian.com/opsgenie/docs/create-a-default-api-integration/) for more information.

Example with OpsGenie:
```
$ bin/check.sh --plan-code 24ska01
> checking 24ska01 availability in all datacenters
> checked  24ska01 available    in fr,gra,rbx,sbg
> sending OpsGenie notification
> sent    OpsGenie notification
```

### Telegram <img src="./assets/telegram.svg" width="24">


In order to recieve notifications for Telegram, the appropriate environment variables must be set:

- `TELEGRAM_CHAT_ID`: chat ID to use Telegram notification service
- `TELEGRAM_BOT_TOKEN`: bot token to use Telegram notification service

See [Telegram bot creation guide](https://core.telegram.org/bots/features#creating-a-new-bot) or [this Gist](https://gist.github.com/nafiesl/4ad622f344cd1dc3bb1ecbe468ff9f8a#file-how_to_get_telegram_chat_id-md)

Example with Telegram:
```
$ bin/check.sh --plan-code 24ska01
> checking 24ska01 availability in all datacenters
> checked  24ska01 available    in fr,gra,rbx,sbg
> sending Telegram notification
> sent    Telegram notification
```

### Examples <img src="./assets/bash.svg" width="24">

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

Check availability of a specific server in all datacenters.

```
$ bin/check.sh --plan-code 24sk50
> checking 24sk50 availability in all datacenters
> checked  24sk50 available    in bhs,fra,gra,rbx,sbg,waw datacenter(s)
```

#### Order a server

Place an order a specific server, the order is only placed and not payed for. The order can then be completed by following the URL provided in the output.

```
$ bin/order.sh --plan-code 24sk50 --datacenter fra --item-configuration region=europe
> cart created id=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> cart updated with item id=xxxxxxxxx
> item auto-configuration dedicated_os=none_64.en
> item user-configuration region=europe
> item user-configuration dedicated_datacenter=fra
> item option ram-32g-ecc-2400-24sk50
> item option bandwidth-300-24sk
> item option softraid-2x2000sa-24sk50
> cart assigned to account
> order completed url=https://www.ovh.com/cgi-bin/order/display-order.cgi?orderId=xxxxxxxxx&orderPassword=xxxxxxxxxx
 ```

 More info on scripts usage can be found in [USAGE.md](USAGE.md).
