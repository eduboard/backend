package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	var envs = []struct {
		key   string
		value string
	}{
		{"HOST", "test"},
		{"MONGO_HOST", "testurl"},
		{"MONGO_PORT", "testport"},
		{"MONGO_DB", "testdb"},
		{"MONGO_USER", "testuser"},
		{"MONGO_PASS", "testpass"},
		{"STATIC_DIR", "testdir"},
		{"FILES_DIR", "testfilesdir"},
		{"LOGFILE", "backend.log"},
	}

	var testCases = []struct {
		name     string
		unset    bool
		expected config
	}{
		{"set", false,
			config{
				"test",
				"testurl",
				"testport",
				"testdb",
				"testuser",
				"testpass",
				"testdir",
				"testfilesdir",
				"backend.log"}},
		{"unset", true,
			config{
				":8080",
				"localhost",
				"27017",
				"eduboard",
				"",
				"",
				"./static",
				"./files",
				""}},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			for _, e := range envs {
				os.Unsetenv(e.key)
				if !v.unset {
					os.Setenv(e.key, e.value)
				}
			}

			c := GetConfig()
			cExpected := v.expected
			assert.Equal(t, cExpected, c, "config does not match expected one")
		})
	}
}
