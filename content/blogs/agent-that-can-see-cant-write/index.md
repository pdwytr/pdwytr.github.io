---
title: "The Agent That Can See Can't Write"
description: "A pattern for scaling enterprise multi-agent systems: split the agent that can see the real system from the agent that writes code, and a six-month project on an audited clinical data service ships in three."
date: 2026-10-25
build:
  publishResources: false  # diagrams are inlined; keep their sources off the site
tags: ["ai-agents", "architecture", "stacked-prs", "audit", "separation-of-duties"]
---

We estimated six months for a project on a clinical data service, where a bad change can become an audit finding. We shipped it in three, with the same review and testing standards.

The speed came from how the agents were arranged, and the arrangement is old news to anyone who has worked near auditors: separation of duties. The agent that can see the real system can't write code. The agent that writes code can't touch the real system.

---

## The project

The system is a file delivery engine we built for a pharmaceutical client. It moves clinical trial data from outside vendors into downstream systems. A watcher notices each file that lands in an inbound folder and records it in an append-only journal. An orchestrator reads that journal and hands each file to a distributor for delivery. Around that core sit an ingestion job that pulls EDC deliveries, a log monitor, a study notifier, and a folder cleaner.

It's a validated system, with installation and operational qualification documents, and the journal *is* the audit trail. If the engine loses a file, delivers one twice, or can't say what happened to one, that's a finding.

The codebase already had the kind of boundaries this pattern needs:

- **No component imports another.** Shared code goes through one common library, so each component is a natural blast radius.
- **Contracts are machine-checked.** The data passed between components is defined by versioned JSON Schemas, and changing one means bumping its version.
- **Invariants are written down as rules.** Only the watcher writes to the journal, and only the orchestrator writes execution state. Every timestamp is UTC, and every file is identified by its SHA-256.
- **Some paths need a human.** The original requirements, the decision records, and operational state files can't be changed without explicit human approval in the PR.
- **Every component has its own CI workflow.**

The rest of this article describes the agent architecture we ran on top of those boundaries.

---

## The problem

AI coding agents are fast, and in a high-integrity codebase that speed comes with three problems.

A long-running agent collects history, half-finished ideas, and stale assumptions, and it reasons worse as the pile grows. Ask it for one fix and it may edit files nobody mentioned. And when the diff comes back large, nobody can review it seriously, so it gets rubber-stamped or thrown away.

Most teams respond by keeping agents away from serious work. In audit systems that's the default, because one unreviewed change can turn into a finding.

This pattern doesn't try to find an agent trustworthy enough for serious work. It limits what any single agent is allowed to do, so the trust comes from the structure rather than the model.

---

## The core principle: split seeing from acting

Each unit of work, called a node, is run by agents with opposite permissions: one **Verifier** and one **Builder** per issue.

| | **Verifier** | **Builder** |
|---|---|---|
| Reads the real system | ✅ | ❌ |
| Runs tests | ✅ Integration and stack tests, on staging | ✅ Unit tests only, inside the sandbox |
| Writes code | ❌ | ✅ |
| Where it runs | Against the real environment, read-only | Only inside a Docker sandbox |
| Where its changes go | Nowhere: it approves, rejects, and merges approved work | Its own git worktree, mounted into its sandbox |

The Verifier has knowledge but no hands, and the Builder has hands but no access, so neither can do damage alone. A Builder that goes off the rails breaks only its own sandbox and worktree. A Verifier that misreads something can wrongly approve or reject, but it can't write the bad change itself.

Auditors will recognize the control. It's the same one that stops a single person from both approving an invoice and paying it.

Every diagram in this article uses the same colors: **blue** is the Verifier, **orange** a Builder, **yellow** an issue branch, **purple** context, **green** a PR, **teal** CI and staging, **pink** a human, and **red** a failure. In the stack diagram, **indigo** boxes are nodes.

Here is what happens inside one node:

{{< diagram src="diagrams/01-node.svg" alt="Inside one node: the Verifier reads the node context and gives issues 104, 105 and 106 to three Builders. Each Builder commits to its own branch in its own worktree, sends its diff and test results back to the Verifier, and approved branches are merged into the node's single PR." >}}

*Inside one node. The node has one context (purple) and produces one PR (green). The Verifier (blue) reads the context and gives each issue to its own Builder, named after the issue number. Each Builder (orange) runs in its own Docker sandbox and commits to its own branch (yellow), checked out in its own worktree. The Verifier reviews each branch's diff and test results, and merges the approved branches into the node branch. That node branch is PR #2, so one PR carries all three issues. The Verifier performs the merges but never writes code.*

