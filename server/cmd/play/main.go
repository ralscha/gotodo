package main

import (
	"fmt"
	"github.com/alexedwards/argon2id"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	// CreateHash returns a Argon2id hash of a plain-text password using the
	// provided algorithm parameters. The returned hash follows the format used
	// by the Argon2 reference C implementation and looks like this:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	params := &argon2id.Params{
		Memory:      1 << 17,
		Iterations:  20,
		Parallelism: 8,
		SaltLength:  16,
		KeyLength:   32,
	}

	hash, err := argon2id.CreateHash("pa$$word", params)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	fmt.Println(hash)
	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	start := time.Now()
	match, err := argon2id.ComparePasswordAndHash("pa$$word", hash)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	elapsed := time.Since(start)
	log.Printf("compare took %s", elapsed)

	log.Printf("Match: %v", match)
}
