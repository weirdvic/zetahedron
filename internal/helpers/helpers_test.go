package helpers

import (
	"os"
	"testing"
)

func TestBase62Encode(t *testing.T) {
	want := "Ok"
	got := Base62Encode(660)

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestBase62Decode(t *testing.T) {
	var want uint64 = 660
	got, _ := Base62Decode("Ok")

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestEnforceHTTPS(t *testing.T) {
	want := "https://nethack.org"
	got := EnforceHTTPS("http://nethack.org")
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
	got = EnforceHTTPS("nethack.org")
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestIsOurDomain(t *testing.T) {
	os.Setenv("DOMAIN", "http://localhost:1323")
	want := true
	got := IsOurDomain("http://localhost:1323")
	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}
	want = false
	got = IsOurDomain("https://nethack.org/")
	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}
}
