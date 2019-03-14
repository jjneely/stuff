Package TSDB
============

A Go package that creates Prometheus TSDB data for testing purposes.

This code originally comes from the [Thanos][1] authors and is licensed under
the Apache License 2.0.  Modifications made by Jack Neely and Jarod Watkins
of [42 Lines, Inc.][2] include:

* Produce time series with labels similar to Prometheus scraping many
  instances of the same application
* Values recorded is the current time in nanoseconds to simulate ever
  increasing Counter type metrics.

Jack Neely <jjneely@42lines.net>
2019/03/13
