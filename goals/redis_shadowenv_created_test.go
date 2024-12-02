package goals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisShadowenvCreatedFileContents(t *testing.T) {
	goal := RedisShadowenvCreated{
		Env: map[string]string{
			"REDIS_HOST_PORT": "{{ .Host }}:{{ .Port }}",
			"REDIS_HOST":      "{{ .Host }}",
		},
	}

	actualFileContents := string(goal.fileContents())

	assert.Contains(t, actualFileContents, `(provide "redis")`)
	assert.Contains(t, actualFileContents, `(env/set "REDIS_HOST_PORT" "localhost:6379")`)
	assert.Contains(t, actualFileContents, `(env/set "REDIS_HOST" "localhost")`)
}
