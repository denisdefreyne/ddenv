package goals

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"denisdefreyne.com/x/ddenv/core"
)

const PostgresqlShadowenvCreated_Path = ".shadowenv.d/300_postgresql.lisp"

type PostgresqlShadowenvCreated struct {
	Version int
}

func (g PostgresqlShadowenvCreated) Description() string {
	return fmt.Sprintf("Adding PostgreSQL %v to Shadowenv", g.Version)
}

func (g PostgresqlShadowenvCreated) IsAchieved() bool {
	_, err := os.Lstat(PostgresqlShadowenvCreated_Path)
	if err != nil {
		return false
	}

	oldContents, err := os.ReadFile(PostgresqlShadowenvCreated_Path)
	if err != nil {
		return false
	}

	if !bytes.Equal(oldContents, g.fileContents()) {
		return false
	}

	return true
}

func (g PostgresqlShadowenvCreated) Achieve() error {
	err := os.WriteFile(PostgresqlShadowenvCreated_Path, g.fileContents(), 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("shadowenv", "trust")

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, string(stdoutStderr))
	}

	return nil
}

func (g PostgresqlShadowenvCreated) PreGoals() []core.Goal {
	return []core.Goal{
		ShadowenvDirCreated{},
	}
}

func (g PostgresqlShadowenvCreated) fileContents() []byte {
	// TODO: Allow configuring the database prefix (not just POSTGRES_)

	data := struct {
		Version      int
		Port         int16
		User         string
		EnvVarPrefix string
	}{Version: g.Version, Port: 5432, User: os.Getenv("USER"), EnvVarPrefix: "POSTGRES"}

	// FIXME: Use `brew --prefix …` instead of hardcoding the path
	templateContent := `(provide "postgresql" "{{ .Version }}")

(env/set "{{ .EnvVarPrefix }}_USER" "{{ .User }}")
(env/set "{{ .EnvVarPrefix }}_PASSWORD" "")
(env/set "{{ .EnvVarPrefix }}_HOST" "localhost")
(env/set "{{ .EnvVarPrefix }}_PORT" "{{ .Port }}")

(env/prepend-to-pathlist "PATH" "/opt/homebrew/opt/postgresql@{{ .Version }}/bin")
`

	tmpl, err := template.New("shadowenv").Parse(templateContent)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		panic(err)
	}

	return b.Bytes()
}
