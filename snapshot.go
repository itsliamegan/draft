package main

import (
	"bytes"
	"crypto/sha256"
	"io"
	"mime"
	"os"
	"path/filepath"
)

type Snapshot struct {
	dir    string
	states map[string]signature
}

type Change struct {
	Filename string `json:"filename"`
	MimeType string `json:"mimeType"`
	Contents string `json:"contents"`
}

type signature []byte

func NewSnapshot(dir string) (*Snapshot, error) {
	states := make(map[string]signature)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if !isRelevant(path) {
			continue
		}

		sig, err := hash(path)
		if err != nil {
			return nil, err
		}

		states[path] = sig
	}

	return &Snapshot{dir: dir, states: states}, nil
}

func (old *Snapshot) Diff(new *Snapshot) ([]Change, error) {
	changes := make([]Change, 0)

	for file, sig := range old.states {
		if !sig.equal(new.states[file]) {
			b, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}
			contents := string(b)

			basename := filepath.Base(file)
			ext := filepath.Ext(file)
			mimeType := mime.TypeByExtension(ext)

			changes = append(changes, Change{Filename: basename, MimeType: mimeType, Contents: contents})
		}
	}

	return changes, nil
}

func (sig signature) equal(other signature) bool {
	return bytes.Equal(sig, other)
}

func hash(file string) (s signature, err error) {
	h := sha256.New()

	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(h, fp)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func isRelevant(path string) bool {
	if filepath.Base(path)[0:1] == "." {
		return false
	}

	if stat, _ := os.Stat(path); stat.IsDir() {
		return false
	}

	return true
}
