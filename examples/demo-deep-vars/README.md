# Deep Variables Example

This example demonstrates how NPKM resolves deeply nested variables across multiple variable scopes (global, group, and host).

## Structure
- `vars/main.yml`: Global variables.
- `inventories/dev/inventory.yml`: The inventory file defining the `servers` group and `server1` host.
- `inventories/dev/group_vars/servers.yml`: Group-specific variables.
- `inventories/dev/host_vars/server1.yml`: Host-specific variables.
- `playbook/deep.yml`: The playbook using these variables.

## Running the Example
```bash
coni main.coni -i examples/demo-deep-vars/inventories/dev/inventory.yml examples/demo-deep-vars/playbook/deep.yml
```
