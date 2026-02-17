package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "bestpicture"
	pathName := CASPathTransformFunc(key)
	expectedPathName := "71056/ad8aa/24742/ea41e/a36fa/2e345/2a316/36e82"
	fmt.Println(pathName)
	if pathName != expectedPathName {
		t.Errorf("have %s, want %s", pathName, expectedPathName)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{PathTransformFunc: CASPathTransformFunc}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some data"))
	if err := s.writeStream("somekey", data); err != nil {
		t.Error(err)
	}
}
