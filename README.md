# kimsufi-notifier

Notify on kimsufi server availability.

Uses [go-ovh](https://github.com/ovh/go-ovh) client to talk to kimsufi API. No crendetials required.

Uses [free.fr](http://mobile.free.fr/) sms api. Get your credentials at: https://mobile.free.fr/moncompte/index.php?page=options

### Usage

* Get hardware and country code

```
$ ./hack/names.sh
> fetching data from https://www.kimsufi.com/fr/serveurs.xml
> fetched data
model=KS-1      hardware=1801sk12       country=FRA
model=KS-2      hardware=1801sk13       country=FRA
model=KS-3      hardware=1801sk14       country=FRA
model=KS-4      hardware=1801sk15       country=FRA
model=KS-5      hardware=1801sk16       country=FRA
model=KS-6      hardware=1801sk17       country=FRA
model=KS-7      hardware=1801sk18       country=FRA
model=KS-8      hardware=1801sk19       country=FRA
model=KS-9      hardware=1801sk20       country=FRA
model=KS-10     hardware=1801sk21       country=FRA
model=KS-11     hardware=1801sk22       country=FRA
model=KS-12     hardware=1801sk23       country=FRA
model=KS-1      hardware=1804sk12       country=BHS
model=KS-5      hardware=1804sk16       country=BHS
model=KS-7      hardware=1804sk18       country=BHS
model=KS-9      hardware=1804sk20       country=BHS
model=KS-10     hardware=1804sk21       country=BHS
model=KS-11     hardware=1804sk22       country=BHS
model=KS-12     hardware=1804sk23       country=BHS
```

* Check availability

```
$ HARDWARE=1801sk12 COUNTRY=FRA NO_SMS=true ./hack/check.sh
> checking 1801sk12 in FRA
> 1801sk12 not available in FRA
```
