# NPKM Variables Example

This example demonstrates how NPKM automatically loads variables from multiple locations.

## Structure
- `vars/main.yml`: Automatically loaded as global variables by `playbook.yml`.
- `group_vars/all.yml`: Automatically loaded and applied to all hosts in the inventory.
- `host_vars/web2.yml`: Automatically loaded for `web2`, overriding the `group_vars`.

## Run the example

```bash
npkm -i inventory.edn playbook/playbook.yml
```

You will see `web1` uses port 80 (from `group_vars`), and `web2` uses port 8080 (from `host_vars`). Both will see the `app_version` from `vars/main.yml`.
