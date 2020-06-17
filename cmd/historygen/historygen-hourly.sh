#!/usr/bin/env bash

echo "A script to generate Prometheus time series data for every hour across the last N hours"

if [[ $# -ne 2 ]]; then
  echo "$USAGE $0 <total-series> <total-hours>"
  exit 1
fi

total_series=$1
total_hours=$2
series_batch=$[total_series/total_hours]

rm -fr data
for i in `seq $[total_hours-1] 0`; do
  start=$[i+1]
  echo "Processing hour -${start}h to -${i}h..."
  ./historygen \
    -C $total_series \
    -c $series_batch \
    -d 1h \
    -b 1h \
    -n $[i*series_batch] -s ${i}h
done
