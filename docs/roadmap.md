# **Roadmap: Velo**

## Alpha Requirements

> *Goal: It boots. It deploys. It doesn't burst into flames.*
> **Target User**: Developer building the platform
> **Focus**: Core deployment logic, basic UI/API, local use

* [ ] Docker Swarm / Nomad integration 
* [ ] Single-node install script (bash/ansible)
* [ ] Web UI skeleton or CLI-only mode
* [ ] Manual deployment definition (YAML/Compose)
* [ ] Basic service CRUD (Create/Update/Delete)
* [ ] Container status dashboard (running, error, restarting)
* [ ] Persisting configs to disk or DB
* [ ] Simple user login with static credentials (no OAuth yet)

---

## Beta Requirements

> *Goal: Tinkerable and safe for personal use*
> **Target User**: Homelabbers, dev teams, early adopters
> **Focus**: Security, installation UX, better management

* [ ] Automated install process

  * [ ] Interactive CLI installer or single-command bootstrapping
  * [ ] Works on major distros (Ubuntu, Debian, Fedora)

* [ ] Security Improvements

  * [ ] Automated HTTPS certificate generation (Let's Encrypt)
  * [ ] Secure secrets storage (encrypted at rest)
  * [ ] User authentication (basic OAuth + roles)
  * [ ] SSH key or token-based node authentication

* [ ] Improved Deployment UX

  * [ ] Web-based YAML editor with validation
  * [ ] Template library (e.g. NGINX, PostgreSQL, Redis)
  * [ ] Git-based deployment sync (GitOps-lite)

* [ ] Cluster Management

  * [ ] Node health checks (CPU, RAM, Disk, status)
  * [ ] Add/remove nodes via UI or CLI
  * [ ] Labels/tags for grouping nodes

* [ ] Basic Clean-up / Rollback

  * [ ] Rollback to previous service version
  * [ ] Auto-cleanup of failed/stuck deployments

---

## Stable Requirements

> *Goal: Production-capable, low-maintenance, self-healing system*
> **Target User**: Small teams or self-hosters running production or pre-prod environments
> **Focus**: Reliability, automation, visibility

* [ ] Self-Healing

  * [ ] Auto-restart failed services
  * [ ] Node drain/rebalance on failure
  * [ ] Deployment retry strategies

* [ ] Automated Observability

  * [ ] Integrated logging (per-service viewer)
  * [ ] Metrics collection (CPU/RAM/Disk/Net per container)
  * [ ] Service health dashboard
  * [ ] Configurable alerts (Slack/email/webhook)

* [ ] Deployment Strategies

  * [ ] Canary releases
  * [ ] Rolling updates with health checks
  * [ ] Manual approval step (optional)

* [ ] Access Control & Audit

  * [ ] Role-based access control (view, deploy, admin)
  * [ ] Audit log of user actions (who deployed what, when)

* [ ] Plugin/Hook System

  * [ ] Pre/post-deploy hooks
  * [ ] Slack/webhook notification support
  * [ ] CLI extensibility with custom scripts

* [ ] Backup & Restore

  * [ ] Backup of configs, volumes, metadata
  * [ ] One-click restore for disaster recovery

