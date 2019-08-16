package fs

import (
	"testing"
)

func TestFake(t *testing.T) {
	ListDir("./")
	t.Error("Noo")
}
