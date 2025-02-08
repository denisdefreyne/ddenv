package goals

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"denisdefreyne.com/x/ddenv/core"
)

const RedisShadowenvCreated_Path = ".shadowenv.d/400_redis.lisp"

type RedisShadowenvCreated struct {
	Env map[string]string
}

func (g RedisShadowenvCreated) Description() string {
	return fmt.Sprintf("Adding Redis to Shadowenv")
}

func (g RedisShadowenvCreated) HashIdentity() string {
	return fmt.Sprintf("RedisShadowenvCreated %v", g)
}

func (g RedisShadowenvCreated) IsAchieved() bool {
	_, err := os.Lstat(RedisShadowenvCreated_Path)
	if err != nil {
		return false
	}

	oldContents, err := os.ReadFile(RedisShadowenvCreated_Path)
	if err != nil {
		return false
	}

	if !bytes.Equal(oldContents, g.fileContents()) {
		return false
	}

	return true
}

func (g RedisShadowenvCreated) Achieve() error {
	err := os.WriteFile(RedisShadowenvCreated_Path, g.fileContents(), 0755)
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

func (g RedisShadowenvCreated) SubGoals() []core.Goal {
	return []core.Goal{
		ShadowenvSetUp{},
	}
}

func (g RedisShadowenvCreated) fileContents() []byte {
	dataForSingleLines := struct {
		Host string
		Port int16
	}{
		Host: "localhost",
		Port: 6379,
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
		EnvVars map[string]string
	}{
		EnvVars: compiledEnvVarMap,
	}

	// FIXME: Use `brew --prefix …` instead of hardcoding the path
	templateContent := `(provide "redis")
{{ range $key, $value := .EnvVars }}
(env/set "{{ $key }}" "{{ $value }}")
{{- end }}
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
