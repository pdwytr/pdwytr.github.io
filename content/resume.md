---
title: "Résumé"
description: "Software engineer building security, observability, and CI/CD for agentic AI. Six years of APIs and data-intensive systems, much of it under GxP."
date: 2026-07-27
lastmod: 2026-09-30
---

[Download PDF](/mohammed-khalid-shaik-resume.pdf) · [moknshaik@gmail.com](mailto:moknshaik@gmail.com) · [LinkedIn](https://www.linkedin.com/in/pdwytr/) · [GitHub](https://github.com/pdwytr) · Dallas, TX

Software engineer with 6 years of experience building APIs and data-intensive
applications. Currently engineering security, observability, and CI/CD
pipelines for a commercial agentic AI platform used in life sciences research.
Previously owned large-scale data infrastructure for companies including Takeda
and Crinetics, along with leading their AWS disaster recovery under formal
change control and GxP audit.

## Experience

### Software Engineer, Integrated Analytics Solutions — May 2023 to present

#### AI agent platform for life sciences research

- Built least-privilege sandboxing for coding agents on a life sciences AI
  agent platform used by 30+ drug researchers. It uses macOS Seatbelt, chosen
  after ruling out Apple's App Sandbox, so each agent can only reach the files
  its task allows.
- Built a 10-trial evaluation that blocks new agent personalities whose
  behavior breaks their guardrails, and added Codex and Hermes to the
  platform's agent monitoring, covering token spend and blocked sessions.
- Designed a multi-agent system for shipping production fixes and security
  patches: controlling what context each agent receives, isolating workers in
  Docker Sandboxes, tracking each agent's token spend, and gating every merge on
  a read-only verifier agent's test proof. Cut reopened tickets by about 50%.

#### Crinetics Pharmaceuticals, clinical-data delivery

- Rebuilt Crinetics' regulated clinical-data delivery (SFTP to SMB) as an
  agent-operable service that AI agents and administrators can configure,
  monitor, diagnose and recover. It handles thousands of files a day and
  delivers ~100-file bursts in under 3 seconds with exactly-once delivery and
  immutable GxP audit trails.

#### Takeda, R&D data infrastructure

- Owned AWS disaster recovery for a single-region R&D SaaS platform (1,000+
  users, ~20 TB), delivering a GxP-validated 2-hour RTO / 4-hour RPO. Authored
  the standby environment in Terraform and Lambda automation that auto-attached
  backup EBS volumes to recovery instances at failover.
- Served as the approving reviewer on the data infrastructure Terraform
  repository, where every applied plan mutated AWS infrastructure. Gated 50+ PRs
  with validated testing, documented releases, and change-control sign-off,
  enabling junior engineers to ship safely in a GxP environment.
- Built a FastAPI and SAS automation service that cut manual FDA clinical
  report RTF-to-PDF assembly from 15+ minutes to under 30 seconds per (~10 MB)
  file, supporting RTFs up to 1 GB for 10,000+ statistical users.
  → [Case study](/blogs/document-conversion-api/)

### Software Engineer, GameStop — Nov 2022 to May 2023

- Automated promotion approvals with Python, SQL, and REST APIs, cutting
  monthly review effort by 80% (100 hours → 20 hours) and eliminating the
  manual gaps that had previously exposed the team to compliance penalties.
  → [Case study](/blogs/promotion-approval-workflow/)

### Software Engineer, M&G — Jan 2021 to Aug 2022

- Developed a service for monitoring and self-healing across 1,500+ servers,
  resolving about half of recurring tickets automatically and saving 100+ team
  hours per month while protecting SLA timelines.

## Skills

- **AI / LLM engineering** — agent orchestration (LangGraph, Pydantic AI), model routing (LiteLLM), semantic caching (LangCache), retrieval (RAG, pgvector), sandboxed agent execution (Docker Sandboxes, macOS Seatbelt, Windows AppContainer)
- **Languages** — Python, Rust, JavaScript/TypeScript, SQL
- **Frameworks** — Tauri, FastAPI, React
- **Backend & APIs** — REST, WebSockets, async & concurrent programming, microservices, distributed systems
- **Cloud & infra** — AWS, Terraform, Kubernetes, Docker
- **Data** — PostgreSQL, Redis, pandas, Pydantic, SQLAlchemy, Airflow
- **Testing & quality** — pytest, Locust load testing, integration & e2e testing, AST static analysis, CI/CD
- **Compliance** — GxP, change control, audit readiness

## Education

- **M.S. Data Science** — The University of Texas at Arlington, Dec 2023
