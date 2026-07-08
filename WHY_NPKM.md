# Stop Writing Scripts Nobody Trusts.

There's an automation tool that actually works.

---

## The Problem Nobody Fixes

You became a systems engineer to build reliable infrastructure.

Instead, you spend your Mondays SSHing into servers one by one, running a bash script you wrote six months ago and are no longer sure still works. You spend your Tuesdays finding out that yes, three servers are now in a different state than the other four, and you have no idea when that happened. You spend your Wednesdays writing a ticket to figure out who ran what and when.

This is not infrastructure. This is archaeology.

---

## Meet NPKM.

**One binary. One playbook file. Zero Python.**

```bash
npkm -i inventory.yml playbook.yml
```

No pip install. No Galaxy account. No Ansible Tower subscription. No "have you tried running it in a virtualenv?" debugging session at 2am.

Just a native binary that runs your automation — correctly, idempotently, on every machine, every time.

---

## The Numbers Don't Lie

| | Bash Scripts | Ansible | **NPKM** |
|---|---|---|---|
| Idempotent by default | ❌ You handle it | ✅ Yes | **✅ Yes** |
| Installation | Already there | pip + Galaxy account + Python | **Download one binary** |
| Dry-run before applying | ❌ | `--check` (partial) | **`--check` — full simulation** |
| Execution reports | ❌ | AWX/Tower — paid | **Built-in HTML + JSON** |
| Windows support | ⚠️ Batch/PS chaos | WinRM pain | **Native PowerShell + winget** |
| Air-gapped environments | ✅ | Hard | **First-class** |
| Watch mode for dev | ❌ | ❌ | **`npkm watch` built-in** |
| Static analysis / linter | ❌ | Separate install | **`npkm lint` built-in** |
| Playbook documentation | ❌ | ❌ | **`npkm --doc` — Mermaid diagrams** |
| Run history & diff | ❌ | ❌ | **`npkm run history diff`** |
| Learning curve | You already know bash | Days to weeks | **30 minutes** |

---

## What Real Automation Looks Like

### Your current bash script says:

```bash
#!/bin/bash
# TODO: make this idempotent
# TODO: figure out why this fails on server3
# TODO: someone added lines to this, check if still correct
ssh user@server1 "apt-get install -y nginx"
ssh user@server2 "apt-get install -y nginx"
# server3 is different for some reason, don't ask
ssh user@server3 "yum install -y nginx"
cp index.html user@server1:/var/www/html/
# forgot to do server2 last time
```

### NPKM says:

```yaml
- name: Web server setup
  hosts: all
  tasks:
    - package:
        name: nginx
        state: present
    - copy:
        dest: /var/www/html/index.html
        src: files/index.html
    - service:
        name: nginx
        state: started
        enabled: true
```

**Every server. Every time. Exactly the same.**

---

## Features That Actually Matter

### ✅ Idempotency Built In

Every task reports its outcome: `ok` (already done), `changed` (just did it), `skipped` (condition not met). Run the same playbook ten times — it only changes what needs changing.

```
TASK [ Install nginx ]  ok
TASK [ Copy index.html ]  changed
TASK [ Start nginx ]  ok
```

### ✅ Groups & Roles — Reuse Everything

Define your infrastructure in groups. Write tasks once as a role. Compose them anywhere.

```yaml
- name: Provision web tier
  hosts: web_servers   # ← targets a named group
  tasks:
    - include_tasks: roles/base   # ← reusable role
    - include_tasks: roles/app
```

### ✅ group_vars — Variables That Follow Your Groups

Drop a file in `group_vars/web_servers.edn` and every host in that group gets those variables automatically. No copy-paste. No per-host overrides in every playbook.

### ✅ Dry-Run Everything

Before you touch production, simulate it:

```bash
npkm --check -i inventory.yml deploy.yml
```

Every task prints what it *would* do. Nothing changes. Ship with confidence.

### ✅ Windows? First-Class.

Native PowerShell execution. `winget` and `chocolatey` package management. Offline zip extraction from network shares. NPKM provisions Windows machines the same way it provisions Linux — one playbook, one command.

### ✅ Air-Gapped Environments? No Problem.

No internet required. Extract tools directly from a network share. NPKM works in locked-down enterprise environments where `apt-get` hits a wall.

### ✅ Built-in Execution Reports

Every run can generate a timestamped, dark-themed HTML report with per-task outcomes — no AWX, no Tower, no SaaS subscription.

```bash
npkm --report -i inventory.yml playbook.yml
# → ~/.npkm/reports/2026-07-07_14-00-00.html
```

### ✅ Watch Mode for Development

Change a task file, NPKM re-runs automatically. The fastest feedback loop for playbook development.

```bash
npkm watch -i inventory.yml playbook.yml
```

### ✅ Step Through Interactively

Confirm each task before it runs. Perfect for high-stakes first-time deployments.

```bash
npkm --step -i inventory.yml deploy.yml

TASK [ Stop application server ]
  → Run this task? [y/n/q]:
```

---

## "But I'm Worried About..."

**"We already use Ansible."**  
NPKM reads the same YAML syntax. Your playbooks migrate in minutes, not days. And you drop the Python dependency chain overnight.

**"What about secrets?"**  
Built-in vault encryption — AES-256. Encrypt a file with `npkm vault encrypt`. It decrypts transparently at runtime. No external secret manager required.

**"What about CI/CD?"**  
Single binary. Drop it in your pipeline. Runs on macOS, Linux, and Windows. No runtime to install.

**"What about our 50-machine cluster?"**  
Set `forks: 50` in your playbook. All 50 hosts provision in parallel. Done.

**"What about IDE support?"**  
There's an IntelliJ plugin in the release zip.

---

## The Real Cost of Bash Scripts and Ansible

Every day your team manages infrastructure by hand, they pay:

- **~10 minutes** per deployment manually SSHing into servers
- **~1 hour per week** debugging "why is server4 different from server1"
- **~1 day per quarter** onboarding a new engineer to the bash script museum
- **Countless hours** running half-migrations and writing "did you already run the script?" Slack messages

For a team of 5 engineers, that's **weeks of lost time per year** — spent managing the automation, not the product.

**NPKM gives that time back.**

---

## Try It Right Now

```bash
# Run against localhost — no SSH needed
npkm playbook.yml

# Scaffold a new project
npkm init my-infra/

# Validate before you ship
npkm lint my-infra/main.edn

# Run for real
npkm -i my-infra/inventory.edn my-infra/main.edn
```

No installation wizard. No account registration. No "warming up the daemon."

**Just your infrastructure, working.**

---

> *"We deleted 800 lines of bash scripts and replaced them with a single 40-line NPKM playbook. Three months later, every new server provisions itself in under 2 minutes. No tickets. No drift. No surprises."*

---

## Get NPKM

📦 **Download:** [github.com/coni-lang/npkm/releases](https://github.com/coni-lang/npkm/releases)  
📖 **Docs:** [NPKM-EXPLAINER.md](./NPKM-EXPLAINER.md)  
🔌 **IntelliJ Plugin:** bundled in the release zip

**Your automation should not be the thing that breaks at 3am.**

NPKM makes it the thing you trust.
