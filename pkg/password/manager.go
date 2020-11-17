package password

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	WorkDir             = filepath.Join(os.Getenv("HOME"), ".matpw")
	DefaultPasswordPath = filepath.Join(WorkDir, "password.json")
	DefaultSecretPath   = filepath.Join(WorkDir, "secret.json")
)

// Manager represents password manager
type Manager struct {
	Passwords []Password
	Passcode
	PasswordDataPath string
	CSVPath          string
	PasscodeDataPath string
	passcodeBuf      io.ReadWriter
	output           io.Writer
}

// Init decode json data and set them
func (m *Manager) Init() error {
	// TODO consider design
	m.passcodeBuf = new(bytes.Buffer)
	m.output = new(bytes.Buffer)

	err := os.MkdirAll(WorkDir, 0755)
	if err != nil {
		return err
	}
	if m.PasswordDataPath == "" {
		m.PasswordDataPath = DefaultPasswordPath
	}
	if m.PasscodeDataPath == "" {
		m.PasscodeDataPath = DefaultPasscodePath
	}
	f, err := os.Open(m.PasswordDataPath)
	if err != nil {
		return err
	}
	err = json.NewDecoder(f).Decode(&m.Passwords)
	if err != nil {
		return err
	}
	f.Close()

	f, err = os.Open(m.PasscodeDataPath)
	if err != nil {
		return err
	}
	err = json.NewDecoder(f).Decode(&m.Passcode)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}

// ToCSV export password data
// TODO  引数をio.Writerにする
func (m *Manager) ToCSV(w io.Writer) error {
	// TODO これはinitで共通処理にする?
	f, err := os.Open(DefaultSecretPath)
	if err != nil {
		return err
	}
	err = m.Passcode.Decript(f)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(w)
	writer.Comma = ','
	writer.Write([]string{
		"ID", "Service", "Acount", "Description", "Password", "CreatedAt", "UpdatedAt",
	})
	// var buf = [][]string{
	// 	{"ID", "Service", "Acount", "Description", "PasswordSet", "CreatedAt", "UpdatedAt"},
	// }
	for _, v := range m.Passwords {
		var pw string
		pw += v.PasswordSet[m.Passcode.Order[0]] + v.PasswordSet[m.Passcode.Order[1]] + v.PasswordSet[m.Passcode.Order[2]] + v.PasswordSet[m.Passcode.Order[3]]
		// buf = append(buf, []string{v.ID.String(), v.Account, v.Descripiton, pw, v.CreatedAt.String(), v.UpdatedAt.String()})
		err = writer.Write([]string{v.ID.String(), v.Account, v.Descripiton, pw, v.CreatedAt.String(), v.UpdatedAt.String()})
		if err != nil {
			return err
		}
	}
	// err = writer.WriteAll(buf)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

// Delete delete sepecified password by service name
func (m *Manager) Delete(p Password) error {
	fmt.Println("Are you sure you want to delete it?")
	var buf []Password
	for _, v := range m.Passwords {
		if p.ID != v.ID {
			buf = append(buf, v)
		}
	}
	return nil
}

func saveJSON(path string, v interface{}) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	// save as pretty json
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	f.Write(buf)
	return nil
}

func loadJSON(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

// Save save passwords to m.Path
func (m *Manager) Save() error {
	return saveJSON(m.PasswordDataPath, m.Passwords)
}

// Load load password from m.Path
func (m *Manager) Load() error {
	return loadJSON(m.PasswordDataPath, &m.Passwords)
}

// SavePasscode save Passcode
func (m *Manager) SavePasscode() error {
	return saveJSON(m.PasscodeDataPath, m.Passcode)
}

// LoadPasscode save Passcode
func (m *Manager) LoadPasscode() error {
	return loadJSON(m.PasscodeDataPath, &m.Passcode)
}
