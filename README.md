<p align="center">
    <img src="assets/kimsufi-logo.webp" alt="Kimsufi logo" height="100px">
</p>

<p align="center">
  <a href="https://github.com/TheoBrigitte/kimsufi-notifier/releases"><img src="https://img.shields.io/github/release/TheoBrigitte/kimsufi-notifier.svg" alt="Github release"></a>
  <a href="https://github.com/TheoBrigitte/kimsufi-notificer/actions/workflows/go.yaml"><img src="https://github.com/TheoBrigitte/kimsufi-notifier/actions/workflows/go.yaml/badge.svg" alt="Github action"></a>
</p>

## About <img src="./assets/info.svg" width="24">

[OVH Eco dedicated servers](https://eco.ovhcloud.com) are known for their low prices and high demand. As a result, they are often out of stock. This program is used to check for server availability, and to place order for servers. Note: the previous Bash script collection was moved to [`bash`](https://github.com/TheoBrigitte/kimsufi-notifier/tree/bash) branch.

This is my playground to learn how the OVH API works, it is then ported into a Telegram Bot at [https://t.me/KimsufiNotifierBot](https://t.me/KimsufiNotifierBot) which is more user-friendly and provides more features.

## Features <img src="./assets/star.svg" width="24">

- [List available servers](#list-available-servers) from OVH Eco catalog
- [Check availability](#check-availability) of a specific server in one or multiple datacenters
- [Order a server](#order-a-server) directly from the command line

## Quickstart <img src="./assets/rocket.svg" width="24">

### Using pre-built binaries

Download the latest release from the [Github releases page](https://github.com/TheoBrigitte/kimsufi-notifier/releases).

### Using Go

```
go install github.com/TheoBrigitte/kimsufi-notifier
kimsufi-notifier
```

### Examples <img src="./assets/bash.svg" width="24">

#### List available servers

List servers from OVH Eco catalog in a specific country and from a specific category.

```
$ kimsufi-notifier list --category kimsufi
planCode          category    name                             price        status         datacenters
--------          --------    ----                             -----        ------         -----------
24ska01           Kimsufi     KS-A | Intel i7-6700k            4.99 EUR     unavailable
25skle01          Kimsufi     KS-LE-1                          9.99 EUR     available      bhs
25skleb01         Kimsufi     KS-LE-B                          9.99 EUR     available      bhs, fra, gra, waw
25sklea01         Kimsufi     KS-LE-A                          9.99 EUR     available      bhs, fra, gra, waw
25skled01         Kimsufi     KS-LE-D                          12.99 EUR    unavailable
25sklec01         Kimsufi     KS-LE-C                          12.99 EUR    unavailable
25sklee01         Kimsufi     KS-LE-E                          14.99 EUR    unavailable
25skle02          Kimsufi     KS-LE-2                          15.99 EUR    unavailable
24sk40            Kimsufi     KS-4 | Intel Xeon-E3 1230 v6     16.99 EUR    available      bhs, fra, gra
24sk10            Kimsufi     KS-1 | Intel Xeon-D 1520         16.99 EUR    available      bhs
24sk20            Kimsufi     KS-2 | Intel Xeon-D 1540         18.99 EUR    unavailable
24sk30            Kimsufi     KS-3 | Intel Xeon-E3 1245 v5     18.99 EUR    unavailable
...
```

#### Check availability

Check availability of a specific server in all datacenters.

```
$ kimsufi-notifier check --plan-code 25skle01
planCode    memory                storage              status         datacenters
--------    ------                -------              ------         -----------
25skle01    ram-16g-noecc-1333    softraid-2x480ssd    unavailable
25skle01    ram-16g-noecc-1333    softraid-2x960ssd    unavailable
25skle01    ram-16g-noecc-1333    softraid-3x2000sa    unavailable
25skle01    ram-16g-noecc-1333    softraid-3x480ssd    unavailable
25skle01    ram-32g-noecc-1333    softraid-2x480ssd    unavailable
25skle01    ram-32g-noecc-1333    softraid-2x960ssd    unavailable
25skle01    ram-32g-noecc-1333    softraid-3x2000sa    available      bhs
25skle01    ram-32g-noecc-1333    softraid-3x480ssd    unavailable
```

#### Order a server

Place an order a specific server, the order is only placed and not payed for. The order can then be completed by following the URL provided in the output.

```
$ kimsufi-notifier order --plan-code 25skle01 --datacenter bhs --item-option memory=ram-32g-noecc-1333-25skle01,storage=softraid-3x2000sa-25skle01
> cart created id=dd413a3a-1eed-473c-bbe4-2a1c4f3d02c0
> cart item added id=299679179
> cart item configured: dedicated_os=none_64.en
> cart item configured: dedicated_datacenter=bhs
> cart item configured: region=europe
> cart option set: memory=ram-32g-noecc-1333-25skle01
> cart option set: storage=softraid-3x2000sa-25skle01
> cart option set: bandwidth=bandwidth-300-unguaranteed-25skle
> cart assigned
> order completed: url=https://www.ovh.com/cgi-bin/order/display-order.cgi?orderId=xxxxxxxxx&orderPassword=xxxxxxxxxx
 ```

 More info on usage can be found in [USAGE.md](USAGE.md).
