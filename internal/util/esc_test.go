package util_test

import (
	"testing"

	"github.com/murtaza-u/conf/internal/util"
)

func TestEscReturns(t *testing.T) {
	in := "hello\nworld\ragain"
	exp := "hello\\nworld\\ragain"
	out := util.EscReturns(in)
	if out != exp {
		t.Fatalf("Expected: %s. Got: %s", exp, out)
	}
}

func TestUnEscReturns(t *testing.T) {
	in := "hello\\nworld\\ragain"
	exp := "hello\nworld\ragain"
	out := util.UnEscReturns(in)
	if out != exp {
		t.Fatalf("Expected: %s. Got: %s", exp, out)
	}
}
