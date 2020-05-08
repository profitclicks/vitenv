package vitenv

import (
	"fmt"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	if err := Load(
		 "fixtures/equals.env",
		); err != nil {
		t.Error(err.Error())
	}
	fmt.Println(envMap)
}