package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"log"
)

func main() {
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatalln(err)
	}
	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	fmt.Println(token)
	hash := sha256.Sum256([]byte(token))
	fmt.Println(hash)
}
