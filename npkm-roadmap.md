# NPKM Feature Audit & Capabilities

## ✅ NPKM Complete Features
| Feature | NPKM | Ansible |
|---|---|---|
| Shell/Command execution | ✅ `shell`, `command`, `powershell` | ✅ |
| File management | ✅ `file`, `copy`, `move`, `remove`, `lineinfile`, `replace` | ✅ |
| Templating (`{{ var }}`) | ✅ | ✅ |
| Static Inventory (YAML, EDN, INI, inline) | ✅ | ✅ |
| Dynamic Inventory (Executable scripts) | ✅ | ✅ |
| SSH remote execution | ✅ | ✅ |
| Parallel host execution (`forks`) | ✅ | ✅ |
| Conditional execution (`when`) | ✅ | ✅ |
| Loops (`loop`, `with_items`, `items`) | ✅ | ✅ |
| Variable `register` | ✅ | ✅ |
| Error handling (`block`, `rescue`, `always`) | ✅ | ✅ |
| Event triggers (`handlers`, `notify`) | ✅ | ✅ |
| Task retry loops (`retry`, `until`, `delay`) | ✅ | ✅ |
| `include_tasks` (local, dir, git URL) | ✅ | ✅ |
| Role Package Manager (`npkm roles install`) | ✅ | ✅ |
| Vault (encrypted secrets & runtime decryption) | ✅ | ✅ |
| Package management | ✅ `package` | ✅ |
| Service management | ✅ `service`, `systemd` | ✅ |
| User management | ✅ `user` | ✅ |
| Cron management | ✅ `cron` | ✅ |
| HTTP file download | ✅ `get_url` | ✅ |
| Git clone/pull | ✅ `git` | ✅ |
| Archive/zip | ✅ `archive`, `unzip` | ✅ |
| Dry-run mode (`--dry-run`, `--check`) | ✅ | ✅ |
| File changes mode (`--diff`) | ✅ | ✅ |
| Idempotent state reporting (`ok`, `changed`) | ✅ | ✅ |
| `become` (sudo escalation) | ✅ | ✅ |
| `--doc` Mermaid flow generation | ✅ 🔥 **UNIQUE** | ❌ |
| Label/name filtering (`--labels`, `--names`) | ✅ | ❌ tags only |
| EDN format support (Tasks, Vars, Inventory) | ✅ 🔥 **UNIQUE** | ❌ |
| Native `coni:` task module (inline scripts) | ✅ 🔥 **UNIQUE** | ❌ |
| Native binary (no Python/runtime) | ✅ 🔥 **UNIQUE** | ❌ |
| Persistent run logs in `~/.npkm/` | ✅ | ❌ |
| Cross-platform (macOS/Linux/Windows) | ✅ | Partial |
