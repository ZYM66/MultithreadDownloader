package statistic

import (
	"fmt"
	"time"
)

func TimeCost(start time.Time) {
	tc := time.Since(start)
	fmt.Printf("time cost = %v\n", tc)
}
