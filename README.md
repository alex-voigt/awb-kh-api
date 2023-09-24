# AWB KH API

Da die vorhandene API des Abfallwirtschaftsbetriebs in Bad Kreuznach ein ziemlich bescheidenes Format ausliefert, habe ich diesen Proxy geschrieben. Dadurch wird die Nutzung z.B. in `Openhab` wesentlich vereinfacht.

## Daten aus der AWB API holen

Beispiel:

```bash
curl --header "Accept: application/json" \
    -XGET "https://blupassionsystem.de/city/rest/garbagestation/getAllGarbageCalendar?appId=44&cityId=23143&fromTime=1693519200000&houseNumberId=32023&regionId=14&streetId=23277&toTime=1696024800000"
```

## Ergebnis AWB API

```json
{
    "resultCode": 200,
    "resultMessage": "Success",
    "id": null,
    "token": null,
    "data": {
        "stations": [
            {
                "id": 169993,
                "fromDate": 1694210400000,
                "toDate": null,
                "beginAt": 47700000,
                "endAt": 51300000,
                "description": null,
                "active": false,
                "name": "Schadstoffsammlung",
                "color": null,
                "icon": "https://blupassionsystem.de/data/koblenz/Icon_Schad_Mobil.png",
                "defaultIcon": false,
                "locationId": 43428,
                "type": "STATION",
                "calendarType": "SELECTDATE",
                "locationName": "Schadstoffsammlung",
                "dateDefined": null,
                "frequency": null,
                "lat": 49.8393,
                "lng": 7.85678,
                "address": "Parkhaus Badeallee, Bad Kreuznach",
                "openingHour": true,
                "numberOfTons": null,
                "stationType": "MOBILE"
            },
            {
                "id": 170005,
                "fromDate": 1693778400000,
                "toDate": null,
                "beginAt": 49500000,
                "endAt": 53100000,
                "description": null,
                "active": false,
                "name": "Schadstoffsammlung",
                "color": null,
                "icon": "https://blupassionsystem.de/data/koblenz/Icon_Schad_Mobil.png",
                "defaultIcon": false,
                "locationId": 43431,
                "type": "STATION",
                "calendarType": "SELECTDATE",
                "locationName": "Schadstoffsammlung",
                "dateDefined": null,
                "frequency": null,
                "lat": 49.8355,
                "lng": 7.87016,
                "address": "Parkplatz Mannheimer Straße gegenüber Friedhof, Bad Kreuznach",
                "openingHour": true,
                "numberOfTons": null,
                "stationType": "MOBILE"
            }
        ]
    }
}           
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
        "typ": "Restmüll"
    },
    {
        "termin": "2022-02-24",
        "wochentag": "Donnerstag",
        "typ": "Bio"
    }
]
```

