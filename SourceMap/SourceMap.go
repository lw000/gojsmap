package SourceMap

import (
	"io/ioutil"

	"github.com/go-sourcemap/sourcemap"
)

type Manager struct {
	ok     bool
	source *sourcemap.Consumer
}

func NewSourceMapManager() *Manager {
	return &Manager{}
}

func (s *Manager) Parse(filename string) error {
	var (
		err  error
		data []byte
	)

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		s.ok = false
		return err
	}

	s.source, err = sourcemap.Parse(filename, data)
	if err != nil {
		s.ok = false
		return err
	}

	// var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// err = json.Unmarshal(data, &s.smap)
	// if err != nil {
	//	return err
	// }

	s.ok = true

	return nil
}

func (s *Manager) Get(r, c int) (source, name string, line, column int, ok bool) {
	return s.source.Source(r, c)
}
