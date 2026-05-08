# NPKM (Nicolas's Playbook Kit Manager)

NPKM is a lightweight, declarative automation and provisioning tool (similar to Ansible or Chef), designed for zero-friction environment bootstrapping. It is written natively in the **Coni** programming language, featuring a custom YAML-to-EDN parser and cross-platform native execution.

## Core Features

- **Cross-OS Build**: Compiles entirely to standalone native binaries (`.exe` and `Mach-O`).
- **YAML Support**: Natively transforms Ansible-style tasks via its zero-dependency `yaml-to-edn` parser.
- **Remote HTTP Playbooks**: Can run playbooks directly via URL.
- **Git Repositories**: Scans cloned repos for playbook yaml/edn (`git clone`).
- **Directory Scanning**: Recursively lists available playbook files.
- **Global Configs**: Interpolation from `config:` blocks into `config.*` variables.
- **Remote SSH Orchestration**: Embedded SSH client allows running playbooks on remote hosts via `inventory.yml`.
- **Conditional Execution**: Support for `when` clauses to target specific OS platforms or custom conditions.

## Supported Tasks

| Task | Description |
| :--- | :--- |
| `file` | directory, touch, link, absent, modes |
| `lineinfile` | Regex matching & replacement in streams |
| `replace` | Replaces all instances of a regex pattern |
| `path` | Modifies the system PATH environment variable |
| `systemd` | start, stop, restart daemons |
| `copy`, `move`, `remove` | Standard IO primitives |
| `get_url` / `unzip` | Downloads and extracts remote assets |
| `shell`, `command`, `powershell`| Shell integration along with inline Powershell |
| `debug`, `fail` | Playbook execution logic and output |
| `package` | Auto-detects brew, apt-get, yum, winget, or choco |
| `service` | Generalizes systemctl, launchctl, and net start |
| `cron` | UNIX crontab -l / - insertion & absent state |
| `user` | Integrates useradd, sysadminctl, net user |
| `archive` | Native `zip` operations without shell dependencies |
| `template` | Deploy templated files with mapped configuration properties |
| `include_tasks` | Include & execute tasks from a local file, directory, or git repo |

## Task Reference & Examples

### `file`
Manage the state of a file, directory, or symlink.
```yaml
- name: Ensure configuration directory exists
  file:
    path: /etc/myapp
    state: directory
    mode: 0755
```

### `copy`
Copy an existing file or directory directly to a specified path.
```yaml
- name: Copy deployment artifact
  copy:
    src: ./build/app.jar
    dest: /opt/myapp/app.jar
```

### `move` / `remove`
Rename, move, or completely delete elements on the disk.
```yaml
- name: Rename old log
  move:
    src: /var/log/app.log
    dest: /var/log/app.old.log

- name: Wipe temporary backups
  remove:
    path: /tmp/backups/*
```

### `get_url` & `unzip`
Download remote assets and seamlessly extract them to the system.
```yaml
- name: Download web app
  get_url:
    url: https://github.com/user/repo/archive/main.zip
    dest: /tmp/app.zip

- name: Extract zip archive
  unzip:
    src: /tmp/app.zip
    dest: /var/www/html/
```

### `archive`
Compress local paths natively into an archive (without shell tools).
```yaml
- name: Backup web directory
  archive:
    src: /var/www/html/
    dest: /backups/html_backup.zip
```

### `package`
Automatically manage OS packages. Will intelligently resolve `brew`, `apt-get`, `yum`, `winget`, or `choco` depending on the platform.
```yaml
- name: Install Git
  package:
    name: git
    state: present
```

### `service` & `systemd`
Manage system-level daemons natively (`systemctl`, `launchctl`, or `net start`).
```yaml
- name: Enable and start Nginx
  service:
    name: nginx
    state: started
    enabled: true

- name: Stop multiple units simultaneously (e.g., to prevent socket activation warnings)
  systemd:
    name: syslog.socket rsyslog.service
    state: stopped
```

