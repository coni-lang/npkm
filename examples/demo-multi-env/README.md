# NPKM Multi-Environment Cluster Demo

> One playbook. Two environments. All nodes in parallel.

## Concept

The key insight: **the playbook never changes**. The environment is 100% defined by the inventory file. DEV1 and DEV2 are the same infrastructure — only the variables differ.

```
provision.edn          ← IDENTICAL for DEV1 and DEV2
inventory/dev1.edn     ← DEV1 hosts + region/AZ vars
inventory/dev2.edn     ← DEV2 hosts + region/AZ vars
group_vars/all.edn     ← shared across all envs
group_vars/dev1.edn    ← DEV1 overrides (db, redis, s3, log level...)
group_vars/dev2.edn    ← DEV2 overrides
roles/base/            ← OS baseline role
roles/app/             ← application deploy role
```

## Run

```bash
# Provision DEV1 cluster (3 nodes in parallel)
npkm -i inventory/dev1.edn provision.edn

# Provision DEV2 cluster (swap inventory — that's it)
npkm -i inventory/dev2.edn provision.edn

# Dry-run first to see what would happen
npkm --dry-run -i inventory/dev1.edn provision.edn

# Step through interactively
npkm --step -i inventory/dev1.edn provision.edn

# Generate an audit report
npkm --report -i inventory/dev1.edn provision.edn

# Watch for changes during active development
npkm watch -i inventory/dev1.edn provision.edn
```

## Variable Resolution Order

```
group_vars/all.edn        (lowest priority — shared defaults)
    ↓
inventory group :vars     (env-level: region, AZ, env name)
    ↓
group_vars/dev1.edn       (env-specific: db, redis, s3, log level)
    ↓
inventory host :vars      (host-specific: node_index, ansible_host)
    ↓
include_tasks :vars       (role-call overrides — highest priority)
```

## What changes between DEV1 and DEV2

| Variable      | DEV1                    | DEV2                    |
|---------------|-------------------------|-------------------------|
| `env`         | `dev1`                  | `dev2`                  |
| `aws_region`  | `us-east-1`             | `us-west-2`             |
| `instance_az` | `us-east-1a`            | `us-west-2b`            |
| `db_host`     | `db.dev1.internal`      | `db.dev2.internal`      |
| `db_name`     | `myapp_dev1`            | `myapp_dev2`            |
| `redis_host`  | `redis.dev1.internal`   | `redis.dev2.internal`   |
| `log_level`   | `DEBUG`                 | `INFO`                  |
| `s3_bucket`   | `myapp-dev1-assets`     | `myapp-dev2-assets`     |
| `replicas`    | `1`                     | `2`                     |

## Scaling to 10 EC2 instances

Add nodes to the inventory — the playbook and roles need zero changes:

```edn
; inventory/dev1.edn  — 10 nodes
{:dev1
 {:vars {:env "dev1" :aws_region "us-east-1"}
  :hosts
  {:dev1-node-1  {:ansible_host "10.0.1.11" :node_index 1}
   :dev1-node-2  {:ansible_host "10.0.1.12" :node_index 2}
   ; ... up to node-10
   :dev1-node-10 {:ansible_host "10.0.1.20" :node_index 10}}}}
```

```edn
; provision.edn — only forks changes (no logic change)
{:name "Cluster Baseline"
 :hosts "dev1"
 :forks 10   ← all 10 nodes provisioned simultaneously
 ...}
```

## Structure

```
demo-multi-env/
  provision.edn           ← single entry point for all envs
  inventory/
    dev1.edn              ← DEV1: 3 nodes, us-east-1
    dev2.edn              ← DEV2: 3 nodes, us-west-2
  group_vars/
    all.edn               ← shared: app_name, app_version, ports
    dev1.edn              ← DEV1: db, redis, s3, log_level
    dev2.edn              ← DEV2: db, redis, s3, log_level
  roles/
    base/
      tasks/main.edn      ← OS baseline: Java, users, directories
      defaults/main.edn
    app/
      tasks/main.edn      ← app config + systemd unit + smoke test
      defaults/main.edn
```
