package vitenv

import (
	"bufio"
	"os"
)

func loadFileAndWriteEnvMap(filename string) (err error) {
	lines, err  := readFile(filename)
	if err != nil {
		return
	}
	return parseLinesAndWriteEnvMap(lines)
}

func readFile(filename string) (lines []string,  err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	bufScanner := bufio.NewScanner(file)
	if err = bufScanner.Err(); err != nil {
		return
	}

	for bufScanner.Scan() {
		lines = append(lines, bufScanner.Text())
	}
	return
}

