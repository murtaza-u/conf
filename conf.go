/*
Package conf provides a high-level abstraction over methods to manage
user configuration stored under $XDG_CONFIG_DIR. Being opinionated, all
configuration is written and read in YAML.
*/
package conf

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/murtaza-u/conf/internal/util"
	"github.com/rogpeppe/go-internal/lockedfile"
	"gopkg.in/yaml.v3"
)

// C encompasses various methods to read, write and query config file.
type C struct{}

// New creates a new instance of C. It is mandatory to run the Init()
// method immediately after.
func New() *C {
	return &C{}
}

// Init creates the necessary config directory and file (if absent). It
// is mandatory to call this method before performing any operations.
func (c *C) Init() error {
	path, err := c.getPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(path), 0751)
	if err != nil {
		return fmt.Errorf(
			"failed to softly create config directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to softly create config file: %w",
			err)
	}
	f.Close()

	return nil
}

// Read reads and marshals the config file into `out`. `out` must
// therefore be passed by reference.
func (c *C) Read(out any) error {
	path, err := c.getPath()
	if err != nil {
		return err
	}

	data, err := lockedfile.Read(path)
	if err != nil {
		return fmt.Errorf("failed to read %q: %w", path, err)
	}

	err = c.unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

// Write writes the configurations to the config path.
func (c *C) Write(in any) error {
	path, err := c.getPath()
	if err != nil {
		return err
	}

	out, err := c.marshal(in)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	err = lockedfile.Write(path, bytes.NewReader(out), 0600)
	if err != nil {
		return fmt.Errorf("failed to write to %q: %w", path, err)
	}

	return nil
}

// Query evaluates the yaml query and returns the result (in string).
func (c *C) Query(q string) (string, error) {
	path, err := c.getPath()
	if err != nil {
		return "", err
	}

	return util.EvaluateToString(q, path)
}

func (C) getApp() string {
	return filepath.Base(os.Args[0])
}

func (c *C) getPath() (string, error) {
	conf, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(
			"could not determine config directory: %w", err)
	}
	path := filepath.Join(conf, c.getApp(), "config.yaml")
	return path, nil
}

func (C) marshal(in any) ([]byte, error) {
	return yaml.Marshal(in)
}

func (C) unmarshal(in []byte, out any) error {
	return yaml.Unmarshal(in, out)
}