---

## The architecture, piece by piece

The design has eight parts, and each one closes off a specific way things go wrong.

### 1. Stacked PRs: units of context and units of review

The work is split into nodes, and each node produces one pull request. When a node creates a child node, two things happen together: the child's PR is stacked on the parent's PR branch, and the child receives a distilled version of the parent's context. Nodes with no children are the leaf nodes, the tips of the stack.

The relationships are fixed:

| Relationship | Rule |
|---|---|
| Node → PR | 1 to 1 |
| Node → context | 1 to 1 |
| Issues → node | many to 1 |
| Issues → PR | many to 1 (through the node) |
| Issue → Builder → worktree | 1 to 1 to 1 |
| Child PR → parent PR | stacked on the parent's branch |

{{< diagram src="diagrams/02-stack.svg" alt="The stack: main is the base of Node A. Node A passes context and its branch to Nodes B and C, and Node B passes them to Nodes D and E. C, D and E are leaf nodes." >}}

*The stack. Each arrow carries the parent's distilled context to the child and makes the parent's PR branch the base of the child's PR. Leaf nodes have a green border. The diagram above shows what happens inside each box.*

Because of this, no node needs the whole project history, only the context its parent distilled for it. Reviewers get the same benefit: one small PR per node instead of one large diff.

### 2. Nodes: one Verifier, one Builder per issue, one PR

A node owns a group of related issues and produces exactly one PR. Inside it, one Verifier directs one Builder per issue, and each Builder is named after its issue (`#104`, `#105`, ...).

The Builders work in parallel. The Verifier decides when each issue is done and merges the approved work into the node's PR branch. The node is finished only when the merged result passes against the real system.

Merging approved branches isn't writing code. The Verifier never writes a line; it only decides which Builder work gets into the PR.

### 3. Links between nodes are context pipelines

Nodes don't share one long conversation. Each node gets exactly one context when it is created, distilled by its parent's Verifier from what the parent learned:

- what the parent delivered: interfaces, contracts, and the decisions behind them
- the test results that show it works
- the child's blast radius, its resource budget, and the constraints it must respect

Dead ends, abandoned approaches, and chat history stay behind. A node deep in the stack starts with a short, deliberate context instead of everything that happened above it.

Inside a node, the context is cut down once more: the Verifier gives each Builder its own issue and scope, not the whole node context.

### 4. Worktrees: every change is isolated

Every Builder gets its own git worktree. The worktree lives in the host repo, and only that worktree is mounted into the Builder's sandbox. The Builder can't see the rest of the host, but the Verifier can see every worktree, which is what lets it review and merge.

Isolating changes this way buys three things. Builders in the same node, and independent nodes in the stack, can work at the same time without stepping on each other. A failed attempt is deleted with its worktree, so nothing leaks into the main checkout. And every diff traces back to exactly one issue, one node, and one PR.

### 5. Blast radius: a declared scope for every node

Before a node starts, it gets a **blast radius**: the files, modules, and interfaces it is allowed to change. Anything outside that scope is rejected, however reasonable it looks.

The Builder knows its boundary before it writes anything. The Verifier checks every diff against the boundary before it looks at anything else. A human reviewer knows what a PR *could* have touched, so the review starts from a known edge.

When a node finds it really does need to change something outside its radius, it doesn't make the change. It escalates, and the stack is re-planned. Scope creep becomes a decision someone makes on purpose instead of a side effect nobody noticed.

{{< diagram src="diagrams/03-blast-radius.svg" alt="The Verifier's check on a submitted diff: outside the blast radius, it is rejected and the stack is re-planned; inside but failing unit tests, it goes back to the Builder with evidence; otherwise it is merged into the node's PR branch." >}}

### 6. Testing: three gates, split by who can reach what

The sandbox that makes a Builder safe also limits what it can test. A Builder can't reach a database, a queue, or another service, so it can prove its change only with unit tests. Integration testing needs a real environment, which puts it on the Verifier's side of the line.

So testing follows the same split as everything else: Builders write the tests, and the Verifier runs the ones that need a real system.

{{< diagram src="diagrams/04-gates.svg" alt="Three test gates in order: unit tests per issue in the Builder's sandbox, integration tests per node on staging run by the Verifier, then after human review, stack tests on staging before merging into main." >}}

