(provide "ruby" "3.3.1")

(when-let ((ruby-root (env/get "RUBY_ROOT")))
(env/remove-from-pathlist "PATH" (path-concat ruby-root "bin"))
(when-let ((gem-root (env/get "GEM_ROOT")))
	(env/remove-from-pathlist "PATH" (path-concat gem-root "bin")))
(when-let ((gem-home (env/get "GEM_HOME")))
	(env/remove-from-pathlist "PATH" (path-concat gem-home "bin"))))

(env/set "GEM_PATH" ())
(env/set "GEM_HOME" ())
(env/set "RUBYOPT" ())

(env/set "RUBY_ROOT" "/Users/denis/.rubies/ruby-3.3.1")
(env/prepend-to-pathlist "PATH" "/Users/denis/.rubies/ruby-3.3.1/bin")
(env/set "RUBY_ENGINE" "ruby")
(env/set "RUBY_VERSION" "3.3.1")
(env/set "GEM_ROOT" "/Users/denis/.rubies/ruby-3.3.1/lib/ruby/gems/3.3.1")

(when-let ((gem-root (env/get "GEM_ROOT")))
	(env/prepend-to-pathlist "GEM_PATH" gem-root)
	(env/prepend-to-pathlist "PATH" (path-concat gem-root "bin")))

(let ((gem-home
			(path-concat (env/get "HOME") ".gem" (env/get "RUBY_ENGINE") (env/get "RUBY_VERSION"))))
	(do
		(env/set "GEM_HOME" gem-home)
		(env/prepend-to-pathlist "GEM_PATH" gem-home)
		(env/prepend-to-pathlist "PATH" (path-concat gem-home "bin"))))