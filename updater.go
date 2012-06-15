package main

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func removeBottom(r *ranks) {
	if r.Len() > max_entries {
		sort.Sort(r)
		for _, k := range r.keys[max_entries:] {
			delete(r.counts, k)
		}
	}
}

func keyMangle(k string) string {
	return k
}

type ranks struct {
	keys   []string
	counts map[string]uint64
}

func (r ranks) Len() int {
	return len(r.keys)
}

func (r ranks) Less(i, j int) bool {
	return r.counts[r.keys[j]] < r.counts[r.keys[i]]
}

func (r *ranks) Swap(i, j int) {
	r.keys[i], r.keys[j] = r.keys[j], r.keys[i]
}

func (r *ranks) updateKeys() {
	r.keys = make([]string, 0, len(r.counts))
	for k := range r.counts {
		r.keys = append(r.keys, k)
	}
}

func (r ranks) MarshalJSON() ([]byte, error) {
	m := map[string]uint64{}
	for k, v := range r.counts {
		m[k] = v
	}
	return json.Marshal(m)
}

func (r *ranks) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &r.counts)
	if err == nil {
		r.updateKeys()
	}
	return err
}

func makeRanks(candidates map[string]uint64) ranks {
	rv := ranks{counts: candidates}
	rv.updateKeys()
	return rv
}

func updateList(client persister, prefix string, totals map[string]uint64) {
	_, err := client.CAS(lPrefix+prefix,
		func(oldBody []byte) []byte {
			if len(oldBody) == 0 {
				b, err := json.Marshal(totals)
				if err != nil {
					log.Fatalf("Error marshaling keys: %v", err)
				}
				return b
			}
			topKeys := map[string]uint64{}
			err := json.Unmarshal(oldBody, &topKeys)
			if err != nil {
				log.Fatalf("Error unmarshaling %s", oldBody)
			}

			// Just shove in our keys to dedup them
			for k, v := range totals {
				topKeys[k] = v
			}

			keys := make([]string, 0, len(topKeys))
			for k := range topKeys {
				keys = append(keys, cPrefix+prefix+k)
			}

			m, err := client.GetBulk(keys)
			if err != nil {
				log.Fatalf("Error getting %v: %v", keys, err)
			}

			candidates := map[string]uint64{}
			for k, resp := range m {
				v, err := strconv.ParseUint(string(resp.Body), 10, 64)
				if err != nil {
					log.Fatalf("Error parsing %#v: %v", resp, err)
				}
				candidates[k[len(cPrefix+prefix):]] = v
			}

			r := makeRanks(candidates)

			removeBottom(&r)

			b, err := json.Marshal(r)
			if err != nil {
				log.Fatalf("Error marshaling %v: %v", m, err)
			}
			return b
		}, int(windowSize.Seconds()))
	if err != nil {
		log.Panicf("Error CASing in a new list: %v", err)
	}
}

func process(client persister, prefix string, t *Tweet) {
	parts := strings.FieldsFunc(strings.ToLower(t.Text), func(r rune) bool {
		return !(unicode.IsLetter(r) || r == '\'')
	})
	counts := map[string]uint64{}
	for _, p := range parts {
		if len(p) < max_word_len {
			counts[p] = counts[p] + 1
		}
	}
	totals := map[string]uint64{}
	for w, count := range counts {
		k := cPrefix + prefix + w
		var err error
		totals[w], err = client.Incr(k, count, 1,
			10*int(windowSize.Seconds()))
		if err != nil {
			log.Fatalf("Error incrementing %v: %v", k, err)
		}
		if *foreverWords {
			_, err = client.Incr(cPrefix+w, count, 1, 0)
		}
	}
	updateList(client, prefix, totals)
}
