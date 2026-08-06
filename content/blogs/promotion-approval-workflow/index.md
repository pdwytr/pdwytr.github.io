---
title: "Promotion Approval Workflow"
description: "An Airflow-orchestrated approval pipeline with database-backed tracking that cut team effort 80% and about $200K a month in leakage."
date: 2026-05-01
tags: ["case-study"]
period:
  start: "2022-09"
  end: "2023-05"
role: "Data Engineer"
org: "GameStop"
stack: ["Python", "Apache Airflow", "SQL", "REST APIs"]
metrics:
  - label: "team effort"
    value: "−80%"
  - label: "monthly saving"
    value: "~$200K"
featured: true
weight: 30
draft: false
---

## Problem

Promotions went through a manual approval process, and promotion-duration limits — how long a discount was allowed to stay active — weren't systematically enforced. That combination meant approvals were slow and inconsistent, and promotions could run longer than intended without anyone catching it until the impact had already been felt.

## Constraints

The workflow had to fit into systems that already existed rather than replace them: approvals, pricing, and promotion data all lived in other tools, so the new workflow had to integrate with them over their existing REST APIs rather than becoming a new system of record that everything else had to migrate to.

## Design

Apache Airflow DAGs orchestrate the approval workflow end to end, calling out to the existing systems' APIs at each step rather than duplicating their data. Workflow state — where a given promotion request is in the approval process — is tracked in a database, which also feeds dashboards so approvers and stakeholders can see the state of pending and active promotions without chasing status by hand. That same tracking layer enforces controls around promotion-duration limits, so a promotion can't silently keep running past its approved window.

## Outcome

The workflow cut the team's manual effort on approvals by 80% and reduced promotion-related leakage — spend lost to promotions that ran longer or looser than intended — by about $200K a month.
