package api

import (
	"testing"
)

// TestAllowedHosts_ExactEntries verifies that AllowedHosts contains exactly
// the three canonical RIS hosts and no others.
func TestAllowedHosts_ExactEntries(t *testing.T) {
	expected := map[string]bool{
		"data.bka.gv.at":    true,
		"www.ris.bka.gv.at": true,
		"ris.bka.gv.at":     true,
	}

	if len(AllowedHosts) != len(expected) {
		t.Fatalf("AllowedHosts hat %d Einträge, erwartet %d", len(AllowedHosts), len(expected))
	}

	for host := range expected {
		if !AllowedHosts[host] {
			t.Errorf("AllowedHosts fehlt erwarteter Host %q", host)
		}
	}

	for host := range AllowedHosts {
		if !expected[host] {
			t.Errorf("AllowedHosts enthält unerwarteten Host %q", host)
		}
	}
}

// TestValidateDocURL_AllowedHosts verifies that all three canonical hosts are accepted.
func TestValidateDocURL_AllowedHosts(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"data.bka.gv.at", "https://data.bka.gv.at/ris/api/v2.6/Bundesrecht"},
		{"www.ris.bka.gv.at", "https://www.ris.bka.gv.at/Dokument.wxe?Abfrage=Bundesnormen&Dokumentnummer=NOR40045109"},
		{"ris.bka.gv.at", "https://ris.bka.gv.at/GeltendeFassung.wxe?Abfrage=Bundesnormen&Gesetzesnummer=10001622"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDocURL(tt.url); err != nil {
				t.Errorf("validateDocURL(%q) = %v, want nil", tt.url, err)
			}
		})
	}
}

// TestValidateDocURL_RejectedHosts verifies that unknown hosts are rejected.
func TestValidateDocURL_RejectedHosts(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"evil.com", "https://evil.com/steal-data"},
		{"localhost", "https://localhost/internal"},
		{"127.0.0.1", "https://127.0.0.1/admin"},
		{"internal metadata", "https://169.254.169.254/latest/meta-data/"},
		{"subdomain spoof", "https://fake.data.bka.gv.at/something"},
		{"similar domain", "https://data.bka.gv.at.evil.com/phish"},
		{"empty host", "https:///no-host"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDocURL(tt.url)
			if err == nil {
				t.Errorf("validateDocURL(%q) = nil, want error for disallowed host", tt.url)
			}
		})
	}
}

// TestValidateDocURL_RejectsNonHTTPS verifies that non-HTTPS schemes are rejected.
func TestValidateDocURL_RejectsNonHTTPS(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"http", "http://data.bka.gv.at/ris/api/v2.6/Bundesrecht"},
		{"ftp", "ftp://data.bka.gv.at/file"},
		{"file", "file:///etc/passwd"},
		{"javascript", "javascript:alert(1)"},
		{"empty scheme", "://data.bka.gv.at/something"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDocURL(tt.url)
			if err == nil {
				t.Errorf("validateDocURL(%q) = nil, want error for non-HTTPS scheme", tt.url)
			}
		})
	}
}

// TestValidateDocURL_MalformedURLs verifies that malformed URLs are rejected.
func TestValidateDocURL_MalformedURLs(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"empty string", ""},
		{"whitespace only", "   "},
		{"no scheme", "data.bka.gv.at/something"},
		{"colon only", ":"},
		{"double scheme", "https://https://data.bka.gv.at"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDocURL(tt.url)
			if err == nil {
				t.Errorf("validateDocURL(%q) = nil, want error for malformed URL", tt.url)
			}
		})
	}
}

// TestValidateDocURL_CaseInsensitiveHost verifies that host matching is case-insensitive.
func TestValidateDocURL_CaseInsensitiveHost(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"uppercase", "https://DATA.BKA.GV.AT/ris/api/v2.6/Bundesrecht"},
		{"mixed case", "https://Data.Bka.Gv.At/ris/api/v2.6/Bundesrecht"},
		{"mixed case ris", "https://WWW.RIS.BKA.GV.AT/Dokument.wxe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDocURL(tt.url); err != nil {
				t.Errorf("validateDocURL(%q) = %v, want nil (case-insensitive match)", tt.url, err)
			}
		})
	}
}
