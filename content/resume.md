---
title: "Résumé"
description: "Software engineer with 5+ years building data and AI systems: agent sandboxing, multi-agent orchestration and observability, plus GxP data infrastructure for Takeda, Crinetics and Jazz."
date: 2026-07-27
lastmod: 2026-10-11
---

[Download PDF](/mohammed-khalid-shaik-resume.pdf) · [moknshaik@gmail.com](mailto:moknshaik@gmail.com) · [LinkedIn](https://www.linkedin.com/in/pdwytr/) · [GitHub](https://github.com/pdwytr)

Software engineer with 5+ years of experience building data and AI systems.
Shipped sandboxing, multi-agent orchestration, and observability for an AI
coding-agent platform for drug research. Also owned GxP data infrastructure
development and deployments for clinical study data at Takeda, Crinetics, and
Jazz, including AWS disaster recovery.

## Experience

### Software Engineer, Integrated Analytic Solutions — May 2023 to present

Clients: Takeda, Crinetics, SOBI, Jazz, Neurocrine · Remote

#### AI coding-agent platform for drug research

- **Agent sandboxing.** Restricted tool calls by coding agents to explicitly
  allowed resources using least-privilege macOS Seatbelt sandboxes as agent
  guardrails.
- **Multi-agent patching.** Built a multi-agent patching system where sandboxed
  workers write fixes, and a read-only verifier must pass tests before merging,
  reducing reopened fixes and speeding up patch releases by 30% compared with
  unorchestrated agents.
- **Agent observability.** Extended AI agent observability to Codex and
  OpenCode, tracking token spend and blocked sessions in the agent monitor.

#### Clinical and R&D data infrastructure

- **Submission PDF assembly.** Reduced clinical-study document preparation from
  15 minutes to under 3 minutes per file for a service with 1,000 active users
  by building a FastAPI service that combines RTF outputs into submission PDFs
  while isolating user sessions.
- **Regulated data delivery.** Delivered clinical-trial data from external
  vendors to Crinetics' biostatistics and clinical teams by building a
  regulated pipeline that processes 100-file bursts in under 2 minutes, with
  immutable GxP audit records.
- **Storage migration.** Migrated Crinetics from legacy S3 to on-prem storage,
  cutting storage costs by 30% and consolidating all vendor deliveries into a
  single audited path.
- **AWS disaster recovery.** Owned AWS disaster recovery for a GxP-validated R&D
  SaaS platform serving 1,000 users, meeting a 2-hour RTO and 4-hour RPO with a
  Terraform-built standby environment and Lambda automation that attaches
  backup EBS volumes at failover.
- **Kubernetes workspaces.** Hardened the Kubernetes engine serving 1,000
  statistical programmers, delivering a SAS workspace in minutes instead of
  hours and fixing race conditions that killed workspaces mid-setup or left
  failed ones running.
- **GxP data consolidation.** Migrated SOBI's GxP data from AWS FSx and
  OneDrive to a single Windows server, cutting storage costs by $120K annually,
  with every file verified by hash check and access controls preserved.

### Software Engineer, GameStop — Sept 2022 to May 2023

Remote

- **API latency.** Reduced tested API latency from 2 seconds to 30 milliseconds
  across 14 FastAPI routers by moving blocking database and LDAP calls off the
  event loop into a thread pool.
- **RAG chatbot.** Built an internal RAG chatbot over company pricing and
  promotions documents, using Haystack for embedding and retrieval on a
  Pinecone vector store.
- **Compliance pipelines.** Enabled compliance tracking for the 90-day legal
  limit on offer runs by building Airflow ETL pipelines that load Quicklizard
  and Zendesk API data into BigQuery.
- **Workflow automation.** Automated promotion approvals with Python, SQL views,
  and REST APIs, cutting monthly review effort by 80%.

### Software Engineer, M&G — Jan 2021 to Aug 2022

India

- **Identity and access.** Kept access aligned with each user's current role
  for 100+ users across 15 teams by running onboarding, role changes, and
  offboarding in on-premises Active Directory, with Linux server access tied
  to AD accounts.
- **Self-healing monitoring.** Automatically resolved half of incident tickets
  affecting a fleet of 150 Ubuntu servers by building monitoring and
  self-healing automation with systemd services and cron jobs.

## Projects

- **[Pigeon](https://github.com/pdwytr/pigeon)** (open source) — an observer
  that flags blocked Claude Code, Codex, and OpenCode sessions in one live
  view, so stalled agents get unblocked fast.
  → [Blog post](/blogs/pigeon-agent-console/)
- **[quickxorhash](https://github.com/pdwytr/quickxorhash)** (open source) — Go
  CLI that computes Microsoft QuickXorHash, enabling OneDrive file-hash
  verification during production migrations.

## Skills

- **Languages** — Python, TypeScript/JavaScript, SQL
- **AI / LLM & agents** — multi-agent orchestration (Docker Sandboxes), LiteLLM (Claude, OpenAI, Bedrock), RAG, Haystack, LlamaIndex, MCP, agent evals and tooling
- **Backend & APIs** — FastAPI, SQLAlchemy, REST APIs, Pydantic AI, async Python
- **Data & workflow** — PostgreSQL, Redis, pandas, Pydantic, Airflow, BigQuery, Kafka, Snowflake, pgvector, ETL/ELT
- **Cloud, DevOps & compliance** — AWS (S3, EBS, FSx, Lambda, EC2, RDS), Terraform, Kubernetes, Docker, CI/CD, GitHub Actions, GxP environments
- **Frontend & testing** — React, pytest, Locust, integration/E2E testing, AST static analysis

## Education

- **M.S. Data Science** — The University of Texas at Arlington, Dec 2023
- **B.Tech.** — Jawaharlal Nehru Technological University, Aug 2020
