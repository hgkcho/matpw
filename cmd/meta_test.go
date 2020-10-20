package cmd

import (
	"testing"

	"github.com/hgkcho/matpw/pkg/password"
)

func Test_meta_init(t *testing.T) {
	type fields struct {
		password []password.Password
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test",
			fields: fields{
				password: []password.Password{
					{
						Service: "aaaaa",
						Account: "aaaaa",
					},
					{
						Service: "bbbb",
						Account: "bbbb",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &meta{
				passwords: tt.fields.password,
			}
			m.init()
		})
	}
}
