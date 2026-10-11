---
title: "Is RAG Dead in 2026? The Evolution of Retrieval"
description: "Claude Code dropped its vector database. Cursor trained its own embedding model. Both work. A look at the papers, benchmarks, and shipped products behind the 'RAG is dead' claim."
date: 2026-10-11
tags: ["rag", "retrieval", "llm", "ai-engineering", "agents"]
---

Claude Code dropped its vector database. Cursor trained its own embedding model. Two of the biggest AI coding tools made opposite bets, and both of them work [1][5]. Anyone who says retrieval-augmented generation (RAG) is dead has read only half that story. So let's put the claim on trial, using the papers, the benchmarks, and what these companies actually shipped.

## A quick recap: what RAG does

A language model knows only what it was trained on. It has never read your company's documents. RAG closes that gap: before the model answers, the system fetches the relevant page and hands it over. Ask a model about your refund window without RAG and you get a confident, made-up number. With RAG, the model reads your policy and quotes it.

In 2023, nearly everyone built RAG the same way. They split documents into chunks, turned each chunk into an embedding, stored the embeddings in a vector database, and pulled back the closest few for every question. Then context windows grew from a few thousand tokens to more than a million, and people started asking why you would retrieve anything when you could paste everything in.

## Argument 1: long context makes retrieval unnecessary

### Long context really did get good

A couple of years ago, models lost track of facts buried in the middle of a long prompt. That problem is now largely fixed. OpenAI's long-context test hides eight similar items in a very large prompt and asks for one specific item. In February, Claude Opus 4.6 scored 76% at one million tokens, a large jump over its predecessor, Sonnet 4.5 [2]. Last month, GPT-6 Astra reached 96.3% in the 512,000-to-1-million-token range, up from 73.8% for GPT-5.6 Sol [3]. On this point, the "RAG is dead" camp wins.

### But a million tokens is smaller than it sounds

One million tokens is roughly 750,000 words, about the size of one large manual. A company's knowledge is spread across wikis, tickets, chat, and email, and it changes every day. It won't fit in a single prompt. Even if it did, you would have to resend it after every edit.

### And it costs too much to send every time

GPT-6 Astra charges $10 per million input tokens [3], so pasting a one-million-token knowledge base into every question costs at least $10 a question. Good retrieval sends perhaps 5,000 tokens, about 5 cents. At 10,000 questions a day, that is **at least $100,000 a day versus $500**.

Caching helps. Gemini 4 Argon discounts cached input by 95% [4]. But a cache only holds while the prompt stays the same. Change one document near the top and you pay full price to resend everything after it. Pasting everything in works for a single contract. For a company wiki that keeps growing, you still need retrieval.

## Argument 2: coding agents abandoned vector search

### The Claude Code case

Boris Cherny, who created Claude Code, said early versions used RAG with a local vector database, but the team soon found that agentic search, mostly `grep` and reading files, worked better [1]. He called it simpler and said it avoided problems with security, privacy, stale indexes, and reliability. That makes sense for code, which is full of exact names like function names and file paths. Plain keyword search finds those immediately. OpenAI's Codex relies on the same kind of search [7]. When Cherny's post spread in February, many people decided vector search was finished.

### The Cursor case

Cursor tested the question directly by running the same agent with and without semantic search. With it, accuracy rose 12.5% on average and improved on every model tested. On large codebases, more of the agent's code was kept [5]. But Cursor didn't drop `grep` either. In March, it built a dedicated index to speed `grep` up, because a single search on a very large repository could take more than 15 seconds [6]. Cursor gets its best results from using both.

### How both teams can be right

At AI Engineer Europe 2026, Turbopuffer gave Claude Code a semantic search tool and tracked which files it opened. Out of the box, about one in three files it read was wasted. Capping each read at 50 lines brought that to about one in five, and adding semantic search on top brought it to about one in eight [8]. So grep alone works, but the agent spends a lot of reads on dead ends.

Vector search alone has its own blind spots. Google DeepMind showed that on very simple queries, such as "who likes apples?", top embedding models recalled less than 20% of the right answers, while BM25, a decades-old keyword method, did much better [9]. Combining the two changes the result. Anthropic's Contextual Retrieval runs embeddings and keyword search side by side, then reranks the results, and cut failed retrievals by 67% [10]. Vectors struggle on their own but earn their place as one part of a search system.

## Argument 3: single-shot retrieval is fading

This one is mostly true. The old "one question, one lookup" setup is on its way out. Take a question like "What changed between our Q2 and Q3 pricing?" Answering it needs both quarters, but a single lookup gets one attempt to find them. What's replacing it is a loop: the model searches, reads, notices what's missing, and searches again.

Anthropic built a research system this way, with a lead agent sending out sub-agents. On internal research evaluations it beat a single agent by about 90%, but it also used roughly 15 times more tokens than a normal chat [11]. That cost is only worth paying for questions that genuinely need it.

## What enterprises are actually doing

### RAG never left

Large companies, the ones sitting on ten years of messy SharePoint, never stopped using RAG. In November 2025, Menlo Ventures surveyed about 500 enterprise AI leaders. RAG was still the second most common way to customize a model, behind plain prompting [12]. A year earlier, RAG appeared in 51% of deployments, up from 31%, while fine-tuning stood at 9% [13]. The 2025 survey also found that only 16% of enterprise deployments were true agents [12]. Online debate is fixated on agentic search, but most companies still run fixed pipelines and keep improving them.

