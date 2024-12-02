package goals

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresqlShadowenvCreatedFileContents(t *testing.T) {
	goal := PostgresqlShadowenvCreated{
		Version: 17,
		Env: map[string]string{
			"DB_HOST_PORT": "{{ .Host }}:{{ .Port }}",
			"DB_USER":      "{{ .User }}",
		},
	}

	actualFileContents := string(goal.fileContents())

	assert.Contains(t, actualFileContents, `(provide "postgresql" "17")`)
	assert.Contains(t, actualFileContents, `(env/set "DB_HOST_PORT" "localhost:5432")`)
	assert.Contains(t, actualFileContents, fmt.Sprintf(`(env/set "DB_USER" %q)`, os.Getenv("USER")))
	assert.Contains(t, actualFileContents, `(env/prepend-to-pathlist "PATH" "/opt/homebrew/opt/postgresql@17/bin")`)
}
