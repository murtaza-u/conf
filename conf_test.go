package conf_test

import (
	"testing"

	"github.com/murtaza-u/conf"
)

type myCfg struct {
	Foo  string `yaml:"foo"`
	Bar  bool   `yaml:"bar"`
	Blah int    `yaml:"blah"`
}

func TestInit(t *testing.T) {
	conf := conf.New()
	err := conf.Init()
	if err != nil {
		t.Fatal("Expected: err == nil. Got: err =", err)
	}
}

func TestWrite(t *testing.T) {
	conf := conf.New()
	err := conf.Init()
	if err != nil {
		t.Fatal("failed to initialize conf:", err)
	}

	cfg := &myCfg{
		Foo:  "foo",
		Bar:  true,
		Blah: 100,
	}
	err = conf.Write(cfg)
	if err != nil {
		t.Fatalf(`failed to write %v to config file.
			Expected: err == nil. Got: err = %s`, cfg, err)
	}
}

func TestRead(t *testing.T) {
	conf := conf.New()
	err := conf.Init()
	if err != nil {
		t.Fatal("failed to initialize conf:", err)
	}

	cfg := new(myCfg)
	err = conf.Read(cfg)
	if err != nil {
		t.Fatalf(`failed to read from config file.
			Expected: err == nil. Got err = %s`, err)
	}
	if cfg.Foo != "foo" || !cfg.Bar || cfg.Blah != 100 {
		t.Fatalf("incorrect data returned from read operation")
	}
}

func TestQuery(t *testing.T) {
	conf := conf.New()
	err := conf.Init()
	if err != nil {
		t.Fatal("failed to initialize conf:", err)
	}

	foo, err := conf.Query(".foo")
	if err != nil {
		t.Fatal("failed to query `.foo`. Expected err = nil. Got err =",
			err)
	}
	if foo != "foo" {
		t.Fatalf(`incorrect data returned from query.
			Expected %s. Got %s`, "foo", foo)
	}
}
