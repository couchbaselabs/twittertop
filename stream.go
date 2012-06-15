package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dustin/gomemcached/client"
)

type Tweet struct {
	Sender struct {
		User string `json:"screen_name"`
		Name string
	} `json:"user"`
	Text string
}

const qtimeFormat = "150405:"
const ptimeFormat = "20060102T150405"
const max_entries = 1000
const max_word_len = 32

const listSuffix = "-list"

var windowSize = flag.Duration("interval", time.Second*30, "Reporting interval")
var numWorkers = flag.Int("workers", 8, "Number of workers")
var recordTo = flag.String("record", "", "Record the stream to a file")
var multiplier = flag.Int("multiplier", 1, "Tweet multiplier")
var mcServer = flag.String("memcached", "localhost:11211",
	"host:port of your memcached server")

func parseNext(d *json.Decoder) (rv Tweet, err error) {
	err = d.Decode(&rv)
	return
}

func handle(ch <-chan Tweet, pch <-chan string) {
	client, err := memcached.Connect("tcp", *mcServer)
	if err != nil {
		log.Fatalf("Error connecting to memcached:  %v", err)
	}
	defer client.Close()
	prefix := ""

	for {
		select {
		case prefix = <-pch:
		case t := <-ch:
			if t.Text != "" {
				process(client, prefix, &t)
			}

		}
	}
}

func moveTime(listeners []chan string) {
	var currentPrefix string

	for {
		t := time.Now().Unix()
		t = t - (t % int64(windowSize.Seconds()))

		oldPrefix := currentPrefix
		currentWindow := time.Unix(t, 0)
		currentPrefix = currentWindow.Format(qtimeFormat)

		for _, ch := range listeners {
			ch <- currentPrefix
		}

		if oldPrefix != "" {
			report(oldPrefix)
		}

		next := time.Unix(t, 0).Add(*windowSize)
		toSleep := next.Sub(time.Now())
		time.Sleep(toSleep)
	}
}

func openStream(path string) (io.ReadCloser, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept-Encoding", "gzip")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}

		return gz, nil
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	panic("Unreachable")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] src\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nsrc can be a path to a local file, or a URL, e.g.\n")
	fmt.Fprintf(os.Stderr, "  https://user:pass@stream.twitter.com/1/statuses/sample.json\n")
	os.Exit(1)
}

func streamTo(path string, ch chan<- Tweet) {
	stream, err := openStream(path)
	if err != nil {
		log.Printf("Error opening stream: %v", err)
		return
	}
	defer stream.Close()
	var r io.Reader = stream
	if *recordTo != "" {
		f, err := os.OpenFile(*recordTo,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Printf("Error opening output stream: %v", err)
			return
		}
		defer f.Close()
		r = io.TeeReader(stream, f)
	}

	d := json.NewDecoder(r)

	for {
		tweet, err := parseNext(d)
		if err != nil {
			log.Printf("Error in stream: %v", err)
			return
		}

		for i := 0; i < *multiplier; i++ {
			ch <- tweet
		}
	}
}

func main() {
	log.SetFlags(log.Lmicroseconds)

	flag.Usage = usage

	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	listenerChans := make([]chan string, *numWorkers)
	for i := 0; i < *numWorkers; i++ {
		listenerChans[i] = make(chan string)
	}

	go moveTime(listenerChans)

	ch := make(chan Tweet, 1000)

	for _, lch := range listenerChans {
		go handle(ch, lch)
	}

	for {
		start := time.Now()
		streamTo(flag.Arg(0), ch)
		if time.Since(start).Seconds() < 5 {
			log.Printf("Crashed in %v, slowing down",
				time.Since(start))
			time.Sleep(time.Second * 5)
		}
		log.Printf("Restarting stream")
	}
}
