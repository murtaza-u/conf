package vars_test

import (
	"testing"

	"github.com/murtaza-u/conf/vars"
)

func TestInit(t *testing.T) {
	vars := vars.New()
	err := vars.Init()
	if err != nil {
		t.Fatal("Expected: err == nil. Got: err =", err)
	}
}

func TestSet(t *testing.T) {
	vars := vars.New()
	err := vars.Init()
	if err != nil {
		t.Fatal("failed to initialize vars:", err)
	}

	err = vars.Set("foo", "bar")
	if err != nil {
		t.Fatal("failed to set `foo = bar`. Expected: err == nil. Got: err =",
			err)
	}
}

func TestExists(t *testing.T) {
	vars := vars.New()
	err := vars.Init()
	if err != nil {
		t.Fatal("failed to initialize vars:", err)
	}

	ok := vars.Exists("foo")
	if !ok {
		t.Fatal("key `foo` exists but -ve was returned")
	}

	ok = vars.Exists("bar")
	if ok {
		t.Fatal("key `bar` does not exists but +ve was returned")
	}
}

func TestGet(t *testing.T) {
	vars := vars.New()
	err := vars.Init()
	if err != nil {
		t.Fatal("failed to initialize vars:", err)
	}

	v := vars.Get("foo")
	if v != "bar" {
		t.Fatalf("Expected: %q. Got: %q", "bar", v)
	}

	v = vars.Get("bar")
	if v != "" {
		t.Fatalf("Expected: %s. Got %q", "", v)
	}
}

func TestDel(t *testing.T) {
	vars := vars.New()
	err := vars.Init()
	if err != nil {
		t.Fatal("failed to initialize vars:", err)
	}

	err = vars.Del("foo")
	if err != nil {
		t.Fatal("Expected: err = nil. Got err =", err)
	}

	v := vars.Get("foo")
	if v == "bar" {
		t.Fatalf("key `foo` not deleted. Expected: %s. Got %q", "", v)
	}
}
