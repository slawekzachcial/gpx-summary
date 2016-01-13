GPX Summary
===========

Command line utility and `golang` package to extract summary information from
GPX tracks.

The following information is extracted:
* Start time
* Duration
* Distance
* Total ascent
* Total descent

In addition, the following data is calculated:
* Average speed
* Average pace

Utility Usage
-------------

```
$> gpx-summary exercise-2015-11-1*.gpx
File: exercise-2015-11-10.gpx
Time: Tue, 10 Nov 2015 11:58:40 CET
Duration: 1h12m42s
Distance: 11.3km
Ascent: 387m
Descent: 400m
Speed: 9.4km/h
Pace: 6.4min/km
---
File: exercise-2015-11-15.gpx
Time: Sun, 15 Nov 2015 10:21:06 CET
Duration: 1h34m4s
Distance: 13.7km
Ascent: 578m
Descent: 570m
Speed: 8.8km/h
Pace: 6.9min/km
```
To show the data as a table with tab-separated columns use `-t` option:
```
$> gpx-summary -t exercise-2015-11-1*.gpx
File	Time	Duration [min]	Distance [km]	Ascent [m]	Descent [m]	Speed [km/h]	Pace [min/km]
exercise-2015-11-10.gpx	Tue, 10 Nov 2015 11:58:40 CET	72.70	11.3	387	400	9.4	6.4
exercise-2015-11-15.gpx	Sun, 15 Nov 2015 10:21:06 CET	94.07	13.7	578	570	8.8	6.9
```

Golang Package Usage
--------------------

The package usage is shown in the [command line utility](cmd/gpx-summary) sources.

