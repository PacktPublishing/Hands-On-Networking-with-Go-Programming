package wordlist

import "testing"

func TestPGPWordSizes(t *testing.T) {
	if len(odd) != 256 {
		t.Errorf("expected odd word length 256, got %d", len(odd))
	}
	if len(even) != 256 {
		t.Errorf("expected even word length 256, got %d", len(even))
	}
}

func TestPGPArray(t *testing.T) {
	data := PGP([]byte{0, 0})
	if data[0] != "aardvark" {
		t.Errorf("expected 'aardvark', got %s", data[0])
	}
	if data[1] != "adroitness" {
		t.Errorf("expected 'adroitness', got %s", data[1])
	}
}
