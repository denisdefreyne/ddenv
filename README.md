# ddenv

**ddenv** (*Denis’ Developer Environment*) is a tool for maintaining a local environment for development.

> [!CAUTION]
> This is pre-release software.

## Requirements

First, ensure you have [Homebrew](https://brew.sh/) installed, and are using a supported shell (zsh, bash, or fish).

Then, run `brew install denisdefreyne/tap/ddenv`.

## Quick start

Create a `ddenv.yaml` file which contains the list of dependencies to manage. For example:[^ruby-version]

```yaml
up:
  - homebrew: overmind
  - ruby
  - bundle
  - node: 20.12.2
  - npm
```

[^ruby-version]: This example relies on a `.ruby-version` file being present, e.g. with the file contents `3.3.6`.

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

-   <code>homebrew: <var>PACKAGENAME</var></code> installs the Homebrew package with the given name.

    e.g. `- homebrew: overmind`

-   <code>ruby</code> installs Ruby (with the version specified in the `.ruby-version` file).

    e.g. `- ruby`

-   <code>bundle</code> runs `bundle install`.

    e.g. `- bundle`

-   <code>node: <var>VERSION</var></code> installs the give Node.js version.

    e.g. `- node: 22.11.0`

-   <code>npm</code> installs packages from package.json using npm.

    e.g. `- npm`
