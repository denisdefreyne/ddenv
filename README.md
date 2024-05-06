# ddenv

**ddenv** (*Denis’ Developer Environment*) is a tool for maintaining a local
environment for development.

> [!CAUTION]
> This software is a pre-alpha work in progress. Do not use just yet.

## Requirements

- Homebrew
- fish shell

> [!NOTE]
> In the future, this will work for bash and zsh as well.

## Installation

> [!NOTE]
> At some point, it’ll be possible simply run `brew install denisdefreyne/ddenv/ddenv`.

1. Ensure you have Go installed (1.22 or later)
2. Clone this repository
3. In this repository, run `make install`

The `ddenv` executable will be placed inside `~/bin`.

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
Installing Homebrew package ‘ruby-install’      skipped
Installing Ruby 3.3.1                           skipped
Adding Shadowenv to shell                       skipped
Creating Shadowenv dir                          done
Adding Shadowenv dir to .gitignore              done
Adding Ruby 3.3.1 to Shadowenv                  done
Installing Ruby gem bundler                     skipped
Installing bundle                               skipped
Installing Homebrew package ‘node-build’        checking...
Installing Node 22.0.0                          pending
Adding Shadowenv to shell                       pending
Creating Shadowenv dir                          pending
Adding Shadowenv dir to .gitignore              pending
Adding Node 22.0.0 to Shadowenv                 pending
Installing npm packages                         pending
```

Now your local developer environment is ready to be used.

## Goals

- <code>homebrew: <var>PACKAGENAME</var></code> installs the Homebrew package with the given name.
- <code>ruby</code> installs Ruby (with the version specified in the `.ruby-version`
  file).
- <code>bundle</code> runs `bundle install`.
- <code>node: <var>VERSION</var></code> installs the give Node.js version.
- <code>npm</code> installs packages from package.json using npm.
