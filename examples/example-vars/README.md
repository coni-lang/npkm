# NPKM Variables Example

This example demonstrates how NPKM resolves variables hierarchically using `group_vars` and `host_vars`.

## Structure

```text
example-vars/
├── inventory.yml             # Defines hosts and groups (webservers, dbservers)
├── group_vars/
│   ├── all.yml               # Applies to all hosts
│   ├── dbservers.yml         # Applies only to the dbservers group
│   └── webservers.yml        # Applies only to the webservers group
├── host_vars/
│   ├── db1.yml               # Applies only to db1
│   └── web1.yml              # Applies only to web1 (overrides webservers group_vars)
└── main.yml                  # Playbook
```

## Running the Example

Run the following command from this directory:

```bash
../npkm -i inventory.yml main.yml
```

## Expected Behavior

- **`all`**: `app_name`, `deploy_user`, `global_env` will be available to all hosts (`web1`, `web2`, `db1`).
- **`group_vars`**: 
  - `webservers` (`web1`, `web2`) get `http_port: 80` and `service_type: frontend`.
  - `dbservers` (`db1`) gets `db_port: 5432` and `service_type: backend`.
- **`host_vars`**: 
  - `web1` overrides `http_port` to `8080` and adds `custom_message`.
  - `db1` overrides `db_port` to `5433` and adds `custom_message`.
  - `web2` receives no `host_vars` and relies on `group_vars` entirely.
