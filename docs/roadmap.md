# **Roadmap: Velo**

## Alpha Requirements

> _Goal: It boots. It deploys. It doesn't burst into flames._ > **Target User**: Developer building the platform
> **Focus**: Core deployment logic, basic UI/API, local use

- [x] Docker Swarm integration
- [x] Single-node install script (bash/ansible)
- [x] CLI-only mode 
- [x] Manual deployment definition (TOML/YAML)
- [x] Basic service CRUD (Create/Update/Delete)
- [x] Container status dashboard (running, error, restarting)
- [x] Persisting configs to disk or DB
- [x] Simple user login with static credentials (no OAuth yet)

---

## Beta Requirements 🚧 IN PROGRESS

> _Goal: Tinkerable and safe for personal use_ > **Target User**: Homelabbers, dev teams, early adopters
> **Focus**: Security, installation UX, better management

- [x] Automated install process

  - [x] Interactive CLI installer or single-command bootstrapping
  - [x] Works on major distros (Ubuntu, Debian, OpenSUSE, Fedora, Arch, Alpine)

- [x] Security Improvements

  - [ ] Automated HTTPS certificate generation (Let's Encrypt)
  - [ ] Secure secrets storage (encrypted at rest)
  - [x] User authentication (basic static credentials)
  - [ ] SSH key or token-based node authentication

- [x] Improved Deployment UX

  - [x] Web-based deployment interface with forms
  - [ ] Template library (e.g. NGINX, PostgreSQL, Redis)
  - [ ] Git-based deployment sync (GitOps-lite)

- [x] Cluster Management

  - [x] Node health checks (CPU, RAM, Disk, status)
  - [x] Add/remove nodes via CLI 
  - [x] Labels/tags for grouping nodes

- [x] Basic Clean-up / Rollback

  - [x] Rollback to previous service version
  - [ ] Auto-cleanup of failed/stuck deployments

---

## Stable Requirements

> _Goal: Production-capable, low-maintenance, self-healing system_ > **Target User**: Small teams or self-hosters running production or pre-prod environments
> **Focus**: Reliability, automation, visibility

- [ ] Self-Healing

  - [ ] Auto-restart failed services
  - [ ] Node drain/rebalance on failure
  - [ ] Deployment retry strategies

- [ ] Automated Observability

  - [ ] Integrated logging (per-service viewer)
  - [ ] Metrics collection (CPU/RAM/Disk/Net per container)
  - [ ] Service health dashboard
  - [ ] Configurable alerts (Slack/email/webhook)

- [ ] Deployment Strategies

  - [ ] Canary releases
  - [ ] Rolling updates with health checks
  - [ ] Manual approval step (optional)

- [ ] Access Control & Audit

  - [ ] Role-based access control (view, deploy, admin)
  - [ ] Audit log of user actions (who deployed what, when)

- [ ] Plugin/Hook System

  - [ ] Pre/post-deploy hooks
  - [ ] Slack/webhook notification support
  - [ ] CLI extensibility with custom scripts

- [ ] Backup & Restore

  - [ ] Backup of configs, volumes, metadata
  - [ ] One-click restore for disaster recovery
