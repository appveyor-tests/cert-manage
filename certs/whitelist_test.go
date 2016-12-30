package certs

import (
	"testing"
)

func TestWhitelist_HexSignature(t *testing.T) {
	certificates := loadPEM(t, "../testdata/example.crt")
	if len(certificates) != 1 {
		t.Fatal("Found != 1 certs in example.crt")
	}
	if certificates[0] == nil {
		t.Fatal("Unable to read cert")
	}

	// Empty signature
	wh1 := HexSignatureWhitelistItem{
		Signature: "",
	}
	if wh1.Matches(*certificates[0]) {
		t.Fatal("wh1 and cert shouldn't match")
	}

	// Signature too short
	wh2 := HexSignatureWhitelistItem{
		Signature: "",
	}
	if wh2.Matches(*certificates[0]) {
		t.Fatal("wh2 and cert shouldn't match")
	}

	// Short signature
	wh3 := HexSignatureWhitelistItem{
		Signature: "96940d99",
	}
	if !wh3.Matches(*certificates[0]) {
		t.Fatal("wh3 and cert don't match, but should")
	}

	// Full signature
	wh4 := HexSignatureWhitelistItem{
		Signature: "96940d991419151450d1e75f66218f6f2594e1df4af31a5ad673c9a8746817ce",
	}
	if !wh4.Matches(*certificates[0]) {
		t.Fatal("wh4 and cert don't match, but should")
	}
}

func TestWhitelist_IssuerCommonName(t *testing.T) {
	certificates := loadPEM(t, "../testdata/example.crt")
	if len(certificates) != 1 {
		t.Fatal("Found != 1 certs in example.crt")
	}
	if certificates[0] == nil {
		t.Fatal("Unable to read cert")
	}

	// Empty whitelist name
	wh1 := IssuersCommonNameWhitelistItem{
		Name: "",
	}
	if wh1.Matches(*certificates[0]) {
		t.Fatal("wh1 and cert should not match")
	}

	// Matching word
	wh2 := IssuersCommonNameWhitelistItem{
		Name: "Starfield",
	}
	if !wh2.Matches(*certificates[0]) {
		t.Fatal("wh2 and cert don't match")
	}

	// Full Match
	wh3 := IssuersCommonNameWhitelistItem{
		Name: "Starfield Secure Certification Authority",
	}
	if !wh3.Matches(*certificates[0]) {
		t.Fatal("wh3 and cert don't match")
	}
}