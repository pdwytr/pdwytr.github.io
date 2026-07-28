---
title: "Enterprise Document Conversion API"
description: "An eight-endpoint async service that took per-file processing from over fifteen minutes to under thirty seconds for about a thousand users."
date: 2026-07-01
period:
  start: "2025-07"
  end: "2026-07"
role: "Backend Engineer"
org: "Takeda"
stack: ["Python", "FastAPI", "asyncio", "pytest", "OpenAPI"]
metrics:
  - label: "processing time"
    value: "−97%"
  - label: "active users"
    value: "~1,000"
  - label: "concurrent sessions"
    value: "100"
  - label: "cross-session leaks"
    value: "0"
featured: true
weight: 10
draft: false
---

## Problem

Document conversion ran synchronously: a caller uploaded a file, held the connection open, and waited for the conversion to finish before getting a response. For anything but the smallest files that meant blocking for more than fifteen minutes per file, which made the API unusable for any workflow that needed to convert more than one document at a time or run unattended.

## Constraints

This was an enterprise, multi-tenant context, so a few things were non-negotiable rather than nice-to-haves:

- **No cross-session data visibility.** Different callers' files and results could never be visible to one another, even transiently, even under concurrent load.
- **Provable file integrity.** A file that was corrupted or tampered with in transit had to be caught, not silently converted.
- **A pollable contract, not a held connection.** Clients needed to submit a job and check on it later — they couldn't be expected to keep a connection open for the duration of a long-running conversion, and infrastructure in between (load balancers, gateways, proxies) couldn't be trusted to hold one open reliably either.

## Design

The service exposes eight endpoints around a bearer-token session model: each session gets its own isolated storage, so one caller's files are never reachable through another caller's token. On upload, the service runs an integrity check before anything else touches the file. Conversion itself happens asynchronously — the upload endpoint returns immediately with a handle, and the actual work runs in the background behind a polling endpoint the client checks until it reports done, failed, or still-processing. Each failure mode — bad upload, integrity check failure, conversion error, timeout — has its own explicit error contract in the OpenAPI spec, rather than a generic 500 that leaves the caller guessing.

## Outcome

Processing time per file dropped 97%, from more than fifteen minutes down to under thirty seconds. The service supported about 1,000 active users and 100 concurrent sessions with zero cross-session data leakage incidents. What I'd revisit: the polling interval is currently a fixed client-side backoff rather than a server-suggested `Retry-After`, which would let the service signal expected completion time instead of making every client guess at a poll cadence.
