package utilities

import (
	"math/rand"
	"time"
)

func RandIntInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	n := min + rand.Intn(max-min+1)
	return n
}
