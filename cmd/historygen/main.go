package main

import (
	"flag"
	"log"
	"time"

	"github.com/qcrisw/stuff/tsdb"
)

var (
	duration = flag.Duration("d", time.Hour*720,
		"Time duration of historical data to generate")
	timeShift = flag.Duration("s", time.Hour*0,
		"Time shift of historical data to generate")
	outDir = flag.String("o", "data/",
		"Output directory to generate TSDB blocks in")
	timeSeries = flag.Int("c", 1,
		"Number of time series to generate")
	totalTimeSeries = flag.Int("C", 0,
		"Total number of time series to generate using multiple invocations (only needed for the zero padding of instance names)")
	timeSeriesStartIndex = flag.Int("n", 0,
		"Start index for time series instance names")
	sampleInterval = flag.Duration("i", time.Second*15,
		"Duration between samples")
	blockLength = flag.Duration("b", time.Hour*2,
		"TSDB block length")
)

func main() {
	log.Printf("Generate Prometheus TSDB test data.")
	flag.Parse()

	endTime := time.Now().Add(-*timeShift)
	err := tsdb.CreateThanosTSDB(tsdb.Opts{
		OutputDir:            *outDir,
		NumTimeseries:        *timeSeries,
		TotalNumTimeSeries:   *totalTimeSeries,
		TimeseriesStartIndex: *timeSeriesStartIndex,
		StartTime:            endTime.Add(-*duration),
		EndTime:              endTime,
		SampleInterval:       *sampleInterval,
		BlockLength:          *blockLength,
	})

	if err != nil {
		log.Fatalf("Error generating data: %s", err)
	}

	log.Printf("TSDB data generation complete")
}
