package config

import (
	"math/rand"
	"time"
)

func InitRandomSeed() {
	rand.Seed(time.Now().UnixNano())
}
