# NPKM — Nuke Playbook Kit Manager

> A native, zero-dependency automation engine written in **Coni**. Deploy, provision, and orchestrate infrastructure with full Ansible parity — and capabilities beyond it.

---

## Release History

### v2.0 "Novae" _(Latest)_
- **[`set_fact` runtime variables](#set_fact)**: Assign variables in one task and reference them with `${var}` in any subsequent task
- **[`register` output capture]**: Save any module's execution output (including stdout/stderr) to a variable for subsequent tasks.
- **Host Filtering**: Use `--limit <host_or_group>` to surgically target specific infrastructure subsets.
- **Config seeding**: All `config:` block keys are automatically available as `${key}` throughout the playbook — no `set_fact` needed
- **Variable chaining**: `set_fact` values can themselves reference earlier `${vars}`, enabling derived variables
- **Mid-playbook overrides**: Call `set_fact` again at any point to update a variable for all following tasks
- **Universal interpolation**: `${var}` works in every string field across all modules (`shell.cmd`, `file.path`, `debug.msg`, `archive.src/dest`, etc.)
- **Enhanced Modules**: 
  - `stat`: Fetch rich file/directory telemetry into nested maps (`{{ file_info.stat.size }}`).
  - `copy`: Now supports `content` mode to write templated strings directly to disk.
  - **Native OS Package Aliases**: Use direct `apt:`, `yum:`, `brew:`, `winget:`, and `choco:` module syntax.
  - **Dry-run (`--check`)**: `copy`, `file`, and `remove` now cleanly simulate their execution without mutating disk state.

### v1.6 "Sentinel"
- **[Role Package Manager](#roles--package-manager)**: Install reusable automation roles from any Git repository with `npkm roles install`
- **[Project Scaffolding](#project-scaffolding-npkm-init)**: Scaffold a complete project skeleton with `npkm init`
- **[Static Analysis](#static-analysis-npkm-lint)**: Validate playbooks before running with `npkm lint`
- **[Watch Mode](#watch-mode-npkm-watch)**: Auto re-run playbooks on file change with `npkm watch`
- **[Interactive Step Mode](#interactive-step-mode---step)**: Execute tasks one-by-one with confirmation via `--step`
- **[Execution Reports](#execution-reports---report)**: Generate JSON + HTML audit reports via `--report`
- **[Run History](#run-history)**: Browse and diff past execution logs with `npkm run history`
- **Keyword var interpolation**: `:vars {:key val}` in `include_tasks` now correctly resolves `{{ key }}` templates
- **Multi-line command safety**: SSH commands with `&&` in block scalars now execute correctly on Debian/Ubuntu (`dash`)

### v1.5 "Quantum Weaver"
- Native Templating (Variables & Loops), Multi-Play Architecture, Documentation Generation (`--doc`), Task Filtering (`--labels`, `--names`), Background Logging

### v1.4 "Flow Control"
- `block` / `rescue` / `always`, Handlers & Notifications, Parallel Host Execution (`forks`)

---

## Core Features

- **Cross-platform binary**: Single static binary for macOS, Linux, and Windows — no Python, JVM, or runtime required
- **YAML + EDN**: Full Ansible-style YAML support alongside native EDN format
- **SSH orchestration**: Built-in SSH client for remote host execution
- **Vault encryption**: AES-256-CBC file encryption with transparent runtime decryption
- **Dynamic inventory**: Executable scripts auto-detected alongside static YAML/EDN/INI inventories
- **Role system**: Reusable, Git-versioned automation modules
- **Zero dependencies**: No pip install, no requirements.txt, no Galaxy account

---

## Quick Start

```bash
# Run a playbook locally
npkm playbook.yml

# Run against remote hosts over SSH
npkm -i inventory.yml playbook.yml

# Scaffold a new project
npkm init my-project/

# Validate before running
npkm lint playbook.yml

# Watch for changes and re-run automatically
npkm watch -i inventory.yml playbook.yml
```

---

## Examples (v2.0 Features)

Here is a quick playbook showcasing the latest module improvements, output capturing (`register`), nested variable interpolation, and dry-run safety:

```yaml
- name: Setup Web Server
  hosts: all
  tasks:
    - name: Fetch details about the existing nginx directory
      stat:
        path: /etc/nginx
      register: nginx_stat

    - name: Print the directory size if it exists
      debug:
        msg: "Nginx config size is {{ nginx_stat.stat.size }} bytes"
      # Conditionally runs only if the nested map evaluation is true
      when: "{{ nginx_stat.stat.exists }}"

    - name: Ensure Nginx is installed (using native OS alias)
      apt:
        name: nginx
        state: present

    - name: Write a templated index file directly to disk
      copy:
        dest: /var/www/html/index.html
        content: |
          <h1>Welcome to {{ hostname }}</h1>
          <p>Managed natively by NPKM</p>
```

**Running with `--check` (Dry Run):**
If you run the above playbook with `npkm --check playbook.yml`, the `apt` and `copy` modules will gracefully simulate execution and return `changed: true` without altering your server state!

**Running with `--limit`:**
You can seamlessly restrict `hosts: all` to a specific target subset:
```bash
npkm --limit web_servers playbook.yml
```

---

## Roles — Package Manager

Roles are reusable, Git-versioned task collections. Install them from any Git repository and reference them in your playbooks via `include_tasks`.

### Installing a role

```bash
# Install from a Git repo — cloned into ~/.npkm/roles/<repo-name>/
npkm roles install git@github.com:myorg/nginx-role.git

# Install a specific version (tag or branch)
npkm roles install git@gitlab.example.com:sys/binet.git --version v1.2.0
```

Roles are stored in `~/.npkm/roles/`. Each role follows this layout:

```
~/.npkm/roles/
  nginx-role/
    tasks/
      main.edn       ← entry point (flat list of tasks)
    defaults/
      main.edn       ← default variable values
```

### Using a role in a playbook

Reference an installed role with `include_tasks:` pointing to the role name under `roles/`:

```yaml
# smb_share.yml
- name: Setup Samba share
  hosts: biner3
  tasks:
    - name: Install and configure Samba
      include_tasks: roles/samba
      vars:
        share_name:  "MY_SHARE"
        share_path:  "/mnt/data/samba/my_share"
        smb_user:    "alice"
        smb_comment: "Production data share"
```

Or in EDN format:

```edn
{:name "Setup Samba share on biner3"
 :hosts "biner3"
 :tasks [{:name "Install and configure Samba"
          :include_tasks "roles/samba"
          :vars {:share_name  "MY_SHARE"
                 :share_path  "/mnt/data/samba/my_share"
                 :smb_user    "alice"
                 :smb_comment "Production data share"}}]}
```

### Role defaults

Variables defined in `defaults/main.edn` act as fallbacks — overridden by anything passed in `:vars`:

```edn
; defaults/main.edn
{:share_name  "DEFAULT_SHARE"
 :smb_user    "guest"
 :smb_password "changeme"}
```

### Role task file format

`tasks/main.edn` must be a **flat vector of tasks** (no `:hosts` or play wrapping):

```edn
[
  {:name "Install samba" :become true :shell {:cmd "apt-get install -y samba"}}
  {:name "Start smbd"    :become true :systemd {:name "smbd" :state "restarted" :enabled true}}
]
```

---

## Project Scaffolding (`npkm init`)

Scaffold a ready-to-run project structure in one command:

```bash
npkm init my-project/
```

Creates:

```
my-project/
  main.edn          ← main playbook
  inventory.edn     ← host inventory
  group_vars/
    all.edn         ← shared variables
  tasks/
    setup.edn       ← example task file
  roles/            ← role directory
```

---

## Static Analysis (`npkm lint`)

Validate playbook structure before executing — catches missing required fields, unknown modules, and structural issues:

```bash
npkm lint playbook.yml
npkm lint smb_share.edn

# Example output:
# ⬡ Linting: smb_share.edn
# ✓ No issues found.
```

---

## Watch Mode (`npkm watch`)

Monitor your playbook and inventory files for changes and re-run automatically — ideal during active role or playbook development:

```bash
# Watch a playbook (re-runs on any file change)
npkm watch playbook.yml

# Watch with a remote inventory
npkm watch -i inventory.edn smb_share.edn

# Example output:
# ⬡ NPKM Watch Mode — watching: smb_share.edn, inventory.edn
#   Press Ctrl+C to stop.
#
# [watch] Change detected — re-running playbook... (run #1)
```

---

## Interactive Step Mode (`--step`)

Execute tasks one at a time with an interactive prompt — ideal for high-risk or first-time runs:

```bash
npkm --step -i inventory.yml deploy.yml
```

```
TASK [ Install nginx ]
  → Run this task? [y/n/q]:
```

- `y` — run the task and continue
- `n` — skip this task
- `q` — quit execution immediately

---

## Execution Reports (`--report`)

Generate a timestamped JSON + dark-themed HTML execution report in `~/.npkm/reports/` after every run:

```bash
npkm --report -i inventory.yml playbook.yml

# --- NPKM Run Report ---
#   ok=12 changed=4 failed=0 skipped=1 duration=8s
#   JSON: ~/.npkm/reports/2026-05-15_09-45-00.json
#   HTML: ~/.npkm/reports/2026-05-15_09-45-00.html
```

---

## Run History

Browse, inspect, and diff past execution logs stored in `~/.npkm/logs/`:

```bash
# List all past runs
npkm run history

# Show the most recent log
npkm run history last

# Diff the last two runs
npkm run history diff
```

---

## New Modules (v2.0 & v1.6)

### `set_fact`

Inject variables into the runtime environment mid-playbook. These variables are immediately available to all subsequent tasks using the new `${var}` or `{{ var }}` syntax.

You can even chain variables, referencing previously defined facts!

```yaml
- name: Compute paths
  set_fact:
    app_root: "/opt/myapp"
    log_dir:  "${app_root}/logs"

- name: Use the variable
  debug:
    msg: "App root is ${app_root} and logs go to ${log_dir}"
```

### `test`

Inline TDD-style assertions on task command output — fail fast if expectations aren't met:

```yaml
- name: Assert samba is running
  test:
    cmd: "systemctl is-active smbd"
    expect: "active"

- name: Assert share is accessible
  test:
    cmd: "smbclient -L localhost -N"
    contains: "MY_SHARE"
```

---

## Supported Modules

| Module | Description |
|---|---|
| `shell`, `command` | Execute shell commands |
| `powershell` | Windows PowerShell execution |
| `file` | Manage files, directories, symlinks |
| `copy`, `move`, `remove` | File I/O primitives |
| `lineinfile`, `replace` | Regex-based file modification |
| `template` | Render templated config files |
| `get_url` | Download remote files |
| `archive`, `unzip` | Compress / extract |
| `package` | Generic package manager abstraction |
| `apt`, `yum`, `brew`, `winget`, `choco` | OS-specific package manager native aliases |
| `service`, `systemd` | Manage system daemons |
| `user` | Create / remove system users |
| `cron` | Manage crontab entries |
| `stat` | Retrieve file or file system status |
| `git` | Clone or pull repositories |
| `path` | Modify `$PATH` |
| `debug`, `fail` | Output and control flow |
| `include_tasks` | Load tasks from file, directory, or Git |
| `block` / `rescue` / `always` | Error handling and cleanup |
| `coni` | Inline Coni scripts with full playbook context |
| `set_fact` | Inject runtime variables |
| `test` | Inline assertions on command output |

---


## Advanced Execution & Templating (v2.1)

### Task Delegation (`delegate_to`)
Execute a specific task on a different host than the one currently being provisioned, while still having access to the target's variables.

```yaml
- name: Remove from load balancer pool
  command: "haproxyctl disable server {{ inventory_hostname }}"
  delegate_to: load_balancer_01
```

### Asynchronous Tasks (`async` & `poll`)
Run long-running tasks in the background without blocking the rest of your playbook execution.

```yaml
- name: Run database migration
  shell:
    cmd: "rake db:migrate"
  async: 300     # Maximum time (in seconds) the task is allowed to run
  poll: 0        # 0 means "fire-and-forget" (don't wait for completion)
```

### Shell Idempotence (`creates` / `removes`)
Make shell commands perfectly idempotent (safe to run multiple times) by checking file existence.

```yaml
- name: Download application binary
  shell:
    cmd: "wget http://example.com/app -O /usr/local/bin/app"
    creates: "/usr/local/bin/app"  # Skip if file already exists

- name: Clean up temporary files
  shell:
    cmd: "rm -rf /tmp/build-cache"
    removes: "/tmp/build-cache"    # Skip if file is already removed
```

### Playbook Tags (`--tags` / `--skip-tags`)
Tag specific tasks and selectively run them.

```yaml
- name: Update database schema
  command: "migrate"
  tags: ["db", "upgrade"]

- name: Drop database
  command: "dropdb"
  tags: ["db", "destructive"]
```
```bash
npkm --tags db --skip-tags destructive playbook.yml
```

### Advanced Template Filters
Format, join, and manipulate variables directly inside templates!

```yaml
- name: Set facts
  set_fact:
    my_list: ["a", "b", "c"]
    my_var: ""

- name: Use inline filters
  debug:
    msg: "Joined list: {{ my_list | join(',') }} or Default var: {{ my_var | default('fallback') }}"
```

---

## Remote SSH Orchestration (Inventories)


```yaml
# inventory.yml
all:
  hosts:
    server1:
      ansible_host: 192.168.1.10
      ansible_user: ubuntu
      ansible_ssh_private_key_file: "~/.ssh/id_rsa"
      ansible_port: 22
```

```bash
npkm -i inventory.yml playbook.yml
```

---

## Flow Control & Error Handling

```yaml
tasks:
  - name: Risky operations
    block:
      - name: Download artifact
        get_url:
          url: "http://example.com/artifact"
          dest: "/tmp/artifact"
    rescue:
      - name: Use fallback
        shell:
          cmd: "echo 'fallback' > /tmp/artifact"
    always:
      - name: Cleanup
        debug:
          msg: "Run complete."
```

---

## Vault Encryption

Encrypt secrets at rest, decrypt transparently at runtime:

```bash
# Encrypt a file
npkm vault encrypt secrets.edn

# Decrypt for inspection
npkm vault decrypt secrets.edn.vault

# Runtime: set the password via environment variable
export NPKM_VAULT_PASSWORD=mysecret
npkm -i inventory.yml playbook.yml
```

---

## Documentation Generation

```bash
# Generate Mermaid flowchart + task table to stdout
npkm --doc playbook.yml

# Save to file
npkm -i inventory.yml --doc deploy.yml > docs/deploy.md
```

---

## Usage Reference

```bash
npkm [options] <playbook.yml | directory | https://... | git@...>

Options:
  -v                    print version
  -h                    show help
  --doc                 generate Mermaid documentation
  --dry-run, --check    simulate without making changes
  --diff                show file diffs
  --report              generate HTML + JSON execution report
  --step                interactive task-by-task confirmation
  --limit <hosts>       limit execution to specific hosts or groups
  --labels <csv>        run only tasks matching labels
  --names  <csv>        run only tasks matching names
  -i <file>             inventory file
  -bw                   disable color output

Commands:
  npkm init [dir]                scaffold a new project
  npkm doctor                    health check and system validation
  npkm lint <playbook>           static analysis
  npkm watch <playbook>          re-run on file change
  npkm run history               list past run logs
  npkm run history last          show most recent log
  npkm run history diff          diff last two runs
  npkm roles install <git-url>   install a role from Git
  npkm vault encrypt <file>      encrypt with AES-256
  npkm vault decrypt <file>      decrypt vault file
```

---

## Directory Layout

```
~/.npkm/
  logs/       ← timestamped execution logs (auto-created)
  reports/    ← JSON + HTML reports (--report)
  roles/      ← installed roles (npkm roles install)
```