**Gate 1: unit tests, per issue.** Each Builder writes its code, its unit tests, and any integration tests the issue needs. It runs the unit tests in its sandbox with dependencies mocked, and nothing reaches the node's PR branch until they pass.

**Gate 2: integration tests, per node.** After the Verifier merges the approved issues, CI deploys the node's PR branch to an ephemeral staging environment. Because the PR is stacked, the branch already contains every ancestor's changes. The Verifier runs the integration suite there, against real services and realistic data. This is the first time the node's issues are tested together, and the first time the Builders' integration tests actually run.

**Gate 3: stack tests, before merging into main.** PRs merge from the root down. Before one merges, it's rebased onto the current `main`, which by then includes its merged ancestors and any merged siblings, redeployed to staging, and run through the full integration and end-to-end suite again. Sibling nodes never saw each other's work, so this is where conflicts between them show up.

When a gate fails, the Verifier still doesn't fix anything. It traces the failure to the issue whose blast radius covers the failing code and sends the logs to that Builder. If the failure spans issues or breaks a contract, it escalates to the parent node.

Tests run on staging, never production. The Verifier reads production to understand it, but an integration test that writes data is still a write, and production is where an audit system controls writes most strictly.

### 7. Resource budgets: a blast radius for compute

The blast radius limits *what* a node can change. A resource budget limits *how much* it can consume. A parent sets each child's budget when it creates the child, alongside the blast radius, and the budget goes into the child's context.

| Resource | Limited where | What it prevents |
|---|---|---|
| **Compute** (CPU, memory, processes) | Each Builder's sandbox, and the node's staging runs | A runaway loop or build starving every other node |
| **Caches** | A shared, read-only base cache, plus a separate writable cache for each node | One node corrupting another node's build cache, or depending on cache state no other node has |
| **Storage** | Disk quota on each sandbox and worktree | Generated files, logs, or test data filling the host |
| **Time** | Each Builder attempt, each test run, and the node as a whole | An agent retrying forever, or a node that never finishes |

When a Builder hits a limit, its sandbox is stopped and the attempt fails. The Verifier gets the evidence and either sends the issue back with a hint or escalates. When a whole node runs out of time or budget, it escalates to its parent. A node that can't fit its budget usually has the wrong scope, and the fix is to split or re-plan it rather than quietly raise the limit.

Budgets also make cost easy to account for. Every unit of compute, storage, and time belongs to one node, one PR, and one issue, which is the trail an audit system wants anyway.

### 8. Tool access: each role gets only its own tools

Separation of duties holds only if the tools enforce it. Every role gets an allowlist, and anything not on it doesn't exist for that role.

| Role | Can use | Can never use |
|---|---|---|
| **Verifier** | Read-only access to the real system (scoped read-only credentials for logs and queries) · read the repo and every worktree (`git log`, `diff`, `show`) · create worktrees and branches · merge approved issue branches into its node's PR branch · push that one branch · trigger staging deploys and test runs · read test results · open and comment on its PR · create child nodes and write their context | File edits · shell commands that change anything on the host · pushing to `main` · writes to production |
| **Builder** | Read and write files in its own mounted worktree · a shell inside its sandbox · the unit test runner · package installs from the read-only cache or an approved mirror · local commits in its worktree | Network access beyond the package mirror · the real system · other Builders' worktrees · pushing · opening PRs · staging |
| **CI + staging** (automation, not an agent) | Build the PR branch · deploy it to an ephemeral staging environment · run integration and stack tests · report results | Deploying to production · merging anything |
| **Human reviewer** | Approve PRs · merge into `main` · change contracts, blast radii, and budgets · resolve escalations | Nothing is technically blocked, but every action is logged like everyone else's |

The allowlists are enforced twice. The agent harness gives each role only its own tools, so it can't even ask for anything else. The infrastructure enforces the same limits underneath: Docker runs with networking off and one worktree mounted, production credentials are read-only, and branch protection blocks pushes to `main`.

The second layer is the one that matters in an audit. A prompt can't grant a permission the infrastructure doesn't have, so a confused or manipulated agent still stays inside its role.

---

## How it fits together

Here is one node from start to finish, with two of its issues:

{{< diagram src="diagrams/05-sequence.svg" alt="Sequence of one node: the Verifier hands issues 104 and 105 to two Builders in parallel, each returns a diff and unit test results, the Verifier merges into PR 2 and runs integration tests on staging, sends any failure back to its Builder, then sends the PR with evidence to a human reviewer who approves it." >}}

