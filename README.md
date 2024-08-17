# kimsufi-notifier

Collection of bash scripts used to check for kimsufi server availability.

It supports sending notifications via [OpsGenie](https://www.atlassian.com/software/opsgenie). Get your API key under Teams > Integrations.

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

Check availability of a specific server in one or multiple datacenters. Exit code is 0 if the server is available, 1 otherwise.

```
$ PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg bin/check.sh
> checking 22sk010 availability in fr,gra,rbx,sbg
> checked  22sk010 available    in fr,gra,rbx,sbg
```

A notification can be sent whenever a server is available.

```
$ PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg OPSGENIE_API_KEY=******** bin/check.sh
> checking 22sk010 availability in fr,gra,rbx,sbg
> checked  22sk010 available    in fr,gra,rbx,sbg
> sending notification
{"result":"Request will be processed","took":0.005,"requestId":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}
> notification sent
