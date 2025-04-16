package pkg

import "math/rand"

func GenerateUUID(name string) int64 {
	// TODO: Generate a random number based on the name
	return rand.Int63()
}

func GeneratePlayerUUID() int64 {
	return rand.Int63()
}
