Generate Test Prometheus TSDB Data
==================================

This tool generates Prometheus 2.x TSDB blocks containing test data.  If you
have ever wanted to test [Prometheus][1] or [Thanos][2] behavior on a specific
time series foot print, say a year's worth of data, then this tool is for you.

    $ ./historygen -h
    Usage of ./historygen:
      -b duration
            TSDB block length (default 2h0m0s)
      -c int
            Number of time series to generate (default 1)
      -C int
            Total number of time series to generate using multiple invocations (only needed for the zero padding of instance names)
      -n int
            Start index for time series instance names (default 0)
      -d duration
            Time duration of historical data to generate (default 720h0m0s)
      -s duration
            Time shift of historical data to generate (default 0h)
      -i duration
            Duration between samples (default 15s)
      -o string
            Output directory to generate TSDB blocks in (default "data/")

Historygen will create the directory named by the `-o` option and populate
it with TSDB blocks containing time series data.  Each block produced will
have the duration specified by `-b` (block length).  The total amount of
history to generate is controlled by `-d` and `-s` for duration.  So, for 720 hours
of history (about a month) the tool would generate 360 TSDB blocks spanning
the time range from the time the command was run until 720 hours ago.

All time series generated have the name `test`, a `job="testdata"` label,
and an `instance="test-metric-XXX"` where XXX is a zero based integer
(unless otherwise specified by the `-n` option) that
uniquely defines each time series as requested by the `-c` option.  Using
`-c 500` would produce time series that look like:

    test{instance="test-metric-000", job="testdata"}
    test{instance="test-metric-001", job="testdata"}
    ...
    test{instance="test-metric-499", job="testdata"}

The `-i` option defines the simulated scrape interval of the data.  No jitter
is added at this writing and the data points are spaced evenly apart at this
duration.  The default matches Prometheus's default of a 15 second scrape
interval.  The data points stored in each time series is a nonosecond
resolution timestamp of when the data was created.  This simulates Prometheus's
Counters.

Why is this Stuff Called Stuff?
===============================

I'm out of names, and I needed a quick place to store some random bits of
Go code that I found useful.

Where Did This Originally Come From?
====================================

The Thanos authors wrote some interesting bench mark tools that this idea
and the `tsdb` package was originally lifted from.  Its been modified to
create the time series naming/labeling scheme as documented above, and to
use nanosecond timestamps as the values for each stored data point.

How Do You Get and Build This?
==============================

With Go, of course.  Golang 1.11 or better is required most likely.

    $ go get -u github.com/jjneely/stuff/cmd/historygen
    $ cd ~/go/src/github.com/jjneely/stuff/cmd/historygen
    $ go build

Ok, I Generated a Pile of Data, What Next?
==========================================

Once the data is generated, insert it into Prometheus.  Ideally, the block
length specified matches the uncompacted block lengths that the Prometheus
instance is creating or is configured to create.  Shut down Prometheus,
copy the TSDB blocks into the data directory, start Prometheus.  Note that
TSDB blocks cannot overlap time ranges with other blocks (although Prometheus
2.8 does introduce some support for this).  The safe thing to do is wipe
Prometheus's TSDB directory, copy in the generated data, then start
Prometheus.  If the Prometheus WAL is left intact, Prometheus will later cut
a new TSDB block from the WAL that will overlap with the generated test data.
Perhaps a Thanos Sidecar component is configured here to upload this data to
GCS or S3 buckets.

You can also use `gsutil` (assuming Google Cloud Platform) and copy this data
directly into a GCS bucket.  Point a Thanos Store and Compact component at it
and the test data will become available to your Thanos Query component.  I've
used this to test Thanos for long data histories.

[1]: https://prometheus.io/
[2]: https://github.com/improbable-eng/thanos
