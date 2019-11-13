package app

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/sethgrid/pester"
)

func createReq(url string) ([]byte, error) {
	client := pester.New()
	client.Concurrency = 1
	client.MaxRetries = 10
	client.Backoff = pester.ExponentialBackoff
	client.KeepLog = true
	client.Timeout = 10 * time.Second

	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.New(client.LogString())
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
