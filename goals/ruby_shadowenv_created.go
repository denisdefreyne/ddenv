package goals

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"denisdefreyne.com/x/ddenv/core"
)

const RubyShadowenvCreated_Path = ".shadowenv.d/100_ruby.lisp"

type RubyShadowenvCreated struct {
	Version string
	Path    string
}

func (g RubyShadowenvCreated) Description() string {
	return fmt.Sprintf("Adding Ruby %v to Shadowenv", g.Version)
}

func (g RubyShadowenvCreated) IsAchieved() bool {
	_, err := os.Lstat(RubyShadowenvCreated_Path)
	if err != nil {
		return false
	}

	oldContents, err := os.ReadFile(RubyShadowenvCreated_Path)
	if err != nil {
		return false
	}

	if !bytes.Equal(oldContents, g.fileContents()) {
		return false
	}

	return true
}

func (g RubyShadowenvCreated) Achieve() error {
	err := os.WriteFile(RubyShadowenvCreated_Path, g.fileContents(), 0755)
	if err != nil {
		return err
	}

	shadowenvTrustCmd := exec.Command("shadowenv", "trust")
	if err := shadowenvTrustCmd.Run(); err != nil {
		return err
	}

	return nil
}

func (g RubyShadowenvCreated) PreGoals() []core.Goal {
	return []core.Goal{
		ShadowenvDirCreated{},
	}
}

func (g RubyShadowenvCreated) fileContents() []byte {
	data := struct {
		RubyVersion string
		RubyPath    string
	}{RubyVersion: g.Version, RubyPath: g.Path}

	templateContent := `(provide "ruby" "{{ .RubyVersion }}")

(when-let ((ruby-root (env/get "RUBY_ROOT")))
(env/remove-from-pathlist "PATH" (path-concat ruby-root "bin"))
(when-let ((gem-root (env/get "GEM_ROOT")))
	(env/remove-from-pathlist "PATH" (path-concat gem-root "bin")))
(when-let ((gem-home (env/get "GEM_HOME")))
	(env/remove-from-pathlist "PATH" (path-concat gem-home "bin"))))

(env/set "GEM_PATH" ())
(env/set "GEM_HOME" ())
(env/set "RUBYOPT" ())

(env/set "RUBY_ROOT" "{{ .RubyPath }}")
(env/prepend-to-pathlist "PATH" "{{ .RubyPath }}/bin")
(env/set "RUBY_ENGINE" "ruby")
(env/set "RUBY_VERSION" "{{ .RubyVersion }}")
(env/set "GEM_ROOT" "{{ .RubyPath }}/lib/ruby/gems/{{ .RubyVersion }}")

(when-let ((gem-root (env/get "GEM_ROOT")))
	(env/prepend-to-pathlist "GEM_PATH" gem-root)
	(env/prepend-to-pathlist "PATH" (path-concat gem-root "bin")))

(let ((gem-home
			(path-concat (env/get "HOME") ".gem" (env/get "RUBY_ENGINE") (env/get "RUBY_VERSION"))))
	(do
		(env/set "GEM_HOME" gem-home)
		(env/prepend-to-pathlist "GEM_PATH" gem-home)
		(env/prepend-to-pathlist "PATH" (path-concat gem-home "bin"))))`

	tmpl, err := template.New("test").Parse(templateContent)
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
