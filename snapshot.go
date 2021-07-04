package main

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
)

type Snapshot struct {
	dir    string
	states map[string]signature
}

type Change struct {
	File     string   `json:"file"`
	Type     FileType `json:"type"`
	Contents string   `json:"contents"`
}

type FileType string

const (
	DocumentFileType   FileType = "document"
	StylesheetFileType          = "stylesheet"
)

type signature []byte

func NewSnapshot(dir string) (*Snapshot, error) {
	states := make(map[string]signature)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		file := filepath.Join(dir, entry.Name())

		if !isRelevant(file) {
			continue
		}

		sig, err := hash(file)
		if err != nil {
			return nil, err
		}

		states[file] = sig
	}

	return &Snapshot{dir: dir, states: states}, nil
}

func (old *Snapshot) Diff(new *Snapshot) ([]Change, error) {
	changes := make([]Change, 0)

	for file, sig := range old.states {
		if !bytes.Equal(sig, new.states[file]) {
			b, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}
			contents := string(b)

			basename := filepath.Base(file)

			change := Change{File: basename, Type: typeOf(file), Contents: contents}
			changes = append(changes, change)
		}
	}

	return changes, nil
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

func typeOf(file string) FileType {
	switch filepath.Ext(file) {
	case ".html":
		return DocumentFileType
	case ".css":
		return StylesheetFileType
	default:
		return ""
	}
}

func isRelevant(file string) bool {
	if filepath.Base(file)[0:1] == "." {
		return false
	}

	return true
}
