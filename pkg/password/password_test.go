package password

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_generateTwoChar(t *testing.T) {
	tests := []struct {
		name    string
		useChar string
	}{
		{
			name:    "OK two strings",
			useChar: "ab",
		},
		{
			name:    "OK with upper string",
			useChar: "abAB",
		},
		{
			name:    "OK with digit",
			useChar: "ab12345",
		},
		{
			name:    "OK with symbol",
			useChar: "ab123&^]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateTwoChar(tt.useChar)
			if len(got) != 2 {
				t.Errorf("generateTwoChar() length = %v, want %v", len(got), 2)
			}
			for _, v := range got {
				if !strings.Contains(tt.useChar, string(v)) {
					t.Errorf("generateTwoChar() generate invalid string. allow strings: %s but, got %s", tt.useChar, string(v))
				}
			}
		})
	}
}

func TestPassword_Create(t *testing.T) {
	type fields struct {
		Service      string
		Account      string
		Descripiton  string
		PasswordSet  []string
		CreatedAt    time.Time
		UpdatedAt    time.Time
		Path         string
		UseUppercase bool
		UseDigit     bool
		UseSymbol    bool
	}
	tests := []struct {
		name        string
		fields      fields
		allowString string
		wantErr     bool
	}{
		{
			name: "[OK]: passwordSet allows lower case only",
			fields: fields{
				UseUppercase: false,
				UseDigit:     false,
				UseSymbol:    false,
			},
			allowString: Lower,
			wantErr:     false,
		},
		{
			name: "[OK]: passwordSet allows lower case and upper case",
			fields: fields{
				UseUppercase: true,
				UseDigit:     false,
				UseSymbol:    false,
			},
			allowString: Lower + Upper,
			wantErr:     false,
		},
		{
			name: "[OK]: passwordSet allows lower case and digits",
			fields: fields{
				UseUppercase: false,
				UseDigit:     true,
				UseSymbol:    false,
			},
			allowString: Lower + Digits,
			wantErr:     false,
		},
		{
			name: "[OK]: passwordSet allows lower case and symbols",
			fields: fields{
				UseUppercase: false,
				UseDigit:     false,
				UseSymbol:    true,
			},
			allowString: Lower + Symbols,
			wantErr:     false,
		},
		{
			name: "[OK]: passwordSet allows lower case and upper case and digits",
			fields: fields{
				UseUppercase: true,
				UseDigit:     true,
				UseSymbol:    false,
			},
			allowString: Lower + Upper + Digits,
			wantErr:     false,
		},
		{
			name: "[OK]: passwordSet allows lower case and upper case and symbols",
			fields: fields{
				UseUppercase: true,
				UseDigit:     false,
				UseSymbol:    true,
			},
			allowString: Lower + Upper + Symbols,
			wantErr:     false,
		},
		{
			name: "[OK]: passwordSet allows lower case and upper case and digits and symbols",
			fields: fields{
				UseUppercase: true,
				UseDigit:     true,
				UseSymbol:    true,
			},
			allowString: Lower + Upper + Digits + Symbols,
			wantErr:     false,
		},
		{
			name:        "[OK]: ID is filled",
			allowString: Lower,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Password{
				Service:      tt.fields.Service,
				Account:      tt.fields.Account,
				Descripiton:  tt.fields.Descripiton,
				PasswordSet:  tt.fields.PasswordSet,
				CreatedAt:    tt.fields.CreatedAt,
				UpdatedAt:    tt.fields.UpdatedAt,
				Path:         tt.fields.Path,
				UseUppercase: tt.fields.UseUppercase,
				UseDigit:     tt.fields.UseDigit,
				UseSymbol:    tt.fields.UseSymbol,
			}
			if err := p.Create(); (err != nil) != tt.wantErr {
				t.Errorf("Password.Create got err = %v, wantErr = %v", err, tt.wantErr)
			}
			if time.Now().Sub(p.CreatedAt) > 1*time.Second {
				t.Errorf("Password.Create() gotten createAt is %v , this is invalid", p.CreatedAt)
			}
			for _, v := range p.PasswordSet {
				if len(v) != 2 {
					t.Errorf("Password set length should be 2 but got %x", len(v))
				}
				for _, vv := range v {
					if !strings.Contains(tt.allowString, string(vv)) {
						t.Errorf("Password set contains %v but is not allowed letter", v)
					}
				}
			}
			// check ID is set
			var b uuid.UUID
			if p.ID == b {
				t.Errorf("[error]: ID is should be filled. got %v", p.ID)
			}

		})
	}
}
