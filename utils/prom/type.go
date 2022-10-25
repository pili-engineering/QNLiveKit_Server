package prom

import (
	"hash/fnv"
)

type trackerValue struct {
	Labels    Labels
	TouchedAt int64 // ms
}

func hash(labelNames []string, labels Labels) uint64 {
	h := fnv.New64()
	for _, labelName := range labelNames {
		h.Write([]byte(labels[labelName]))
	}
	return h.Sum64()
}
