package id

import (
	"crypto/rand"
	"fmt"
)

// Generates pseudo-unique identifier.
//
// Generation is oversimplified because real uniqueness is not actually required.
// Collisions are not a problem as much as application goals are concerned.
func NewId() string {
	buf := make([]byte, 16)

	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Sprintf("Cannot read random bytes: %s", err))
	}

	return fmt.Sprintf("%x", buf)
}
