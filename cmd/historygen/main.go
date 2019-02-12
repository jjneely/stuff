package main

import (
	"flag"
	"log"
	"time"

	"github.com/improbable-eng/thanos/benchmark/pkg/tsdb"
)

var (
	duration = flag.Duration("d", time.Hour*720,
		"Time duration of historical data to generate")
	outDir = flag.String("o", "data/",
		"Output directory to generate TSDB blocks in")
	timeSeries = flag.Int("c", 1,
		"Number of time series to generate")
	sampleInterval = flag.Duration("i", time.Second*15,
		"Duration between samples")
	blockLength = flag.Duration("b", time.Hour*2,
		"TSDB block length")
)

func main() {
	log.Printf("Generate Prometheus TSDB test data.")
	flag.Parse()

	endTime := time.Now()
	err := tsdb.CreateThanosTSDB(tsdb.Opts{
		OutputDir:      *outDir,
		NumTimeseries:  *timeSeries,
		StartTime:      endTime.Add(-*duration),
		EndTime:        endTime,
		SampleInterval: *sampleInterval,
		BlockLength:    *blockLength,
	})

	if err != nil {
		log.Fatalf("Error generating data: %s", err)
	}

	log.Printf("TSDB data generation complete")
}