### LinkedIn's knowledge graph

LinkedIn's customer service team turned years of support tickets into a knowledge graph so the system could see which issues were connected. Retrieval ranking improved by about 77%, and after six months in production, the median time to resolve a ticket fell by almost 29% [14].

### The MCP debate

Inside companies, the live argument is about the Model Context Protocol (MCP). Instead of copying everything into an index, the agent queries tools like Slack or Jira directly, at the moment someone asks. Search vendors push back. Glean and Elastic both argue that a live query is only as fast as the slowest app it hits, and that it can't rank results across apps [15][16]. Both companies sell indexes, so take their view with some skepticism. Still, anyone who has waited on a slow internal tool knows the speed concern is real.

### The real bottleneck: data and permissions

For most companies, the model isn't what's holding them back. Their data is. Gartner predicts that through 2026, organizations will abandon 60% of AI projects that aren't backed by AI-ready data [17].

Permissions are the other problem. Microsoft Copilot can surface anything a user is technically able to open. In a July report from CoreView, which sells Microsoft 365 governance tools, two-thirds of companies said they had delayed or cancelled Copilot over fears it would expose confidential SharePoint data [18]. For most large companies, the hard RAG work in 2026 is cleaning up data and sorting out who is allowed to see what.

## A practical decision guide

Pick the retrieval method that fits the data and the question:

| Your situation | What to use |
| --- | --- |
| Exact identifiers such as error codes, SKUs, and ticket IDs | Keyword search |
| Users describe things in different words than your documents do | Add embeddings |
| Small datasets | Load everything into the context window |
| Connected records such as tickets or customer accounts | Consider a knowledge graph |
| Hard research questions | Save the agentic search loop for these |

Most real systems mix several of these, and every one of them needs to check who is asking.

## Verdict

The 2023 pipeline, one vector database and one lookup per question, is effectively dead. Retrieval is everywhere, in coding tools and inside large companies. Claude Code went with `grep`. Cursor uses `grep` and embeddings together. LinkedIn built a graph. Every one of them is still retrieving.

If you're learning RAG now, keep going, and learn the other tools as well as vector databases. The next time someone tells you RAG is dead, ask them which RAG they mean.

---

## References

1. Boris Cherny, post on Claude Code's retrieval approach, X. <https://x.com/bcherny/status/2017824286489383315>
2. Anthropic, "Claude Opus 4.6." <https://www.anthropic.com/news/claude-opus-4-6>
3. OpenAI, "GPT-6 Astra: A new generation of intelligence." <https://openai.com/index/gpt-6-astra/>
4. eesel AI, "Gemini 4 Argon: benchmarks, pricing, and how to get access." <https://www.eesel.ai/blog/gemini-4-argon>
5. Cursor, "Improving agent with semantic search." <https://cursor.com/blog/semsearch>
6. Cursor, "Fast regex search: indexing text for agent tools." <https://cursor.com/blog/fast-regex-search>
7. Yage, "Why Coding Agents Still Use grep as Their Search Backbone." <https://yage.ai/share/why-coding-agents-still-use-grep-en-20260327.html>
8. Kuba Rogut (turbopuffer), "Benchmarking semantic code retrieval on Claude Code," AI Engineer Europe 2026. <https://ai.engineer/talks/zKk7sDMGDEQ>
9. Google DeepMind, "On the Theoretical Limitations of Embedding-Based Retrieval," arXiv:2508.21038. <https://arxiv.org/abs/2508.21038>
10. Anthropic, "Introducing Contextual Retrieval." <https://www.anthropic.com/engineering/contextual-retrieval>
11. Anthropic, "How we built our multi-agent research system." <https://www.anthropic.com/engineering/multi-agent-research-system>
12. Menlo Ventures, "2025: The State of Generative AI in the Enterprise." <https://menlovc.com/perspective/2025-the-state-of-generative-ai-in-the-enterprise/>
13. Menlo Ventures, "2024: The State of Generative AI in the Enterprise." <https://menlovc.com/2024-the-state-of-generative-ai-in-the-enterprise/>
14. LinkedIn, "Retrieval-Augmented Generation with Knowledge Graphs for Customer Service Question Answering," arXiv:2404.17723. <https://arxiv.org/abs/2404.17723>
15. Glean, "Is MCP + federated search killing the index?" <https://www.glean.com/blog/federated-indexed-enterprise-ai>
16. Elastic, "The future of search engines: Does MCP make indexed search obsolete?" <https://www.elastic.co/search-labs/blog/future-of-search-engines-indexed-search-mcp>
17. Gartner, "Lack of AI-Ready Data Puts AI Projects at Risk." <https://www.gartner.com/en/newsroom/press-releases/2025-02-26-lack-of-ai-ready-data-puts-ai-projects-at-risk>
18. CoreView, "66% of Enterprises Delay Microsoft Copilot Over Concerns AI Could Expose Confidential SharePoint Data." <https://www.coreview.com/news/66-percent-of-enterprises-delay-microsoft-copilot>
