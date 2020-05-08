package vitenv

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type SyncError struct {
	sync.RWMutex
	description string
	code        int
}

func NewError(code int, description string) SyncError {
	return SyncError{
		description: description,
		code:        code,
	}
}

func (se *SyncError) SetError(code int, description string) int {
	se.Lock()
	defer se.Unlock()

	se.code = code
	se.description = description

	return se.code
}

func (se *SyncError) AppendError(description ... string) int {
	se.Lock()
	defer se.Unlock()

	se.description = strings.Join([]string{se.description, strings.Join(description, ` `)}, "\n")

	return se.code
}


func (se *SyncError)GetCode() int {
	se.RLock()
	defer se.RUnlock()

	return se.code
}

func (se *SyncError)IsError() bool {
	se.RLock()
	defer se.RUnlock()

	return len(se.description) != 0
}

func (se *SyncError) GetDescription() string {
	se.RLock()
	defer se.RUnlock()

	return se.description
}

func (se *SyncError) GetError() error {
	se.RLock()
	defer se.RUnlock()

	return errors.New(fmt.Sprintf("code:%d - %s", se.code, se.description))
}

