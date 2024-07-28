# kimsufi-notifier

Notify on kimsufi server availability.

Uses [opsgenie](https://www.atlassian.com/software/opsgenie) to send notifications. Get your API key under Teams > Integrations.

### Usage

* Check availability

```
$ PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg hack/check.sh
> checking 22sk010 availability in fr,gra,rbx,sbg
> checked  22sk010 available    in fr,gra,rbx,sbg
> sending notification
{"result":"Request will be processed","took":0.005,"requestId":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}
> notification sent
```
