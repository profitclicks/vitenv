package vitenv

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

var files []string
var remoteFiles []string
var envMap map[string]string

func Load(filenames ...string) (err error) {
	for _, file := range filenames {
		if strings.Contains(file, "http") {
			remoteFiles = append(remoteFiles, file)
			continue
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		files = []string{".env"}
	}
	return load()
}

func OnUpdate() (err error) {
	return load()
}



func load() (err error) {
	envMap = make(map[string]string)

	for _, file := range files {
		if err := loadFileAndWriteEnvMap(file); err != nil {
			return errors.New(fmt.Sprintf("file %s %s", file, err.Error()))
		}
	}
	if _ , ok := envMap["REMOTE_ENV"]; len(remoteFiles) == 0 && ok {
		remoteFiles = append(remoteFiles, GetEnv("REMOTE_ENV", ""))
	}


	var wg sync.WaitGroup
	var se SyncError

	for _, file := range remoteFiles {
		wg.Add(1)

		go func(file string) {
			defer wg.Done()

			errFormat := "file %s %s"

			lines, err := loadRemoteFileAdnWriteEnvMap(file)
			if err != nil{
				se.AppendError(fmt.Sprintf(errFormat, file, err.Error()))
				return
			}

			if err := parseLinesAndWriteEnvMap(lines); err != nil {
				se.AppendError(fmt.Sprintf(errFormat, file, err.Error()))
			}

		}(file)
	}
	wg.Wait()
	if se.IsError(){
		return errors.New(se.GetDescription())
	}
	for key, value := range envMap {
		//fmt.Printf("%s=%s\n", key, value)
		if err = os.Setenv(key, value); err != nil {
			return
		}
	}
	return
}
