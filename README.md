# NPKM (Nicolas's Playbook Kit Manager)

NPKM is a lightweight, declarative automation and provisioning tool (similar to Ansible or Chef), designed for zero-friction environment bootstrapping. It is distributed across two implementations providing exact feature parity.

## Implementations

- **[npkm-go](./npkm-go/)**: The original Go-based implementation built on `gopkg.in/yaml.v3` and `go-git`. Robust, strongly typed, and compiled easily into standalone binaries.
- **[npkm-coni](./npkm-coni/)**: A Drop-in replacement implementation written natively in the **Coni** programming language. Features a custom YAML-to-EDN parser and relies on shell-based native abstractions.

## Feature Parity Matrix

| Feature / Task | `npkm-go` | `npkm-coni` | Notes |
| :--- | :---: | :---: | :--- |
| **Core Architecture** | Go | Coni (Lisp-syntax) | |
| **Cross-OS Build** | ✅ (Mac, Win, Linux) | ✅ (Mac, Win, Linux) | Both compile entirely to `.exe` and `Mach-O` |
| **Remote HTTP Playbooks** | ✅ | ✅ | Can run playbooks directly via URL |
| **Git Repositories** | ✅ (`go-git`) | ✅ (`git clone`) | Scans cloned repo for playbook yaml/edn |
| **Directory Scanning** | ✅ | ✅ | Recursively lists available playbook files |
| **YAML Support** | ✅ (Strict) | ✅ (`yaml-to-edn`) | Natively transforms Ansible-style tasks |
| `file` | ✅ | ✅ | directory, touch, link, absent, modes |
| `lineinfile` | ✅ | ✅ | Regex matching & replacement in streams |
| `replace` | ✅ | ✅ | Replaces all instances of a regex pattern |
| `path` | ✅ | ✅ | Patches `.bashrc` / Powershell registry |
| `systemd` | ✅ | ✅ | start, stop, restart daemons |
| `copy`, `move`, `remove` | ✅ | ✅ | Standard IO primitives |
| `get_url` / `unzip` | ✅ | ✅ | Downloads and extracts remote assets |
| `shell`, `command`, `pwsh`| ✅ | ✅ | Shell integration along with Powershell |
| `debug`, `fail` | ✅ | ✅ | Playbook execution handling |
| `package` | ✅ | ✅ | Auto-detects brew, apt-get, yum, or choco |
| `service` | ✅ | ✅ | Generalizes systemctl, launchctl, and net start |
| `cron` | ✅ | ✅ | UNIX crontab -l / - insertion & absent state |
| `user` | ✅ | ✅ | Integrates useradd, sysadminctl, net user |
| `archive` | ✅ | ✅ | tar and zip abstraction across paths |
| `template` | ✅ | ✅ | Deploy templated files with mapped vars |

## Usage
Provide either a local YAML file, a directory, a remote HTTP/HTTPS link, or an SSH Git path:
```bash
# NPKM Go
cd npkm-go
./npkm playbook.sample.yml

# NPKM Coni
cd npkm-coni
./npkm-coni ssh://git@s5:2222/hellonico/my-playbook.git
```