### `shell`, `command` & `powershell`
Execute raw OS-dependent instructions.
```yaml
- name: Run raw bash script
  shell:
    cmd: "rm -rf /tmp/cache && echo 'Cleared'"
    cwd: /tmp/

- name: Run Windows powershell instruction
  powershell:
    inline: "Get-Process | Where-Object {$_.Name -eq 'node'} | Stop-Process"
```

### `lineinfile` & `replace`
Modify and parse file streams based on regex.
```yaml
- name: Ensure memory limit is correct
  lineinfile:
    path: /etc/php.ini
    regexp: "^memory_limit="
    line: "memory_limit=512M"

- name: Swap default port anywhere in config
  replace:
    path: /opt/app/config.json
    regexp: "8080"
    replace: "9000"
```

### `path`
Append a directory natively to the global OS `$PATH` configuration.
```yaml
- name: Install java to path
  path:
    path: /opt/java/bin
```

### `user` & `cron`
Manage system-level profiles and periodic tasks.
```yaml
- name: Add worker user
  user:
    name: worker
    state: present

- name: Setup midnight backup
  cron:
    name: "DB Backup"
    state: present
    job: "0 0 * * * /opt/backup.sh"
```

### `debug` & `fail`
Provide real-time execution outputs or forcefully term execution conditions.
```yaml
- name: Print variables
  debug:
    msg: "Current root path is {{ config.root }}"

- name: Stop on unsupported OS
  fail:
    msg: "Halting execution: OS not supported."
```

### `include_tasks`
Dynamically include a list of tasks from a separate `.yml` file, a local directory (first `.yml` found), or a remote git repository. Combine with `when:` to load tasks conditionally.

**Local file:**
```yaml
tasks:
  - name: Include web server setup
    include_tasks: tasks/web_tasks.yml
    when: "ansible_os_family == 'Unix'"
```

**Local directory (first `.yml` file is used):**
```yaml
tasks:
  - name: Include all tasks in the db folder
    include_tasks: tasks/database/
```

**Remote git repository:**
```yaml
tasks:
  - name: Pull shared tasks from private repo
    include_tasks: git@github.com:myorg/common-tasks.git
    when: "env == 'production'"
```

The included file must be a flat YAML list of tasks (no `hosts:` or `plays:` wrapping):
```yaml
# web_tasks.yml
- name: Install nginx
  package:
    name: nginx
    state: present

- name: Start nginx
  service:
    name: nginx
    state: started
```

## Global Configuration Interpolation

NPKM supports dynamic global string replacement. You can define variables in an inline `config:` block at the top of your playbook (or placed alongside it as a separate `config.yml`), and they will be injected wherever `config.your_key` is referenced in the tasks.

```yaml
config:
  deploy_path: /opt/production
  service_user: nginx

tasks:
  - name: Ensure deployment directory exists
    file:
      path: config.deploy_path
      state: directory
```

## Conditional Execution (OS Detection)

NPKM provides built-in conditional execution using the `when:` clause. It automatically populates the `ansible_os_family` runtime variable (`Unix` or `Windows`) for both local and remote executions.

```yaml
tasks:
  - name: Install dependencies on Linux/macOS
    shell:
      cmd: curl -fsSL https://example.com/install.sh | sh
    when: "ansible_os_family == 'Unix'"

  - name: Install dependencies on Windows
    powershell:
      inline: irm https://example.com/install.ps1 | iex
    when: "ansible_os_family == 'Windows'"
```

## Privilege Escalation (become / sudo)

If a task requires root privileges on a Linux or macOS target (e.g., restarting a system daemon or installing a package), you can use the `become: true` flag. This will automatically prefix the command with `sudo`.

```yaml
tasks:
  - name: Restart rsyslog using systemd
    become: true
    systemd:
      name: rsyslog
      state: restarted
      enabled: true
```

**Note on passwords:** NPKM currently executes SSH commands non-interactively and does not pause to prompt for a sudo password. If your remote user requires a password to use `sudo`, the command will fail. To use `become: true`, you must configure your target machine's `/etc/sudoers` file to allow passwordless sudo for the user (e.g., `ubuntu ALL=(ALL) NOPASSWD:ALL`).

