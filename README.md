# AWB KH API

Da die vorhandene API des Abfallwirtschaftsbetriebs in Bad Kreuznach ein ziemlich bescheidenes Format ausliefert, habe ich diesen Proxy geschrieben. Dadurch wird die Nutzung z.B. in `Openhab` wesentlich vereinfacht.

## Daten aus der AWB API holen

Beispiel:

```bash
curl --header "Accept: application/json" \
    -XGET "https://blupassionsystem.de/city/rest/garbageorte/getAllGarbageCalendar?appId=44&cityId=23143&fromTime=1704063600000&houseNumberId=32023&regionId=14&streetId=23277&toTime=1706655600000"
```

## Proxy benutzen

```bash
curl -XGET "http://<host>:8010/getDates?appId=44&cityId=23143&houseNumberId=32023&regionId=14&streetId=23277"
```

## Ergebnis

Es werden alle vergangenen Termine rausgefiltert und nur die nächsten 10 Termine im folgenden Format ausgegeben:

```json
[
    {
        "termin": "2022-02-17",
        "wochentag": "Donnerstag",
        "typ": "Restabfall"
    },
    {
        "termin": "2022-02-24",
        "wochentag": "Donnerstag",
        "typ": "Bioabfall"
    }
]
```

