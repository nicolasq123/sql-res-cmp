package cmp

import (
	"testing"
)

func TestNewDB_InvalidDSN(t *testing.T) {
	_, err := NewDB("invalid://dsn")
	if err == nil {
		t.Error("expected error for invalid dsn")
	}
}

func TestParseMySQLURLToDSN(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "simple format",
			url:  "mysql://root:pass@127.0.0.1:3306/mydb?parseTime=true",
			want: "mysql://root:pass@127.0.0.1:3306/mydb?parseTime=true",
		},
		{
			name: "tcp format",
			url:  "mysql://root:pass@tcp(127.0.0.1:3306)/mydb?parseTime=true",
			want: "mysql://root:pass@127.0.0.1:3306/mydb?parseTime=true",
		},
		{
			name: "with query params",
			url:  "mysql://root:pass@127.0.0.1:3306/mydb",
			want: "mysql://root:pass@127.0.0.1:3306/mydb",
		},
		{
			name: "preserve existing parseTime",
			url:  "mysql://root:pass@127.0.0.1:3306/mydb?parseTime=false",
			want: "mysql://root:pass@127.0.0.1:3306/mydb?parseTime=false",
		},
		{
			name: "with special chars in password",
			url:  "mysql://root:pass@tcp(127.0.0.1:3306)/mysql?parseTime=true&sql_mode=TRADITIONAL",
			want: "mysql://root:pass@127.0.0.1:3306/mysql?parseTime=true&sql_mode=TRADITIONAL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMySQLURL(tt.url)
			if got != tt.want {
				t.Errorf("GOT %v, want %v", got, tt.want)
			}
		})
	}
}
