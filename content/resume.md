---
title: "Résumé"
description: "Software engineer shipping security and observability for agentic AI. Six years of scalable systems for data-intensive and AI applications, much of it under GxP."
date: 2026-07-27
lastmod: 2026-10-03
---

[Download PDF](/mohammed-khalid-shaik-resume.pdf) · [moknshaik@gmail.com](mailto:moknshaik@gmail.com) · [LinkedIn](https://www.linkedin.com/in/pdwytr/) · [GitHub](https://github.com/pdwytr) · Dallas, TX

Software engineer with 6 years of experience building scalable systems for
data-intensive and AI applications. Shipped security and observability features
for an enterprise AI agent harness used in life sciences research. Previously
owned large-scale data infrastructure at Takeda and Crinetics, and supervised
their AWS disaster recovery under GxP audit. Open to relocation.

## Experience

### Software Engineer, Integrated Analytics Solutions — May 2023 to present

#### AI agent harness for life sciences research

- **Least-privilege agent sandboxing.** Built least-privilege sandboxing for
  coding agents with macOS Seatbelt on a life sciences AI agent platform for
  drug research. Each agent can reach only the files, tools, MCP servers and
  context it's allowed to access.
- **Multi-agent patch pipeline.** Designed a multi-agent system for shipping
  production fixes and security patches. Workers run isolated in Docker
  Sandboxes with managed context and observability, and no merge lands until a
  read-only verifier agent produces test proof. Cut agent errors by 50% and
  tripled patch releases.
- **Agent observability.** Extended Codex and OpenCode support to the harness's
  agent monitor, so it tracks their token spend and blocked sessions, too.

#### Crinetics Pharmaceuticals, clinical-data delivery

- **Regulated data delivery.** Rebuilt Crinetics' regulated clinical-data
  delivery (SFTP to SMB) as an agent-operable service that AI agents and
  administrators can configure, monitor, diagnose and recover. It handles
  thousands of files a day and delivers 100-file bursts in under 3 seconds with
  exactly-once delivery and immutable GxP audit trails.

#### Takeda, R&D data infrastructure

- **AWS disaster recovery.** Owned AWS disaster recovery for a single-region R&D
  SaaS platform (10,000+ users, 20 TB), delivering a GxP-validated 2-hour RTO /
  4-hour RPO. Authored the standby environment in Terraform and Lambda
  automation that auto-attached backup EBS volumes to recovery instances at
  failover.
- **Infrastructure-as-code governance.** Served as the approving reviewer on the
  data infrastructure Terraform repository, where every applied plan mutated
  AWS infrastructure. Gated 50+ PRs with validated testing, documented
  releases, and change-control sign-off, enabling junior engineers to ship
  safely in a GxP environment.
- **FDA submission automation.** Built and owned an enterprise
  document-assembly service (FastAPI) that compiles clinical study RTF outputs
  into one FDA submission-ready PDF, using a spreadsheet control file to set
  section order, a clickable table of contents, and bookmarks. Assembly time
  fell by 97%, from 15+ minutes to 30 seconds per file.
  → [Case study](/blogs/document-conversion-api/)
- **Secure multi-user architecture.** Engineered the document-assembly service
  for regulated multi-user production use: isolated per-session storage, file
  integrity checks, and async processing for up to 100 concurrent sessions and
  1 GB workloads. It ran under GxP change control with zero cross-session data
  leaks, supporting 10,000+ users.

### Software Engineer, GameStop — Nov 2022 to May 2023

- **Workflow automation.** Automated promotion approvals with Python, SQL views,
  and REST APIs, cutting monthly review effort by 80% (100 hours → 20 hours)
  and eliminating the manual gaps that had previously exposed the team to
  compliance penalties.
  → [Case study](/blogs/promotion-approval-workflow/)
- **RAG chatbot.** Built an internal retrieval-augmented generation (RAG)
  chatbot over company documents, using Haystack for indexing, embedding, and
  retrieval, and PostgreSQL with pgvector for vector search.

### Software Engineer, M&G — Jan 2021 to Aug 2022

- **Self-healing monitoring.** Developed a service for monitoring and
  self-healing across 1,500+ servers. Resolved half of recurring tickets
  automatically, saving 100+ team hours per month and protecting SLA timelines.
- **Identity and access management.** Managed IAM for 100+ users across 15
  teams, provisioning and maintaining their access to applications, servers,
  storage, and other IT infrastructure.

### Data Scientist Intern, INSOFE — June 2020 to Dec 2020

- **Demand forecasting.** Built product demand forecasting models in R on
  Walmart's public retail sales data. Compared ARIMA and ARIMAX time series
  models against Random Forest, CNN, and LSTM models, and tuned hyperparameters
  to improve forecast accuracy.

## Projects

- **[Pigeon](https://github.com/pdwytr/pigeon)** (open source) — unified
  Claude Code, Codex, and OpenCode session observer.
  → [Blog post](/blogs/pigeon-agent-console/)
- **[quickxorhash](https://github.com/pdwytr/quickxorhash)** (open source) — Go
  CLI computing Microsoft QuickXorHash for OneDrive hash parity in production
  migrations.
- **Shell.ai Hackathon** (competition) — optimized a 50-turbine offshore wind
  farm layout to maximize annual energy output.

## Skills

- **AI / LLM engineering** — sandboxed agent execution, agent orchestration (Pydantic AI), model routing (LiteLLM), semantic caching (LangCache, LangChainWiki), RAG (Haystack, LlamaIndex)
- **Languages** — Python, Rust, JavaScript/TypeScript, SQL
- **Frontend** — Tauri, React
- **Backend & APIs** — FastAPI, SQLAlchemy, REST, WebSockets, async & concurrent programming, microservices, distributed systems
- **Cloud & infra** — AWS, Terraform, Kubernetes, Docker
- **Data** — PostgreSQL, Redis, pandas, Pydantic, Airflow
- **Testing & quality** — pytest, Locust, integration & e2e testing, AST static analysis, CI/CD
- **Compliance** — GxP, change control, audit readiness

## Education

- **M.S. Data Science** — The University of Texas at Arlington, Dec 2023
- **B.Tech.** — Jawaharlal Nehru Technological University, Hyderabad, Aug 2020
