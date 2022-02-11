# AWB KH API

Da die vorhandene API des Abfallwirtschaftsbetriebs in Bad Kreuznach ein ziemlich bescheidenes Format ausliefert habe ich diesen Proxy geschrieben. Dadurch wird die Nutzung z.B. in `Openhab` wesentlich vereinfacht.

## Daten aus der AWB API holen

Beispiel:

```bash
curl 'https://app.awb-bad-kreuznach.de/api/loadDates.php' \
  --data-raw 'ort=Bad Kreuznach&strasse=Europaplatz' \
  --compressed
```

## Ergebnis AWB API

```json
[
    {
        "id": "17794",
        "state": "1",
        "termin": "2022-02-03",
        "restmuell": "1",
        "bio": "0",
        "wert": "0",
        "papier": "0",
        "TRIM(wochentag)": "Donnerstag",
        "wochentag": "Donnerstag"
    },
    {
        "id": "16973",
        "state": "1",
        "termin": "2022-02-10",
        "restmuell": "0",
        "bio": "0",
        "wert": "0",
        "papier": "7",
        "TRIM(wochentag)": "Donnerstag",
        "wochentag": "Donnerstag"
    }
]
```

## Proxy benutzen

```bash
curl -XGET "http://<host>:8010/getDates?ort=Bad%20Kreuznach&strasse=Europaplatz"
```

## Ergebnis

Es werden alle vergangenen Termine rausgefiltert und nur die nächsten 10 Termine im folgenden Format ausgegeben:

```json
[
    {
        "termin": "2022-02-17",
        "wochentag": "Donnerstag",
        "typ": "Restmüll"
    },
    {
        "termin": "2022-02-24",
        "wochentag": "Donnerstag",
        "typ": "Bio"
    }
]
```

