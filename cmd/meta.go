package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hgkcho/matpw/pkg/password"
)

var (
	WorkDir             = filepath.Join(os.Getenv("HOME"), ".matpw")
	DefaultPasswordPath = filepath.Join(WorkDir, "password.json")
)

type meta struct {
	passwords  []password.Password
	passcode   password.Passcode
	pwFilePath string
}

func (m *meta) init() error {
	os.MkdirAll(WorkDir, 0755)
	m.pwFilePath = DefaultPasswordPath
	f, err := os.OpenFile(m.pwFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	m.passwords = []password.Password{}
	// ioutil.ReadFile(filepath.Join(workDir, "password.json"))
	return json.NewDecoder(f).Decode(&m.passwords)
}

// Save save passwords to password.json
func (m *meta) Save() error {
	f, err := os.Create(m.pwFilePath)
	// f, err := os.OpenFile(c.meta.pwFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
	defer f.Close()
	if err != nil {
		return err
	}
	// save as pretty json
	buf, err := json.MarshalIndent(m.passwords, "", "  ")
	// err = json.NewEncoder(f).Encode(c.meta.passwords)
	if err != nil {
		return err
	}
	f.Write(buf)
	return nil
}
