package rp

import (
	"regexp"
	"testing"
)

func BenchmarkGo(b *testing.B) {
	emails := []string{
		"user@example.com",
		"john.doe@gmail.com",
		"alice@outlook.co",
		"bob@company.org",
		"test+label@domain.net",
		"pretty-long-email-address@subdomain.example.co.uk",
		"a@b.c",
		"invalid-email",
		"@domain.com",
		"user@",
	}

	e, err := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		for _, email := range emails {
			e.MatchString(email)
		}
	}
}
