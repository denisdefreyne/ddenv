package goals

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"denisdefreyne.com/x/ddenv/core"
)

const NodeShadowenvCreated_Path = ".shadowenv.d/200_node.lisp"

type NodeShadowenvCreated struct {
	Version string
	Path    string
}

func (g NodeShadowenvCreated) Description() string {
	return fmt.Sprintf("Adding Node %v to Shadowenv", g.Version)
}

func (g NodeShadowenvCreated) HashIdentity() string {
	return fmt.Sprintf("NodeShadowenvCreated %v", g)
}

func (g NodeShadowenvCreated) IsAchieved() bool {
	_, err := os.Lstat(NodeShadowenvCreated_Path)
	if err != nil {
		return false
	}

	oldContents, err := os.ReadFile(NodeShadowenvCreated_Path)
	if err != nil {
		return false
	}

	if !bytes.Equal(oldContents, g.fileContents()) {
		return false
	}

	return true
}

func (g NodeShadowenvCreated) Achieve() error {
	err := os.WriteFile(NodeShadowenvCreated_Path, g.fileContents(), 0755)
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

func (g NodeShadowenvCreated) SubGoals() []core.Goal {
	return []core.Goal{
		ShadowenvSetUp{},
	}
}

func (g NodeShadowenvCreated) fileContents() []byte {
	data := struct {
		NodeVersion string
		NodePath    string
	}{NodeVersion: g.Version, NodePath: g.Path}

	templateContent := `(provide "node" "{{ .NodeVersion }}")

(env/prepend-to-pathlist "PATH" "{{ .NodePath }}/bin")`

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
