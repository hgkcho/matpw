package password

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestManager_ToCSV(t *testing.T) {
	type fields struct {
		Passwords []Password
		Path      string
		CSVPath   string
		Passcode  Passcode
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "[OK]: export successfully",
			fields: fields{
				CSVPath: "output.csv",
				Path:    filepath.Join(os.Getenv("HOME"), ".matpw", "password.json"),
				Passcode: Passcode{
					Path: filepath.Join(os.Getenv("HOME"), ".matpw", "secret.json"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				Passwords:        tt.fields.Passwords,
				PasswordDataPath: tt.fields.Path,
				CSVPath:          tt.fields.CSVPath,
				Passcode:         tt.fields.Passcode,
			}
			if err := m.Load(); err != nil {
				t.Fatalf("[error]: %v", err)
			}
			if err := m.ToCSV(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Init(t *testing.T) {
	type fields struct {
		Passwords        []Password
		PasswordDataPath string
		CSVPath          string
		Passcode         Passcode
		PasscodeDataPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "[OK]: password data is decoded",
			fields: fields{
				PasswordDataPath: filepath.Join(os.Getenv("TMPDIR"), "matpassword.json"),
				PasscodeDataPath: filepath.Join(os.Getenv("TMPDIR"), "passsecret.json"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Cleanup(func() {
			err := os.Remove(tt.fields.PasswordDataPath)
			if err != nil {
				log.Fatal(err)
			}
			err = os.Remove(tt.fields.PasscodeDataPath)
			if err != nil {
				log.Fatal(err)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				Passwords:        tt.fields.Passwords,
				PasswordDataPath: tt.fields.PasswordDataPath,
				CSVPath:          tt.fields.CSVPath,
				Passcode:         tt.fields.Passcode,
				PasscodeDataPath: tt.fields.PasscodeDataPath,
			}
			setupTestData(t, *m)
			if err := m.Init(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, v := range m.Passwords {
				var zero uuid.UUID
				if v.ID == zero {
					t.Errorf("[error]: failed to parse %v ", tt.fields.PasswordDataPath)
				}
			}
			var zero []byte
			if bytes.Equal(m.Passcode.OrderCode, zero) {
				t.Errorf("[error]: failed to parse %v ", tt.fields.PasscodeDataPath)
			}

		})
	}
}

func setupTestData(t *testing.T, m Manager) {
	t.Helper()
	var pw = Password{}
	var pc = Passcode{}
	pw.Create()
	m.Passwords = []Password{pw}
	pc.Order = []int{1, 4, 18, 23}
	pc.Encript()
	m.Passcode = pc
	m.Save()
	m.SavePasscode()
}
