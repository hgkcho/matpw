package password

import "testing"

func TestPassword_Create(t *testing.T) {
	type fields struct {
		Title       string
		Account     string
		Descripiton string
		PasswordSet []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test ok",
			fields: fields{
				Title: "title",
				Account: "aaa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Password{
				Title:       tt.fields.Title,
				Account:     tt.fields.Account,
				Descripiton: tt.fields.Descripiton,
				PasswordSet: tt.fields.PasswordSet,
			}
			p.Create()
		})
	}
}
