package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) isPasswordCompromised(ctx context.Context, password string) (bool, error) {
	alg := sha1.New()
	alg.Write([]byte(password))
	hash := strings.ToUpper(hex.EncodeToString(alg.Sum(nil)))
	prefix := strings.ToUpper(hash[:5])
	suffix := strings.ToUpper(hash[5:])

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.pwnedpasswords.com/range/"+prefix, nil)
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

	hashList := strings.SplitSeq(string(body), "\r\n")

	for line := range hashList {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		if strings.TrimSpace(parts[0]) == suffix {
			occurrence, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
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
