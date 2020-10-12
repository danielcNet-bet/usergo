package core_test

import (
	"testing"
	"time"

	"github.com/thelotter-enterprise/usergo/core"
)

type innerStruct struct {
	ID int
}

type outerStruct struct {
	Name     string
	Inner    innerStruct
	Duration int64
}

const (
	tenSeconds time.Duration = time.Second * 10
)

func TestDecode(t *testing.T) {
	wantName := "guy kolbis"
	wantID := 555
	duration := tenSeconds.Milliseconds()
	input := map[string]interface{}{
		"Name":     wantName,
		"Duration": duration,
		"Inner": map[string]interface{}{
			"ID": wantID,
		},
	}

	var output outerStruct
	decoder := core.NewDecoder()
	err := decoder.Decode(input, &output)

	if err != nil {
		t.Error(err)
	}

	if output.Name != wantName {
		t.Error("Name does not match")
	}

	if output.Inner.ID != wantID {
		t.Error("ID does not match")
	}
}
