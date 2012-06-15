package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
)

func report(prefix string) {
	log.Printf("Reporting on %v", prefix)

	client, err := getPersister(*mcServer)
	if err != nil {
		log.Printf("Error connecting to persister:  %v", err)
		return
	}
	defer client.Close()

	resp, err := client.Get(lPrefix + prefix)
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
