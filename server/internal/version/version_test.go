package version

import "testing"

func TestShortRevision(t *testing.T) {
	tests := map[string]string{
		"":           "",
		"abc":        "abc",
		"12345678":   "12345678",
		"1234567890": "12345678",
	}

	for revision, want := range tests {
		if got := shortRevision(revision); got != want {
			t.Errorf("shortRevision(%q) = %q, want %q", revision, got, want)
		}
	}
}
