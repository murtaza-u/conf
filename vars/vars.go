package vars

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/murtaza-u/conf/internal/util"

	"github.com/rogpeppe/go-internal/lockedfile"
)

// Map encompasses an internal concurrency-safe map object as well as
// different operations to be performed on it.
type Map struct {
	mu sync.Mutex
	m  map[string]string

	path string
}

// New creates a new map type internally. It is mandatory to run the
// Init() method immediately after.
func New() *Map {
	m := &Map{
		m:  make(map[string]string),
		mu: sync.Mutex{},
	}
	return m
}

// Init creates the necessary cache directory and file (if absent) and
// loads the keys and values into the internal map object. It is
// mandatory to call this method before performing any operations.
func (m *Map) Init() error {
	cache, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("could not determine cache directory: %w",
			err)
	}

	err = os.MkdirAll(cache, 0751)
	if err != nil {
		return fmt.Errorf("failed to softly create %q: %w", cache, err)
	}

	m.path = m.getPath(cache)
	f, err := os.OpenFile(m.path, os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to softly create %q: %w", m.path, err)
	}
	f.Close()

	err = m.read()
	if err != nil {
		return err
	}

	return nil
}

// Get returns the value associated with the given key. If the key does
// not exist, it returns a blank string ("").
func (m *Map) Get(key string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if v, ok := m.m[key]; ok {
		return v
	}

	return ""
}

// Exists checks if a key is present.
func (m *Map) Exists(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.m[key]
	return ok
}

// Del deletes a key. If the key does not exist, no operation is
// performed.
func (m *Map) Del(key string) error {
	m.mu.Lock()
	delete(m.m, key)
	m.mu.Unlock()

	return m.write()
}

// Set adds/updates a key-value pair.
func (m *Map) Set(key, value string) error {
	m.mu.Lock()
	m.m[key] = value
	m.mu.Unlock()

	return m.write()
}

func (*Map) getPath(cacheDir string) string {
	app := filepath.Base(os.Args[0])
	return filepath.Join(cacheDir, app)
}

func (m *Map) read() error {
	if m.path == "" {
		return fmt.Errorf("vars not yet initialised")
	}

	data, err := lockedfile.Read(m.path)
	if err != nil {
		return fmt.Errorf("failed to read %q: %w", m.path, err)
	}

	err = m.unmarshalText(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal text: %w", err)
	}

	return nil
}

func (m *Map) write() error {
	out := m.marshalText()
	err := lockedfile.Write(m.path, bytes.NewReader(out), 0600)
	if err != nil {
		return fmt.Errorf("failed to write to %q: %w", m.path, err)
	}
	return nil
}

func (m *Map) marshalText() []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	var out string
	for k, v := range m.m {
		out += fmt.Sprintf("%s=%s\n", k, util.EscReturns(v))
	}

	return []byte(out)
}

func (m *Map) unmarshalText(in []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sc := bufio.NewScanner(bytes.NewReader(in))
	for sc.Scan() {
		txt := sc.Text()
		splits := strings.SplitN(txt, "=", 2)
		if len(splits) == 2 {
			k := splits[0]
			v := util.UnEscReturns(splits[1])
			m.m[k] = v
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("failed to parse file %q: %w", m.path, err)
	}

	return nil
}