*One node from start to finish. In the build phase, each Builder works on its own issue and returns a diff with unit test results. In the integrate phase, the Verifier checks scope, merges the work into the node's PR, runs integration tests on staging, and sends any failure back to the Builder that owns it. In the review phase, a human approves the PR. The node starts from the context its parent gave it and ends by distilling the context its own children will start from.*

---

## System scalability

When a single agent does the work, it carries more as the project grows: a longer context, more of the codebase, bigger diffs. This system grows by adding nodes instead. Each node keeps its own context, sandboxes, worktrees, and budget, so a new node adds capacity without adding load to the nodes already running.

### What stays constant as the system grows

Three things don't grow with the project, and they're the reason it scales:

- **Context per node.** A node sees only what its parent distilled, so a leaf in a 50-node project starts with about as much context as a leaf in a 5-node one.
- **PR size.** One node produces one PR, so review effort grows with the number of nodes, not with the size of each diff.
- **Risk per change.** Every node has a blast radius and a budget, so a bigger system has more boundaries, not looser ones.

### Three ways to scale

| Direction | How it scales | What makes it safe |
|---|---|---|
| **Within a node** | More issues → more Builders, all running at once | Each Builder has its own sandbox and worktree, and the Verifier merges their work |
| **Across the tree** | More sibling nodes, all running at once | Siblings share a parent's contract, not each other's work, and Gate 3 catches conflicts at merge time |
| **Down the tree** | Parent and child nodes overlap instead of waiting for each other | Contract-first context, explained below |

The first two come with the architecture. The third needs one change.

### Scaling down the tree: contract-first context

If a child has to wait for its parent to *finish* before its context exists, every path from root to leaf runs one node after another. The longest of those paths, A then B then D in the stack above, sets the minimum project time however many agents you add.

The change is in **when a child's context is created**:

1. **Write context from the contract, not the result.** When a node creates its children, its Verifier writes their context from the *planned* interfaces, decisions, blast radius, and budget, so the children can start right away.
2. **Branch from the parent's current state, then restack.** A child's PR is based on its parent's branch as it stands now. When the parent's branch moves, the child rebases onto it and re-runs Gate 2.
3. **Treat a contract change as an escalation.** If a parent has to change a contract it already handed down, the affected children get a fresh context and re-run. Without that rule, parallel work goes on building against a contract that no longer exists.
4. **Keep merging in order.** Building overlaps, but merging into `main` still goes root first, through Gate 3, so review stays as simple as before.

{{< diagram src="diagrams/06-timeline.svg" alt="Gantt chart comparing a root-to-leaf path of Nodes A, B and D run one after another with the same nodes starting from the contract and overlapping, followed by a short restack-and-retest step." >}}

*An illustration, not measured data: the same root-to-leaf path run both ways. Starting from the contract lets the nodes overlap, at the cost of a restack-and-retest step at the end.*

### Where it stops scaling

| Bottleneck | Why it appears | What to do |
|---|---|---|
| **The Verifier** | One Verifier reviews every diff in its node | Cap the number of issues per node. When a node outgrows the cap, split it into child nodes |
| **Contract churn** | Every contract change sends rework down the tree | Settle contracts early, at the nodes near the root, and keep them small |
| **Staging capacity** | Every node runs integration tests, and every restack runs them again | Use ephemeral staging environments and queue test runs within each node's budget |
| **Human review** | People can only review so many PRs a day | Keep PRs small and attach test evidence, so each review is quick |

Budgets make capacity planning a matter of addition. The system's total compute, storage, and time is the sum of its nodes' budgets, so you can decide how many nodes run at once before anything starts, instead of finding out from a slow build or a full disk.

Overlapping nodes pays off only when contracts are stable. When they keep changing, the overlap turns into rework. Audit systems tend to design their interfaces carefully up front, which is why the overlap works there.

---

## Observability

Many agents working in parallel are manageable only if you can see what each one did and why. Here the structure does most of that work, because it already gives every action an owner.

Every event carries the same four IDs: node, issue, PR, and attempt. A failed integration test, a merge, a budget alert, and a Builder's shell command all trace back to one issue in one node.

{{< diagram src="diagrams/07-observability.svg" alt="Verifier, Builders, and CI plus staging all send events to a telemetry collector, which feeds an append-only audit log, dashboards, and alerts. The Builders' events are recorded by the sandbox runtime." >}}