## Remote SSH Orchestration (Inventories)

NPKM allows you to execute your playbooks seamlessly over SSH to remote targets using an `inventory.yml` file. Just provide the inventory alongside your playbook!

```yaml
# inventory.yml
all:
  hosts:
    server1:
      ansible_host: 192.168.1.10
      ansible_user: root
      ansible_ssh_pass: "mysecret"                   # Optional: Password authentication
      ansible_ssh_private_key_file: "~/.ssh/id_rsa"  # Optional: SSH Key authentication
      ansible_port: 22
```

In your playbook, define `hosts: all` or explicitly target `hosts: server1`:

```yaml
# playbook.yml
name: Deploy Web Server
hosts: server1

tasks:
  - name: Install nginx
    package:
      name: nginx
      state: present
```

Execute by passing the inventory file using the `-i` flag to run via SSH:
```bash
# Run a playbook on remote hosts via SSH
./npkm-coni -i inventory.yml playbook.yml

# Example: Run the bundled install_ollama.yml on your remote SSH inventory
./npkm-coni -i inventory.yml install_ollama.yml
```

## Advanced Features

### Loops & Iteration
NPKM supports native task iteration using `with_items` and `loop` constructs. You can loop over inline lists or variables defined in your configuration, and dynamically interpolate the `{{ item }}` reference throughout your task properties.

**Using `with_items` (Inline List):**
```yaml
tasks:
  - name: Install required packages
    package:
      name: "{{ item }}"
      state: present
    with_items:
      - curl
      - git
      - docker
```

**Using `loop` (Variable Reference):**
```yaml
config:
  app_files:
    - index.html
    - app.js
    - style.css

tasks:
  - name: Copy app files
    copy:
      src: "./src/{{ item }}"
      dest: "/var/www/html/{{ item }}"
    loop: config.app_files
```

### Advanced Templating & Nesting
The YAML parser perfectly maps complex YAML structures into nested dictionaries. You can use the `template` task to inject a full dictionary of key-value pairs (using the `vars:` map) into your configuration templates seamlessly:

```yaml
tasks:
  - name: Configure Nginx Site
    template:
      src: ./templates/nginx.conf.j2
      dest: /etc/nginx/nginx.conf
      vars:
        port: 8080
        server_name: mysite.local
        worker_processes: 4
```

# Usage

Provide a single local YAML/EDN file, a directory containing playbooks, a mix of files and folders, a remote HTTP/HTTPS link, or an SSH/Git path. When you pass a directory, NPKM recursively lists and evaluates all playbook files inside it!

```bash
# Run a specific local playbook
./npkm-coni test-playbook.yml

# Run all playbooks inside a directory
./npkm-coni ./playbooks/

# Mix and match individual files and folders at the same time
./npkm-coni deploy-web.yml ./database_setup/ ./monitoring/

# Clone from Git and run
./npkm-coni ssh://git@s5:2222/hellonico/my-playbook.git

# Run directly from a remote web server
./npkm-coni https://raw.githubusercontent.com/user/npkm/main/playbook.yml
```

## Documentation Generation

You can automatically generate Markdown documentation with Mermaid graphs for your playbooks and inventory using the `--doc` flag.

```bash
# Generate documentation for a playbook and print to stdout
./npkm-coni --doc test-playbook.yml

# Generate documentation for multiple playbooks with an inventory and save to a file
./npkm-coni -i inventory.yml --doc web.yml db.yml > doc.md
```

## Automatic Logging

NPKM-Coni automatically records and archives the output of every playbook execution natively! 

Every time you run the tool, your complete execution trace is intercepted in the background. Once the run finishes (or upon failure), the logs are automatically stripped of ANSI color codes and saved as a plain-text log inside your local `~/.npkm/` directory.

- **Log Path Format:** `~/.npkm/YYYY-MM-DD_HH-MM-SS.log`
- **Clean output:** The log preserves all standard output minus the terminal color formatting for perfect readability in text editors.
