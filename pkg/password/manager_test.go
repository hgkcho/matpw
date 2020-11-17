package password

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestManager_ToCSV(t *testing.T) {
	type fields struct {
		Passwords        []Password
		Path             string
		CSVPath          string
		PasswordDataPath string
		PasscodeDataPath string
		Passcode         Passcode
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "[OK]: export successfully",
			fields: fields{
				CSVPath:          filepath.Join(os.Getenv("TMPDIR"), "output.csv"),
				Path:             filepath.Join(os.Getenv("HOME"), ".matpw", "password.json"),
				PasswordDataPath: filepath.Join(os.Getenv("TMPDIR"), "matpassword.json"),
				PasscodeDataPath: filepath.Join(os.Getenv("TMPDIR"), "passsecret.json"),
				Passcode: Passcode{
					Path: filepath.Join(os.Getenv("TMPDIR"), "secret.json"),
					rw:   &bytes.Buffer{},
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
			// setupTestData(t, *m)
			var p = Passcode{orderByte: []byte{0, 4, 6, 19}, rw: &bytes.Buffer{}}
			var w io.Writer
			if err := p.Encript(w); err != nil {
				log.Fatal(err)
			}
			m.Passcode.rw = p.rw
			if err := m.Load(); err != nil {
				t.Fatalf("[error]: %v", err)
			}
			var buf = new(bytes.Buffer)
			if err := m.ToCSV(buf); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			s := bufio.NewReader(buf)
			l, isP, err := s.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("isPrefix: %v, line: %v", isP, l)
			tests := []string{
				"ID", "Service", "Acount", "Description", "PasswordSet", "CreatedAt", "UpdatedAt",
			}
			actual := strings.Split(string(l), ",")
			if !reflect.DeepEqual(actual, tests) {
				t.Errorf("get unexpected header: get %v", tests)
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
			// m.Passcode.Encript()
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
	// pc.Encript()
	// m.Passcode = pc
	m.Save()
	m.SavePasscode()
}

func cleanUpTestData(t *testing.T, m Manager) {

}
