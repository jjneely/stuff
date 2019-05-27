package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func fixMetas(bucket, promenv string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	it := client.Bucket(bucket).Objects(ctx, nil)
	for objAttrs, err := it.Next(); err != iterator.Done; objAttrs, err = it.Next() {
		if err != nil {
			log.Printf("Error iterating over bucket gs://%s: %s", bucket, err)
			return err
		}
		if !strings.HasSuffix(objAttrs.Name, "/meta.json") {
			continue
		}

		log.Printf("Found: %s", objAttrs.Name)
		buf, err := getMeta(ctx, bucket, objAttrs.Name, client)
		if err != nil {
			log.Printf("Failed to get gs://%s/%s: %s",
				bucket, objAttrs.Name, err)
			continue
		}
		if writeLocalMeta(buf, objAttrs.Name) != nil {
			log.Printf("Failed to save a backup of gs://%s/%s: %s",
				bucket, objAttrs.Name, err)
			continue
		}
		buf, err = morphBuf(buf, promenv)
		if err != nil {
			log.Printf("Failed to alter meta.json: %s", err)
			continue
		}
		if buf != nil {
			log.Printf("Uploading modifed gs://%s/%s",
				bucket, objAttrs.Name)
			writeLocalMeta(buf, objAttrs.Name+".fixed")
			if err = setMeta(ctx, bucket, objAttrs.Name, buf, client); err != nil {
				log.Printf("Error writing to gs://%s/%s: %s",
					bucket, objAttrs.Name, err)
			}
		} else {
			log.Printf("Changes not needed for gs://%s/%s",
				bucket, objAttrs.Name)
			continue
		}
	}

	return nil
}

func getMeta(ctx context.Context, bucket, name string, client *storage.Client) ([]byte, error) {
	bkt := client.Bucket(bucket)
	obj := bkt.Object(name)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func setMeta(ctx context.Context, bucket, name string, buf []byte, client *storage.Client) error {
	bkt := client.Bucket(bucket)
	obj := bkt.Object(name)
	w := obj.NewWriter(ctx)

	_, err := w.Write(buf)
	if err != nil {
		w.Close()
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func writeLocalMeta(buf []byte, name string) error {
	fn := strings.ReplaceAll(name, "/", "-")
	log.Printf("Writing backup to: %s", fn)
	return ioutil.WriteFile(fn, buf, 0666)
}

func morphBuf(buf []byte, promenv string) ([]byte, error) {
	blob := make(map[string]interface{})
	err := json.Unmarshal(buf, &blob)
	if err != nil {
		return nil, err
	}

	thanos, ok := blob["thanos"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("meta.json does not contain a 'thanos' map")
	}
	labels, ok := thanos["labels"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("meta.json does not contain a 'thanos.labels' map")
	}
	_, ok = labels["promenv"]
	if ok {
		// No modifications needed
		return nil, nil
	}

	labels["promenv"] = promenv
	buf, err = json.MarshalIndent(blob, "", "\t")
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s [options] <GCS Bucket>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	confirm := flag.Bool("confirm", false, "Confirm you want to make changes to GCS")
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
	}

	bucket := flag.Arg(0)
	// Sanitize the bucket string for uniformity
	if strings.HasPrefix(bucket, "gs://") {
		bucket = bucket[len("gs://"):]
	}
	for bucket[len(bucket)-1] == '/' {
		bucket = bucket[:len(bucket)-1]
	}
	if strings.HasPrefix(bucket, "bruce-thanos-lts-expiring-backup") {
		log.Fatalf("Refusing to alter Thanos LTS backup bucket gs://%s", bucket)
	}

	// Regex out the prometheus environment string
	envexp := regexp.MustCompile(`^[gs://]*[a-zA-Z0-9]+-thanos-lts-[-a-zA-Z0-9]+-([a-z]+)/*$`)
	if !envexp.MatchString(bucket) {
		log.Fatalf("Refusing to alter GCE bucket %s does not match expected regexp", bucket)
	}
	promenv := envexp.ReplaceAllString(bucket, "$1")
	if promenv == "backup" {
		log.Fatalf("Refusing to alter Thanos LTS backup bucket gs://%s", bucket)
	}

	// Report what we found
	log.Printf("GCS Bucket   : %s", bucket)
	log.Printf("Promenv Value: %s", promenv)
	if !*confirm {
		log.Printf("Making no changes with -confirm")
		os.Exit(0)
	}

	// Search and destroy
	err := fixMetas(bucket, promenv)
	if err != nil {
		log.Printf("Error: %s", err)
	}
}
