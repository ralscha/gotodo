package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) isPasswordCompromised(password string) (bool, error) {
	alg := sha1.New()
	alg.Write([]byte(password))
	hash := strings.ToUpper(hex.EncodeToString(alg.Sum(nil)))
	prefix := strings.ToUpper(hash[:5])
	suffix := strings.ToUpper(hash[5:])

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.pwnedpasswords.com/range/"+prefix, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Add-Padding", "true")

	response, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected http status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	hashList := strings.Split(string(body), "\r\n")

	for _, line := range hashList {
		if line[:35] == suffix {
			occurrence, err := strconv.ParseInt(line[36:], 10, 64)
			if err != nil {
				return false, err
			}
			if occurrence > 0 {
				return true, nil
			}
			break
		}
	}

	return false, nil
}
