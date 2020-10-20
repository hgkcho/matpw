package password

import (
	"crypto/cipher"
	"os"
	"path/filepath"
	"testing"
)

func TestPasscode_Encript(t *testing.T) {
	type fields struct {
		Order     []int
		OrderByte []byte
		OrderCode []byte
		Path      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Order: []int{0, 4, 20, 24},
				Path:  filepath.Join(os.Getenv("HOME"), ".matpw", "secret.json"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Passcode{
				Order:     tt.fields.Order,
				OrderCode: tt.fields.OrderCode,
				Path:      tt.fields.Path,
			}
			if err := p.Encript(); (err != nil) != tt.wantErr {
				t.Errorf("Passcode.Encript() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasscode_Decript(t *testing.T) {
	type fields struct {
		Order     []int
		OrderCode []byte
		Path      string
		cipher    cipher.Block
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Path: filepath.Join(os.Getenv("HOME"), ".matpw", "secret.json"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Passcode{
				Order:     tt.fields.Order,
				OrderCode: tt.fields.OrderCode,
				Path:      tt.fields.Path,
				cipher:    tt.fields.cipher,
			}
			if err := p.Decript(); (err != nil) != tt.wantErr {
				t.Errorf("Passcode.Decript() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
