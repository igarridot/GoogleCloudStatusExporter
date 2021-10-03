package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestCompareMetrics(t *testing.T) {
	err := testutil.GatherAndCompare()
}