*Builders don't report on themselves. The sandbox runtime records their tool calls and output from outside the sandbox, so a Builder can't edit or hide its own record.*

| Signal | Examples | Used for |
|---|---|---|
| **Audit log** | Every tool call, merge, approval, escalation, and context handed to a node | Answering "who changed this, and on what basis?" months later |
| **Traces** | One trace per node, with a span per Builder attempt and per gate | Seeing where time goes and which path through the tree is the slowest |
| **Metrics** | Attempts per issue · gate pass rates · budget used vs. allowed · escalations · time waiting for review | Spotting badly scoped nodes and bottlenecks before they slow everything down |
| **Alerts** | Budget nearly spent · repeated failures on one issue · blast radius violations · contract changes | Bringing a human in only when something needs one |

Each node's context is saved as a versioned snapshot too. When a Builder makes a strange choice, you can read exactly what it was told instead of guessing.

---

## A maintainable system, not just a fast one

The same structure that makes the work fast also makes the result easier to live with afterward:

- **Every line has a recorded reason.** Code traces to an issue, the issue to a node, the node to a PR and the context it was given, so "why is this here?" has an answer you can look up.
- **Modules stay decoupled.** Blast radii keep each change inside declared boundaries, so a dependency between modules appears only when a contract calls for it.
- **Contracts double as documentation.** The interfaces and decisions passed down the tree are written down by design and kept with the code they describe.
- **Tests arrive with the code.** Every issue ships with unit tests and every node with integration tests, so coverage grows with the codebase instead of trailing it.
- **History is easy to read and to undo.** One small PR per node, merged in order, means any node can be reviewed, reverted, or rebuilt on its own.
- **The node tree is a map.** A new engineer can see which part was built when, against which contract, and with which tests.
- **Old changes can be rebuilt.** Sandboxes, pinned caches, and versioned context snapshots let you rebuild and re-test a change as it was.

In an audit system, that's the part that lasts: the next audit, feature, and engineer inherit a codebase that explains itself.

---

## Where the time went

On paper this is more process than one agent with a big prompt: more agents, more handoffs, more checks. The project still took half its estimate, and the structure explains why.

- **Out-of-scope work is caught when the diff is submitted.** The blast radius check is mechanical, so a wandering change is rejected before anyone spends review time on it.
- **Reviews are small and come with evidence.** A reviewer reads one node's PR with its test results attached, not a sprawling diff.
- **Work overlaps at every level.** Builders within a node, sibling nodes, and, with contract-first context, parents and children all run at the same time.
- **Failures surface where they're cheapest.** Unit tests run in the sandbox before anything reaches staging, and integration runs happen once per node, not once per change.
- **Reviewers can trust the boundary.** They don't have to wonder what else an agent touched, because it couldn't touch anything else.

The result, on a production audit system:

- **Estimated timeline:** 6 months
- **Actual timeline:** 3 months
- **Review and testing standards:** unchanged

---

## What broke

Halving the timeline didn't mean nothing went wrong. In production, two vendor files landed in the inbound folder and were never picked up.

The vendors upload to a temporary name, then rename the file to its final name. An earlier fix had taught the watcher to treat that rename as the arrival signal, and its tests passed. But the production SFTP server performed the rename as *create a link, then delete the old name*, so the operating system reported a new file being created rather than a rename. The watcher deliberately ignored that kind of event. The fix's tests had only exercised the other style of rename.

No data was lost, and a service restart delivered both files the same morning. Within a day there was a confirmed root cause, a written root-cause analysis, and a client summary. The corrective actions were to recognize both event types, test both rename styles against the production server's real behavior, re-scan the inbound folder periodically, and alert when a file sits unprocessed.

The lesson fits the architecture. Unit tests in a sandbox can only prove what their mocks believe, and here the mocks believed in the wrong rename. That's the argument for Gate 2: integration tests have to run somewhere that behaves like production, not like the test author's idea of it.

---

## When not to use this

This pattern pays off when mistakes are expensive. It's overkill for prototypes and throwaway scripts, for small single-file changes, and for greenfield exploration, where the shape of the work isn't known yet, and a fixed blast radius would get in the way.

If a bad merge costs you only a revert, one agent with good tests is probably enough. If it costs you an audit finding, don't look for an agent you can trust with everything. Build the system so that no agent ever has to be.
