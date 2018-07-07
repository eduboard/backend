package url

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestStringifyURLs(t *testing.T) {
	u, _ := url.Parse("https://example.com")
	var testCases = []struct {
		name     string
		input    []url.URL
		expected []string
	}{
		{"empty", []url.URL{}, []string{}},
		{"success", []url.URL{*u}, []string{"https://example.com"}},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			s := StringifyURLs(v.input...)
			assert.ElementsMatch(t, v.expected, s, "stringified urls do not match")
		})
	}
}

func TestURLifyStrings(t *testing.T) {
	u, _ := url.Parse("https://example.com")
	var testCases = []struct {
		name     string
		input    []string
		expected []url.URL
		error    bool
	}{
		{"empty", []string{}, []url.URL{}, false},
		{"success", []string{"https://example.com"}, []url.URL{*u}, false},
		{"bad url", []string{"htttp\\:.orgcom"}, []url.URL{}, true},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			o, e := URLifyStrings(v.input...)
			assert.ElementsMatch(t, v.expected, o, "elements do not match")
			if v.error {
				assert.Error(t, e, "no error returned")
			}
		})
	}
}
