package utility

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func GetHashForObject(o any) string {
	b, _ := json.Marshal(o)
	sum := sha256.Sum256(b)
	return fmt.Sprintf("%x", sum[:])
}
