package statistic

import (
	"fmt"
	"time"
)

func timeCost(start time.Time) {
	tc := time.Since(start)
	fmt.Printf("time cost = %v\n", tc)
}
