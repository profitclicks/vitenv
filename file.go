package vitenv

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var lock sync.RWMutex

func GetEnv(name string, def string) string {
	lock.RLock()
	defer lock.RUnlock()
	if value, ok := envMap[name]; ok {
		return value
	}
	return def
}

func parseLinesAndWriteEnvMap(lines []string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	for i, line := range lines {
		if isIgnoredLine(line) {
			continue
		}
		var key, value string

		_, key, value, err = parseLine(line)
		if err != nil {
			return errors.New(fmt.Sprintf("line %d - %s", i+1, err.Error()))
		}
		if len(key) == 0 {
			continue
		}
		envMap[key] = value
	}
	return
}
func parseLine(line string) (export bool, key string, value string, err error) {
	items := strings.Split(line, `=`)
	if len(items) == 1 {
		err = errors.New("invalid line")
		return
	}
	var isEmpty bool
	if value, isEmpty = parseValue(strings.Join(items[1:], `=`)); isEmpty {
		return
	}
	key = strings.TrimSpace(items[0])
	if strings.HasPrefix(key, "export") {
		export = true
		key = strings.TrimSpace(strings.TrimPrefix(key, "export"))
	}
	return
}
func parseValue(value string) (result string, isEmpty bool) {
	value = strings.TrimSpace(value)
	var isQuotesO bool
	var isQuotesW bool
	var isSkip bool
	var isBreak bool

	endSlice := 0

	var re = regexp.MustCompile(`(?m)\${([^}]*)}`)
	isReg := false
	value = re.ReplaceAllStringFunc(value, func(s string) string {
		isReg = true
		if value, ok := envMap[s[2:len(s)-1]]; ok {
			return value
		}
		return s
	})
	if isReg {
		return value, false
	}

	var trimQuotes string
	for i := 0; i < len(value); i++ {
		if !isSkip {
			switch value[i] {
			case '\'':
				if isQuotesO {
					isQuotesO = false
				} else if i == 0 {
					trimQuotes = `'`
					isQuotesO = true
				}

			case '"':
				if isQuotesW {
					isQuotesW = false
				} else if i == 0 {
					trimQuotes = `"`
					isQuotesW = true
				}
			case '\\':
				isSkip = true
			case '#':
				if !isQuotesO && !isQuotesW {
					isBreak = true
				}
			}
		}
		if isBreak {
			endSlice = i
			break
		}
		endSlice = i + 1
	}
	if endSlice > 0 {
		result = value[:endSlice]

		result = strings.Trim(result, fmt.Sprintf(" #\n\t%s", trimQuotes))
	} else {
		isEmpty = true
	}

	return
}

func isIgnoredLine(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#")
}
