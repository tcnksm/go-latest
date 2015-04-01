package latest

import (
	"testing"
)

func TestJSON_implement(t *testing.T) {
	var _ Source = &JSON{}
}
