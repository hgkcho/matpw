package password

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var (
	WorkDir             = filepath.Join(os.Getenv("HOME"), ".matpw")
	DefaultPasswordPath = filepath.Join(WorkDir, "password.json")
)

// Manager represents password manager
type Manager struct {
	Passwords        []Password
	PasswordDataPath string
	CSVPath          string
	Passcode
	PasscodeDataPath string
}

// Init decode json data and set them
func (m *Manager) Init() error {
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
func (m *Manager) ToCSV() error {
	f, err := os.Create(m.CSVPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = m.Passcode.Decript()

	if err != nil {
		return err
	}

	writer := csv.NewWriter(f)
	writer.Comma = ','
	var buf = [][]string{
		{"ID", "Service", "Acount", "Description", "PasswordSet", "CreatedAt", "UpdatedAt"},
	}
	for _, v := range m.Passwords {
		var pw string
		pw += v.PasswordSet[m.Passcode.OrderByte[0]] + v.PasswordSet[m.Passcode.OrderByte[1]] + v.PasswordSet[m.Passcode.OrderByte[2]] + v.PasswordSet[m.Passcode.OrderByte[3]]
		buf = append(buf, []string{v.ID.String(), v.Account, v.Descripiton, pw, v.CreatedAt.String(), v.UpdatedAt.String()})
	}
	err = writer.WriteAll(buf)
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

func save(path string, v interface{}) error {
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

func load(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

// Save save passwords to m.Path
func (m *Manager) Save() error {
	// f, err := os.Create(m.PasswordDataPath)
	// defer f.Close()
	// if err != nil {
	// 	return err
	// }
	// // save as pretty json
	// buf, err := json.MarshalIndent(m.Passwords, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// f.Write(buf)
	// return nil
	return save(m.PasswordDataPath, m.Passwords)
}

// Load load password from m.Path
func (m *Manager) Load() error {
	// f, err := os.Open(m.PasswordDataPath)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()
	// return json.NewDecoder(f).Decode(&m.Passwords)
	return load(m.PasswordDataPath, &m.Passwords)
}

func (m *Manager) SavePasscode() error {
	return save(m.PasscodeDataPath, m.Passcode)
}

func (m *Manager) LoadPasscode() error {
	return load(m.PasscodeDataPath, &m.Passcode)
}
