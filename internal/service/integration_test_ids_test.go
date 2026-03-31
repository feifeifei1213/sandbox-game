package service

import (
	"sync/atomic"
	"time"
)

var integrationTestCounter atomic.Int64

func nextIntegrationUniqueSeed() int64 {
	return time.Now().UnixNano() + integrationTestCounter.Add(1)
}

func integrationGroupNoFromSeed(seed int64) int {
	if seed < 0 {
		seed = -seed
	}
	return int(seed%900000000) + 100000000
}
