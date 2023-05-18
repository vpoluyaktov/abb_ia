package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
)

var lock sync.Mutex

var Marshal = func(obj interface{}) (io.Reader, error) {
  b, err := json.MarshalIndent(obj, "", "  ")
  if err != nil {
    return nil, err
  }
  return bytes.NewReader(b), nil
}

var Unmarshal = func(r io.Reader, obj interface{}) error {
  return json.NewDecoder(r).Decode(obj)
}

// Save saves a representation of object to the file at path.
func DumpJson(file string, obj interface{}) error {
  lock.Lock()
  defer lock.Unlock()
  f, err := os.Create(file)
  if err != nil {
    return err
  }
  defer f.Close()
  r, err := Marshal(obj)
  if err != nil {
    return err
  }
  _, err = io.Copy(f, r)
  return err
}

// loads the file at path into obj.
func LoadJson(file string, obj interface{}) error {
  lock.Lock()
  defer lock.Unlock()
  f, err := os.Open(file)
  if err != nil {
    return err
  }
  defer f.Close()
  return Unmarshal(f, obj)
}