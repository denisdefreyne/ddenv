# ddenv

**ddenv** (_Denis’ Developer Environment_) is a tool for setting up a local development environment.

Watch the 75-second introduction:

[![YouTube video: “An introduction to ddenv”](https://img.youtube.com/vi/SVd8HdS02yY/0.jpg)](https://www.youtube.com/watch?v=SVd8HdS02yY)

## Requirements

First, ensure you have [Homebrew](https://brew.sh/) installed, and are using a supported shell (zsh, bash, or fish).

Then, run `brew install denisdefreyne/tap/ddenv`.

## Quick start

Create a `.config/ddenv.yaml` or `ddenv.yaml` file which contains the list of dependencies to manage. For example:[^ruby-and-node-version]

```yaml
up:
  - homebrew: overmind
  - postgresql:
      version: 17
      env:
        DB_URL: "postgres://{{ .User }}:{{ .Password }}@{{ .Host }}:{{ .Port }}/mydb"
  - redis:
      env:
        REDIS_URL: "redis://{{ .Host }}:{{ .Port }}/0"
  - ruby
  - bundle
  - node
  - npm
```

[^ruby-and-node-version]: This example relies on a `.ruby-version` file being present, e.g. with the file contents `3.3.6`, and a `.node-version` file, e.g. with the file contents `20.12.2`.

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

-   <code>homebrew: <var>PACKAGENAME</var></code> installs the Homebrew package with the given name. Example:

    ```yaml
    up:
      - homebrew: overmind
    ```

-   <code>ruby</code> installs Ruby (with the version specified in the `.ruby-version` file). Example:

    ```yaml
    up:
      - ruby
    ```

-   <code>bundle</code> runs `bundle install`. Example:

    ```yaml
    up:
      - ruby
      - bundle
    ```

-   <code>node</code> installs Node.js (with the version specified in the `.node-version` file). Example:

    ```yaml
    up:
      - node
    ```

-   <code>npm</code> installs packages from package.json using npm. Example:

    ```yaml
    up:
      - node
      - npm
    ```

-   <code>pnpm</code> installs packages from package.json using pnpm. Example:

    ```yaml
    up:
      - node
      - pnpm
    ```

-   <code>yarn</code> installs packages from package.json using yarn. Example:

    ```yaml
    up:
      - node
      - yarn
    ```

-   <code>postgresql: <var>…</var></code> installs the given version of PostgreSQL (the <var>version</var> key), starts it, and sets up environment variables based on the <var>env</var> key (`User`, `Password`, `Host` and `Port` are available as keys). Example:

    ```yaml
    up:
      - postgresql:
          version: 17
          env:
            DB_URL: "postgres://{{ .User }}:{{ .Password }}@{{ .Host }}:{{ .Port }}/mydb"
    ```

-   <code>redis: <var>…</var></code> installs Redis, starts it, and sets up environment variables based on the <var>env</var> key (`Host` and `Port` are available as keys). Example:

    ```yaml
    up:
      - redis:
          env:
            REDIS_URL: "redis://{{ .Host }}:{{ .Port }}/0"
    ```

## License

Copyright (C) 2025–… Denis Defreyne

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

## Maintenance terms

Current as of 2025.

The code, text and data in this repository are provided as-is under the terms of the repository's LICENSE.txt file, as a gift to the commons and the common good. In providing this software as-is, its author admit no further obligations from anyone using the software for any reason, particularly with respect to:

- Response time,
- Change review and integration,
- Disclosure schedules,
- Discretionary, proprietary or otherwise secretive communications, and
- Any other non-contractual obligations or conventions, regardless of their presumed urgency or severity.

The author hope you find it valuable on those terms.

(Adapted from https://github.com/mhoye/maintenance-terms.)
