package whitelist

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	// "time"

	"github.com/adamdecaf/cert-manage/tools/file"
)

// TOOD(adam): Read and review this code
// https://blog.hboeck.de/archives/888-How-I-tricked-Symantec-with-a-Fake-Private-Key.html

// todo: dedup certs already added by one whitelist item
// e.g. If my []Item contains a signature and Issuer.CommonName match
// don't add the cert twice

// Item can be compared against an x509 Certificate to see if the cert represents
// some value presented by the whitelist item. This is useful in comparing specific fields of
// Certificate against multiple whitelist candidates.
type Item interface {
	Matches(x509.Certificate) bool
}

// findRemovable a list of x509 Certificates against whitelist items to
// retain only the certificates that are disallowed by our whitelist.
// An empty slice of certificates is a possible (and valid) output.
func findRemovable(incoming []*x509.Certificate, whitelisted []Item) []*x509.Certificate {
	// Pretty bad search right now.
	var removable []*x509.Certificate

	for _, inc := range incoming {
		remove := true
		// If the whitelist matches on something then don't remove it
		for _, wh := range whitelisted {
			if inc != nil && wh.Matches(*inc) {
				remove = false
			}
		}
		if remove {
			removable = append(removable, inc)
		}
	}

	return removable
}

// Json structure in struct form
type jsonWhitelist struct {
	Signatures jsonSignatures `json:"Signatures,omitempty"`
}
type jsonSignatures struct {
	Hex []string `json:"Hex"`
}

// loadFromFile reads a whitelist file and parses it into Items
func loadFromFile(path string) ([]Item, error) {
	if !validWhitelistPath(path) {
		return nil, fmt.Errorf("The path '%s' doesn't seem to contain a whitelist.", path)
	}

	// read file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var parsed jsonWhitelist
	err = json.Unmarshal(b, &parsed)
	if err != nil {
		return nil, err
	}

	// Read parsed format into structs
	var items []Item
	for _, s := range parsed.Signatures.Hex {
		items = append(items, fingerprint{Signature: s})
	}

	return items, nil
}

// validWhitelistPath verifies that the given whitelist filepath is properly defined
// and exists on the given filesystem.
func validWhitelistPath(path string) bool {
	if !file.Exists(path) {
		fmt.Printf("The path %s doesn't seem to exist.\n", path)
	}

	isFlag := strings.HasPrefix(path, "-")
	if len(path) == 0 || isFlag {
		fmt.Printf("The given whitelist file path '%s' doesn't look correct.\n", path)
		if isFlag {
			fmt.Println("The path looks like a cli flag, -whitelist requires a path to the whitelist file.")
		} else {
			fmt.Println("The given whitelist file path is empty.")
		}
		return false
	}

	return true
}
