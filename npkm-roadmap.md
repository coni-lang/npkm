# NPKM Feature Audit & Roadmap

## 1. Feature Audit

### ✅ What NPKM Has (Solid Foundation)
| Feature | NPKM | Ansible |
|---|---|---|
| Shell/Command execution | ✅ `shell`, `command`, `powershell` | ✅ |
| File management | ✅ `file`, `copy`, `move`, `remove`, `lineinfile`, `replace` | ✅ |
| Templating (`{{ var }}`) | ✅ | ✅ |
| Inventory (YAML, INI, inline) | ✅ | ✅ |
| SSH remote execution | ✅ | ✅ |
| Conditional execution (`when`) | ✅ | ✅ |
| Loops (`loop`, `with_items`, `items`) | ✅ | ✅ |
| Variable `register` | ✅ | ✅ |
| `include_tasks` (local, dir, git URL) | ✅ | ✅ |
| Package management | ✅ `package` | ✅ |
| Service management | ✅ `service`, `systemd` | ✅ |
| User management | ✅ `user` | ✅ |
| Cron management | ✅ `cron` | ✅ |
| HTTP file download | ✅ `get_url` | ✅ |
| Git clone/pull | ✅ `git` | ✅ |
| Archive/zip | ✅ `archive`, `unzip` | ✅ |
| `--doc` Mermaid flow generation | ✅ 🔥 **UNIQUE** | ❌ |
| Label/name filtering (`--labels`, `--names`) | ✅ | ❌ tags only |
| EDN format support | ✅ 🔥 **UNIQUE** | ❌ |
| Native binary (no Python/runtime) | ✅ 🔥 **UNIQUE** | ❌ |
| Persistent run logs in `~/.npkm/` | ✅ | ❌ |
| `become` (sudo escalation) | ✅ | ✅ |
| Cross-platform (macOS/Linux/Windows) | ✅ | Partial |

---

### ❌ What Ansible Has That You Don't
These are the real gaps, in priority order:

| Gap | Impact | Effort |
|---|---|---|
| **Parallel host execution** (`forks`) | ✅ Done | Medium |
| **Handlers + `notify`** | ✅ Done | Low |
| **`block` / `rescue` / `always`** | ✅ Done | Medium |
| **`retry` / `until`** | ✅ Done | Low |
| **Vault (encrypted secrets)** | 🟡 Medium — secure credential storage | Medium |
| **`check_mode` (dry-run)** | ✅ Done | Low |
| **Idempotent state reporting** | ✅ Done — via `changed_when` | Low |
| **Dynamic inventory** | 🟠 Nice to have | Medium |

---

## 2. Best Plan of Action

We can structure the upcoming work into sprints to rapidly close the core gaps and emphasize NPKM's unique strengths over Ansible.

| Phase / Sprint | Goal | Sub-Tasks |
|---|---|---|
| ✅ **Sprint 1: Core Reliability** | Close basic operational gaps | <ul><li>[x] Implement `--dry-run` / `--check` mode</li><li>[x] Implement `retry: 3` and `delay: 5` (until success)</li><li>[x] Add support for `ok`, `changed`, and `skipped` states per task</li><li>[x] Windows compatibility in demo playbooks</li></ul> |
| ✅ **Sprint 2: Flow Control** | Advanced playbook structure | <ul><li>[x] Implement `handlers` and `notify`</li><li>[x] Implement `block`, `rescue`, `always` for error boundaries</li></ul> |
| ✅ **Sprint 3: The Multi-Node Killer Feature** | True parallel execution | <ul><li>[x] Refactor SSH loop to use goroutines (channels) for concurrent host execution</li><li>[x] Add `forks: 5` playbook parameter</li><li>[x] Implement `parallel: true` task groups</li></ul> |
| **Sprint 4: Ecosystem & Uniqueness** | Lean into Coni/EDN | <ul><li>[x] Create native `coni:` task module (inline scripts inside playbooks)</li><li>[ ] Build `npkm-galaxy` style hub (git repo convention)</li><li>[ ] Add `--diff` mode for showing file changes</li></ul> |
