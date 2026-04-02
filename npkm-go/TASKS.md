# npkm-go Tasks Overview

This document describes the tasks available in the `npkm-go` playbook runner. The tasks ported from the previous `coni` version include all common system, file manipulation, and Git management actions.

## Task Reference Table

| Task | Description | Fields | Example |
|------|-------------|--------|---------|
| `shell` | Execute a shell command string | `cmd`<br>`cwd` (optional) | `- shell: { cmd: "echo $USER" }` |
| `file` | Manage files and directories (create, symlink, touch, remove) | `path`<br>`state` (directory, touch, link, absent)<br>`src` (for link)<br>`mode` (optional) | `- file: { path: "/tmp/foo", state: "directory" }` |
| `debug` | Print a debug message to standard output | `msg` | `- debug: { msg: "Hello World" }` |
| `copy` | Copy a file from a local source path to a destination path | `src`<br>`dest` | `- copy: { src: "./file.txt", dest: "/opt/file.txt" }` |
| `remove`| Completely delete a file or directory tree | `path` | `- remove: { path: "/tmp/old_dir" }` |
| `fail` | Abort playbook execution with a custom error message | `msg` | `- fail: { msg: "Pre-condition failed!" }` |
| `unzip` | Extract a zip archive to a destination directory | `src`<br>`dest` | `- unzip: { src: "archive.zip", dest: "/tmp" }` |
| `git` | Clone or pull a remote git repository | `repo`<br>`dest` | `- git: { repo: "https://gitea/r.git", dest: "./opt" }` |
| `move` | Move or rename a file (with cross-device fallback) | `src`<br>`dest` | `- move: { src: "/tmp/a.txt", dest: "/tmp/b.txt" }` |
| `path` | Persistently append a new path to the user's PATH (supports Windows, macOS, Linux) | `path` | `- path: { path: "/opt/bin/custom" }` |

### Other Built-in Tasks

| Task | Description | Fields | Example |
|------|-------------|--------|---------|
| `command` | Execute a command directly without invoking a shell | `cmd`<br>`cwd` (optional) | `- command: { cmd: "ls -la" }` |
| `get_url` | Download a file via HTTP/HTTPS | `url`<br>`dest` | `- get_url: { url: "http://..", dest: "./out" }` |
| `lineinfile` | Ensure a specific line exists in a file (with optional regex substitution) | `path`<br>`line`<br>`regexp` (optional) | `- lineinfile: { path: "/etc/hosts", line: "127.0.0.1 db" }` |
| `replace` | Find and replace text directly within a file using RegEx | `path`<br>`regexp`<br>`replace` | `- replace: { path: "conf", regexp: "foo", replace: "bar" }` |
| `systemd` | Manage systemd services | `name`<br>`state`<br>`enabled` | `- systemd: { name: "nginx", state: "restarted", enabled: true }` |
