package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hgkcho/matpw/pkg/password"
)

type meta struct {
	passwords []password.Password
	pwFilePath string
}

func (m *meta) init() error {
	workDir := filepath.Join(os.Getenv("HOME"), ".matpw")
	os.MkdirAll(workDir, 0755)
	m.pwFilePath = filepath.Join(workDir, "password.json")
	f, err := os.OpenFile(filepath.Join(workDir, "password.json"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	m.passwords = []password.Password{}
	// ioutil.ReadFile(filepath.Join(workDir, "password.json"))
	err = json.NewDecoder(f).Decode(&m.passwords)
	if err != nil {
		return err
	}
	// json.NewEncoder(f).Encode(m.password)
	return nil
}
