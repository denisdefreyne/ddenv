# ddenv

**ddenv** (*Denis’ Developer Environment*) is a tool for maintaining a local
environment for development.

> [!CAUTION]
> This software is a pre-alpha work in progress. Do not use just yet.

## Requirements

- [Homebrew](https://brew.sh/)
- a support shell (zsh, bash, or fish)

## Installation

Run `brew install denisdefreyne/tap/ddenv`.

## Quick start

Create a `ddenv.yaml` file which contains the list of dependencies to manage:

```yaml
up:
  - homebrew: overmind
  - ruby: 3.3.0
  - bundle
  - node: 20.12.2
  - npm
```

Then, run `ddenv`:

```
% ddenv
Installing Homebrew package ‘overmind’       skipped
Installing Homebrew package ‘ruby-install’   skipped
Installing Ruby 3.3.6                        skipped
Adding Shadowenv to shell                    skipped
Creating Shadowenv dir                       done
Adding Shadowenv dir to .gitignore           done
Adding Ruby 3.3.6 to Shadowenv               done
Installing Ruby gem bundler                  done
Installing bundle                            done
Installing Homebrew package ‘node-build’     checking...
Installing Node 20.12.2                      pending
Adding Node 20.12.2 to Shadowenv             pending
Installing npm packages                      pending
```

Now your local developer environment is ready to be used.

## Goals

- <code>homebrew: <var>PACKAGENAME</var></code> installs the Homebrew package with the given name.
- <code>ruby</code> installs Ruby (with the version specified in the `.ruby-version`
  file).
- <code>bundle</code> runs `bundle install`.
- <code>node: <var>VERSION</var></code> installs the give Node.js version.
- <code>npm</code> installs packages from package.json using npm.
