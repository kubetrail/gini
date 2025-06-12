package flags

import (
	"fmt"
	"testing"
)

func TestModels(t *testing.T) {
	for i, model := range Models {
		fmt.Printf("%d: %s\n", i, model)
	}
}
