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
}

func (g RedisShadowenvCreated) Description() string {
	return fmt.Sprintf("Adding Redis to Shadowenv")
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

func (g RedisShadowenvCreated) PreGoals() []core.Goal {
	return []core.Goal{
		ShadowenvDirCreated{},
	}
}

func (g RedisShadowenvCreated) fileContents() []byte {
	// TODO: Allow configuring the database prefix (not just REDIS_)

	data := struct {
		EnvVarPrefix string
	}{EnvVarPrefix: "REDIS"}

	templateContent := `(provide "redis")

(env/set "{{ .EnvVarPrefix }}_URL" "redis://@127.0.0.1:6379/0")
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
