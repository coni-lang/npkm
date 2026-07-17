# NPKM vs Ansible: Syntax & Feature Parity

NPKM aims to provide a zero-dependency, ultra-fast alternative to Ansible while maintaining extremely high syntax parity. Playbooks written for Ansible can often be executed by NPKM with zero modifications.

## 1. Core Playbook Structure

| Feature | Ansible Syntax | NPKM Support | Notes |
|---------|---------------|--------------|-------|
| **Playbook Structure** | List of plays (`- name: ...`) | ✅ Supported | NPKM also supports single-map playbooks natively. |
| **Hosts Definition** | `hosts: webservers` | ✅ Supported | Groups and specific hosts are mapped via `inventory.yml`. |
| **Tasks Definition** | `tasks:` list | ✅ Supported | Deeply nested and inline maps supported. |
| **Handlers** | `handlers:` and `notify:` | ✅ Supported | Same event-driven task resolution. |
| **Variables Definition**| `vars:` block | ✅ Supported | Scoped to the play context automatically. |
| **Includes** | `include_tasks:`, `import_tasks:` | ✅ Supported | Seamlessly includes modular tasks. |
| **Roles** | `roles:` list | ✅ Supported | Full directory traversal (`tasks/main.yml`, `vars/main.yml`). |

## 2. Variables & Jinja2 Templating Engine

NPKM features a standalone `j2` module built natively in Coni. It eliminates the need for a Python dependency while maintaining advanced Jinja2 macro processing.

### Supported Jinja2 Filters & Features
| Feature | Ansible Syntax | NPKM Support | Notes |
|---------|---------------|--------------|-------|
| **Variable Injection** | `{{ my_var }}` | ✅ Supported | Standard string interpolation. |
| **Nested Variables** | `{{ user.name }}` | ✅ Supported | Full map traversal and property dot-notation. |
| **String Filters** | `{{ var &#124; upper }}`, `lower` | ✅ Supported | Converts strings dynamically. |
| **Default Fallback** | `{{ missing &#124; default('X') }}` | ✅ Supported | Supports fallback for undefined values. |
| **Ternary Operator** | `{{ bool &#124; ternary('T', 'F') }}` | ✅ Supported | Boolean evaluation directly in template. |
| **Data Serialization** | `{{ obj &#124; to_json }}` | ✅ Supported | Additionally supports `to_edn`. |
| **List Joining** | `{{ list &#124; join(',') }}` | ✅ Supported | Formats arrays as delimited strings. |
| **Native Execution** | *N/A (Python Eval)* | 🚀 **NPKM Exclusive** | Execute raw Coni functions: `{{ var &#124; (fn [x] ...) }}` |
| **Magic Variables** | `inventory_hostname` | ✅ Supported | Auto-injected (`npkm_os_family`, `groups`, etc.) |

### Jinja2 Examples

Here are some detailed examples of how you can leverage Jinja2 templating natively in NPKM without any Python dependencies:

**1. Basic Variable Injection & Defaulting:**
```yaml
- name: Greet the user
  debug:
    msg: "Hello {{ user.name | default('Admin') }}, welcome to NPKM!"
```

**2. Serialization and Conditionals:**
```yaml
- name: Show API config
  debug:
    msg: "Config: {{ api_config | to_json }} - Enabled: {{ is_active | ternary('YES', 'NO') }}"
```

**3. The Power of Native Coni Filters:**
Because NPKM runs on Coni, you aren't restricted by standard Jinja filters. If a filter isn't recognized, NPKM evaluates it as a raw Coni anonymous function!
```yaml
- name: Uppercase an entire list dynamically
  debug:
    # (fn [x] ...) executes arbitrary Coni language logic inline!
    msg: "Roles: {{ user_roles | (fn [x] (str/join \", \" (map str/upper x))) }}"
```

## 3. Inventory Management

| Feature | Ansible Syntax | NPKM Support | Notes |
|---------|---------------|--------------|-------|
| **YAML Inventory** | `all: hosts: ...` | ✅ Supported | NPKM consumes standard Ansible YAML inventories. |
| **INI Inventory** | `[webservers]` | ✅ Supported | Native INI parsing support. |
| **Host Variables** | Defined under `vars:` | ✅ Supported | Evaluated and merged per-host. |
| **Group Variables** | `group_vars/` directory | ✅ Supported | |

## 4. Execution & Flow Control

| Feature | Ansible Syntax | NPKM Support | Notes |
|---------|---------------|--------------|-------|
| **Privilege Escalation**| `become: yes` | ✅ Supported | Evaluates sudo natively across platforms. |
| **Looping** | `loop:` / `with_items:` | ✅ Supported | Resolves lists and loops the specific task. |
| **Conditionals** | `when: var == 'test'` | ✅ Supported | Full boolean conditional skipping. |
| **Delegation** | `delegate_to:` | 🚧 In Progress | Planned for next major milestone. |
| **Parallel Execution** | `strategy: free` | 🚀 **Enhanced** | NPKM supports `parallel: true` groups via Go channels. |

## 5. Modules

NPKM implements a robust list of core Ansible modules directly in native Coni. These run instantly with Go concurrency, drastically reducing overhead compared to Ansible's Python bootstrapping.

### Full List of Supported Modules
- **System**: `command`, `shell`, `powershell`, `win_shell`, `coni` (run native Coni scripts!)
- **Files & Directories**: `file`, `copy`, `template`, `remove`, `move`, `stat`, `path`
- **File Contents**: `lineinfile`, `replace`
- **Network & Source Control**: `get_url`, `git`
- **Packaging & Archives**: `package`, `unzip`, `archive`
- **System Configuration**: `systemd`, `service`, `cron`, `user`
- **Debugging & Control**: `debug`, `fail`

### Module Syntax Example
```yaml
- name: Deploy web application
  hosts: webservers
  vars:
    app_version: "1.0.4"
  tasks:
    - name: Clone repository
      git:
        repo: "https://github.com/my-org/my-app.git"
        dest: "/var/www/app"
        version: "{{ app_version }}"

    - name: Configure systemd service
      template:
        src: "app.service.j2"
        dest: "/etc/systemd/system/app.service"
      notify: restart_app

  handlers:
    - name: restart_app
      systemd:
        name: "app"
        state: "restarted"
```

### Recommended Future Modules (Roadmap)
To achieve even higher parity with enterprise Ansible deployments, we recommend adding support for:
1. **`apt` / `yum` / `brew` explicit aliases** (Currently handled dynamically by the `package` module, but explicit modules improve backwards compatibility).
2. **`wait_for` / `wait_for_connection`** (Crucial for deployments involving reboots or waiting for application ports to open).
3. **`docker_container` / `docker_image`** (Highly requested for containerized deployments).
4. **`uri` / `htpasswd`** (Expanding upon `get_url` to handle complex API interactions or basic auth setups).
5. **`set_fact`** (To dynamically store computed Coni values during the playbook run).
