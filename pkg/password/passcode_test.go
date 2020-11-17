package password

import (
	"bytes"
	"reflect"
	"testing"
)

func TestPasscode_Encript(t *testing.T) {
	type fields struct {
		Order []int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "[OK]",
			fields: fields{
				Order: []int{0, 4, 20, 24},
			},
			wantErr: false,
		},
		{
			name:   "[Fail]: when p.Order is nil",
			fields: fields{
				// Order: []int{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Passcode{
				Order: tt.fields.Order,
			}
			var w = new(bytes.Buffer)
			if err := p.Encript(w); (err != nil) != tt.wantErr {
				t.Errorf("Passcode.Encript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if p.OrderCode == nil && !tt.wantErr {
				t.Error("Passcode.OrderCode is zero value")
				return
			}
		})
	}
}

func TestPasscode_Decript(t *testing.T) {
	tests := []struct {
		name     string
		order    []int
		wantErr  bool
		resetBuf bool
	}{
		{
			name:    "[OK]",
			order:   []int{0, 5, 13, 18},
			wantErr: false,
		},
		{
			name:    "[OK]",
			order:   []int{5, 1, 3, 21},
			wantErr: false,
		},
		{
			name:     "[Fail]: cannot read data from io.Reader",
			order:    []int{5, 1, 3, 18},
			wantErr:  true,
			resetBuf: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepase faze: create encoded data and write it on buffer
			tmp := &Passcode{
				Order: tt.order,
			}
			var buf = new(bytes.Buffer)
			if err := tmp.Encript(buf); err != nil {
				t.Fatal(err)
			}

			if tt.resetBuf {
				buf.Reset()
			}

			// Core test
			actual := &Passcode{}
			if err := actual.Decript(buf); (err != nil) != tt.wantErr {
				t.Errorf("Passcode.Decript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.order, actual.Order) && !tt.wantErr {
				t.Errorf("[error]: expected is %v, but actual is %v", tt.order, actual.Order)
				return
			}
		})
	}
}
