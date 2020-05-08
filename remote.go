package vitenv

import (
	"bufio"
	"errors"
	"net/http"
	"strings"
	"time"
)

var retryСount int

func loadRemoteFileAdnWriteEnvMap(filename string) (lines []string, err error) {
	retryСount = 5
	lines, err = readRemoteFile(filename)
	return
}

func readRemoteFile(filename string) (lines []string, err error) {
	file := strings.Split(filename, "@")
	req, err := http.NewRequest("GET", file[0], nil)
	if err != nil {
		return
	}

	if len(file) == 2 {
		req.Header.Add("Authorization", "Basic " + file[1])
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		if retryСount > 5 {
			retryСount--
			time.Sleep(time.Second)
			return readRemoteFile(filename)
		}
		return
	}
	defer resp.Body.Close()
	if 200 != resp.StatusCode {
		err = errors.New(resp.Status)
		return
	}
	bufScanner := bufio.NewScanner(resp.Body)
	if err = bufScanner.Err(); err != nil {
		return
	}

	for bufScanner.Scan() {
		lines = append(lines, bufScanner.Text())
	}
	return
}
