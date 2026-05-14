# NPKM Feature Reference

> **NPKM** — Nuke Playbook Manager  
> A native, zero-dependency automation engine written in Coni with full Ansible parity and unique capabilities beyond it.

---

## ✅ Core Execution Engine

| Feature | Detail |
|---|---|
| Shell/Command execution | `shell`, `command`, `powershell` |
| File management | `file`, `copy`, `move`, `remove`, `lineinfile`, `replace` |
| Templating | `{{ var }}` interpolation across tasks, vars, and templates |
| Static Inventory | YAML, EDN, INI, inline hosts |
| Dynamic Inventory | Executable scripts (JSON or YAML output auto-detected) |
| SSH remote execution | Native SSH via `sys-ssh-exec`, host vars, keys, passwords |
| Parallel host execution | `forks:` per-play parallelism via goroutine fan-out |
| Conditional execution | `when:` clauses with `==` / `!=` operators |
| Loops | `loop:`, `with_items:`, `items:` with `{{ item }}` replacement |
| Variable `register` | Capture task output into named variables |
| Error handling | `block:` / `rescue:` / `always:` structured error boundaries |
| Event triggers | `handlers:` + `notify:` with deduplication |
| Task retry | `retries:`, `until:`, `delay:` |
| Task inclusion | `include_tasks:` — local file, directory, or git URL |
| Role package manager | `npkm roles install <git-url> [version]` → `~/.npkm/roles/` |
| Vault encryption | AES-256-CBC via `npkm vault encrypt/decrypt`; transparent runtime decryption |
| Package management | `package:` — brew, apt-get, yum, winget, choco auto-detected |
| Service management | `service:`, `systemd:` — Linux, macOS, Windows |
| User management | `user:` — useradd/sysadminctl/net user |
| Cron management | `cron:` — idempotent via marker comments |
| HTTP download | `get_url:` |
| Git clone/pull | `git:` |
| Archive/unzip | `archive:`, `unzip:` |
| Dry-run mode | `--dry-run`, `--check` |
| File diff mode | `--diff` |
| Idempotent reporting | `ok`, `changed`, `skipped` per task |
| `become` (sudo) | `become: true` on any task or play |
| Cross-platform | macOS, Linux, Windows (PowerShell path) |

---

## ✅ Sprint 6 — Beyond Ansible

| Feature | Command / Usage |
|---|---|
| **`set_fact:`** | Set runtime variables mid-playbook: `:set_fact {:my_var "value" :count 42}` |
| **`test:` module** | Inline assertions: `:test {:cmd "echo hi" :expect "hi" :contains "..."}` |
| **`--step` mode** | Interactive task-by-task confirmation: `npkm --step playbook.edn` (y/n/q) |
| **`--report` flag** | Generates JSON + dark-themed HTML report in `~/.npkm/reports/` after every run |
| **`npkm init`** | Scaffolds a new project: `npkm init [dir]` creates `main.edn`, `inventory.edn`, `group_vars/`, `roles/`, `tasks/` |
| **`npkm lint`** | Static analysis before running: `npkm lint playbook.yml` — checks missing names, unknown modules, required fields |
| **`npkm run history`** | Browse past runs: `npkm run history` / `last` / `diff` |
| **`npkm watch`** | Re-runs playbook on file change: `npkm watch playbook.edn` (1s polling) |

---

## 🔥 NPKM-Unique Capabilities

| Feature | Detail |
|---|---|
| `--doc` Mermaid flowcharts | Generates visual playbook flow with `block/rescue/always` subgraphs |
| EDN format support | Tasks, vars, and inventory can be written in EDN or YAML |
| `coni:` inline scripting | Embed arbitrary Coni code as a task module |
| Native binary | Single static binary — no Python, no JVM, no runtime |
| Persistent run logs | All output captured to `~/.npkm/logs/` automatically |
| Label/name filtering | `--labels`, `--names` — run only specific tasks |
| HTML execution reports | Dark-themed, color-coded per-task run summaries |
| Inline test assertions | `test:` module for TDD-style playbook verification |
| Project scaffolding | `npkm init` — one command from zero to running |
| Interactive step mode | `--step` — surgical task-by-task execution with abort |
| Run history & diff | `npkm run history diff` — compare last two execution logs |

---

## 📁 Directory Layout

```
~/.npkm/
  logs/           # timestamped run logs (auto-migrated)
  reports/        # JSON + HTML execution reports (--report)
  roles/          # installed roles (npkm roles install)
```

---

## 📋 Module Quick Reference

| Module | Key Fields |
|---|---|
| `shell` / `command` | `cmd`, `cwd?` |
| `file` | `path`, `state` (directory/touch/link/absent), `mode?` |
| `copy` | `src`, `dest` |
| `move` | `src`, `dest` |
| `remove` | `path` |
| `debug` | `msg` |
| `template` | `src`, `dest`, `vars?` |
| `lineinfile` | `path`, `line`, `regexp?` |
| `replace` | `path`, `regexp`, `replace` |
| `get_url` | `url`, `dest` |
| `git` | `repo`, `dest` |
| `archive` / `unzip` | `src`, `dest` |
| `package` | `name`, `state`, `manager?` |
| `service` / `systemd` | `name`, `state`, `enabled?` |
| `user` | `name`, `state` |
| `cron` | `name`, `job`, `minute/hour/day/month/weekday?`, `state?` |
| `path` | `path` |
| `powershell` | `inline?`, `file?` |
| `fail` | `msg` |
| `set_fact` | `{key: value, ...}` — merges into runtime vars |
| `test` | `cmd`, `expect?`, `contains?` |
| `coni` | `script` — inline Coni expression |
| `include_tasks` | path, directory, or git URL |
| `block` | `block:`, `rescue:?`, `always:?` |
