package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "bestpicture"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "71056ad8aa24742ea41ea36fa2e3452a31636e82"
	expectedPathName := "71056/ad8aa/24742/ea41e/a36fa/2e345/2a316/36e82"
	if pathKey.Pathname != expectedPathName {
		t.Errorf("have %s, want %s", pathKey.Pathname, expectedPathName)
	}
	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("have %s, want %s", pathKey.Filename, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{PathTransformFunc: CASPathTransformFunc}
	s := NewStore(opts)
	key := "specials"

	data := []byte("some data")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
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
}
