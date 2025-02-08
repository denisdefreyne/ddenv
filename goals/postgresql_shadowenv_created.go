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
	Env     map[string]string
}

func (g PostgresqlShadowenvCreated) Description() string {
	return fmt.Sprintf("Adding PostgreSQL %v to Shadowenv", g.Version)
}

func (g PostgresqlShadowenvCreated) HashIdentity() string {
	return fmt.Sprintf("PostgresqlShadowenvCreated %v", g)
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
		ShadowenvSetUp{},
	}
}

func (g PostgresqlShadowenvCreated) fileContents() []byte {
	dataForSingleLines := struct {
		Version  int
		Host     string
		Port     int16
		User     string
		Password string
	}{
		Version:  g.Version,
		Host:     "localhost",
		Port:     5432,
		User:     os.Getenv("USER"),
		Password: "",
	}

	// Evaluate each template
	compiledEnvVarMap := make(map[string]string)
	for envName, envValueTemplate := range g.Env {
		// Template line
		tmpl, err := template.New("shadowenvLine").Parse(envValueTemplate)
		// TODO: handle err gracefully
		if err != nil {
			panic(err)
		}

		// Compile line
		var b bytes.Buffer
		err = tmpl.Execute(&b, dataForSingleLines)
		if err != nil {
			panic(err)
		}

		compiledEnvVarMap[envName] = string(b.Bytes())
	}

	dataForEntireTemplate := struct {
		Version int
		EnvVars map[string]string
	}{
		Version: g.Version,
		EnvVars: compiledEnvVarMap,
	}

	// FIXME: Use `brew --prefix …` instead of hardcoding the path
	templateContent := `(provide "postgresql" "{{ .Version }}")
{{ range $key, $value := .EnvVars }}
(env/set "{{ $key }}" "{{ $value }}")
{{- end }}

(env/prepend-to-pathlist "PATH" "/opt/homebrew/opt/postgresql@{{ .Version }}/bin")
`

	tmpl, err := template.New("shadowenv").Parse(templateContent)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, dataForEntireTemplate)
	if err != nil {
		panic(err)
	}

	return b.Bytes()
}
