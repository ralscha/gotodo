package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	password := "skjfksjfkfaasdf2323"

	alg := sha1.New()
	alg.Write([]byte(password))
	hash := strings.ToUpper(hex.EncodeToString(alg.Sum(nil)))
	prefix := strings.ToUpper(hash[:5])
	suffix := strings.ToUpper(hash[5:])

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.pwnedpasswords.com/range/"+prefix, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Add-Padding", "true")

	response, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if rerr := response.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if response.StatusCode != http.StatusOK {
		log.Fatalln(response.StatusCode)
	}

	// Parse our resp.Body.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	hashList := strings.Split(string(body), "\r\n")
	found := false
	for _, line := range hashList {
		if line[:35] == suffix {
			occurrence, err := strconv.ParseInt(line[36:], 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			if occurrence > 0 {
				fmt.Println("found it:" + strconv.FormatInt(occurrence, 10))
				found = true
			}
			break
		}
	}
	if !found {
		fmt.Println("not in HIBP")
	}

}
