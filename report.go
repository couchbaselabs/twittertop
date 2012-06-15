package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/dustin/gomemcached/client"
)

func report(prefix string) {
	log.Printf("Reporting on %v", prefix)

	client, err := memcached.Connect("tcp", *mcServer)
	if err != nil {
		log.Printf("Error connecting to memcached:  %v", err)
		return
	}
	defer client.Close()

	resp, err := client.Get(0, prefix+listSuffix)
	if err != nil {
		log.Printf("Error reporting on %s:  %v", prefix, err)
		return
	}

	topKeys := ranks{}
	err = json.Unmarshal(resp.Body, &topKeys)
	if err != nil {
		log.Printf("Error unmarshaling %s", resp.Body)
		return
	}
	sort.Sort(&topKeys)

	tw := tabwriter.NewWriter(os.Stdout, 8, 8, 2, '\t', 0)
	for i, k := range topKeys.keys {
		fmt.Fprintf(tw, "   %v\t%v\n", k, topKeys.counts[k])
		if i > 20 {
			break
		}
	}
	tw.Flush()
}
