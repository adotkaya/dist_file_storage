package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{PathTransformFunc: CASPathTransformFunc}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "bestpicture"
	pathKey := CASPathTransformFunc(key)
	expectedFilename := "71056ad8aa24742ea41ea36fa2e3452a31636e82"
	expectedPathName := "71056/ad8aa/24742/ea41e/a36fa/2e345/2a316/36e82"
	if pathKey.Pathname != expectedPathName {
		t.Errorf("have %s, want %s", pathKey.Pathname, expectedPathName)
	}
	if pathKey.Filename != expectedFilename {
		t.Errorf("have %s, want %s", pathKey.Filename, expectedFilename)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)
	key := "fooandbars"
	data := []byte("some foobar data")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s", key)
	}

	r, err := s.Read(key)

	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	fmt.Println(string(b))

	if string(b) != string(data) {
		t.Errorf("have %s, want %s", data, b)
	}
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
	if ok := s.Has(key); ok {
		t.Errorf("expected to not have key %s", key)
	}
}
