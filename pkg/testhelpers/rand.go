package testhelpers

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func init() {
	randSeed := time.Now().UnixNano()
	if seed := os.Getenv("TEST_SEED"); seed != "" {
		log.Printf("Using TEST_SEED as random seed")
		i, err := strconv.ParseInt(seed, 10, 64)
		if err != nil {
			log.Println("error parsing TEST_SEED", err)
		} else {
			randSeed = i
		}
	}
	rand.Seed(randSeed)
	log.Printf("Random Seed is %d", randSeed)
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

//GenerateRandomEmail return a random email suitable for testing
func GenerateRandomEmail() string {
	return RandString(16) + "@example.com"
}

// RandString is a convenience function that returns the string of RandBytes(n)
func RandString(n int) string {
	return string(RandBytes(n))
}

// RandBytes generates a random []byte of length n
func RandBytes(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; {
		// #nosec
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return b
}
