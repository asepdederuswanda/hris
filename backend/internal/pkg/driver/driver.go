// Package driver mendefinisikan tipe database driver dan utilitas
// yang digunakan bersama oleh config dan database package.
package driver

import "strings"

// Type mendefinisikan tipe database driver yang didukung.
type Type string

const (
	Postgres Type = "postgres"
	MySQL    Type = "mysql"
)

// Parse mengkonversi string driver ke Type.
// Return Postgres jika tidak dikenal (default).
func Parse(driver string) Type {
	switch strings.ToLower(driver) {
	case "mysql":
		return MySQL
	default:
		return Postgres
	}
}

// String mengembalikan representasi string dari driver.
func (t Type) String() string {
	return string(t)
}

// IsValid memeriksa apakah driver didukung.
func (t Type) IsValid() bool {
	return t == Postgres || t == MySQL
}
