# moknshaik.com — Personal Site Design

**Date:** 2026-07-27
**Repo:** `/Users/khalid/Documents/Projects/website`
**Domain:** `moknshaik.com` (GitHub Pages, custom domain already configured)
**Status:** Approved design, pending implementation plan

---

## 1. Purpose

A personal site that does three jobs, in priority order:

1. **Publish writing** with near-zero friction — markdown (or HTML) plus pictures, committed to this repo.
2. **Substantiate the résumé** — turn résumé claims into linked case studies a recruiter can verify in one click.
3. **Be genuinely useful** to a reader who arrives from a Reddit, LinkedIn, or X link.

Success criteria:

- Publishing a post is: create file → write → preview locally → `git push`. Live in about a minute.
- Adding a picture to a post is: drop the file next to the post → one line of markdown.
- A shared link renders a proper social card with title and site identity.
- The deployed site is HTML and CSS. No client-side framework, no runtime to break.
- Toolchain is two compiled binaries. No `package.json`, no lockfile, no dependency-update treadmill.

Non-goal: this is not a web app. Nothing here requires a server.

---

## 2. Stack

**Hugo** (Go), static output, deployed to GitHub Pages by GitHub Actions.

| Concern | Decision |
|---|---|
| Site generator | Hugo, pinned version, no third-party theme |
| Content format | `.md` or `.html` — both are native Hugo content formats |
| Images in content | Page bundles + an image render hook (§6) |
| Syntax highlighting | Chroma (built into Hugo), class-based, dual light/dark |
| Taxonomies, RSS, sitemap, pagination, reading time | Hugo built-ins, no configuration beyond enabling tags |
| Social/OG cards | Generated at build time with Hugo's `images.Text` filter |
| Frontmatter validation | Validation partial that calls `errorf`, failing the build |
| Styling | Hand-written CSS, concatenated + minified + fingerprinted via Hugo Pipes |
| Fonts | Self-hosted variable woff2 (Inter for text, JetBrains Mono for metadata/code) |
| Build verification | A small stdlib-only Go program in `tools/verify` |
| Link checking | `lychee` (Rust binary) in CI |
| Comments | None — `discuss` frontmatter links to the Reddit/LinkedIn/X thread instead |
| Analytics | Deferred (GoatCounter is a one-line add later) |
| Search | Deferred (Pagefind ships a standalone binary; add past ~20 posts) |

### Why Hugo over Astro

Astro's real advantages over Hugo are: friendlier templating, Zod-typed frontmatter, and first-class interactive islands. Weighed against this project:

- The friendlier templating is JSX-flavored, which is the thing being avoided by choosing Go.
- Typed frontmatter is replaceable with a ~20-line validation partial (§5.3).
- Interactivity is speculative; the described site is essays, case studies, and pictures.

Against that, Hugo gives a one-binary toolchain, native `.html` content, and taxonomies/feeds/pagination/image-processing already built. It covers 100% of the requirements.

**Accepted cost:** layouts are written in Go `html/template`, which is verbose and has fiddly lookup rules. Mitigated by owning only ~8 small layout files and no theme.

### Deliberately out of scope (YAGNI)

Comments, search, newsletter, analytics, tag-filtering UI, dark/light beyond a single toggle, pagination tuning, i18n, webmentions, a CMS admin UI. Each has a named on-ramp in §11; none is built now.

---

## 3. Information architecture

```
/                      identity, what I work on, 3 recent posts, 3 pinned projects
/writing/              all posts, reverse-chron, year-grouped
/writing/<slug>/       a post
/projects/             case-study index with metrics visible
/projects/<slug>/      a case study
/resume/               full résumé as HTML + PDF download
/about/                bio, contact, links
/tags/<tag>/           tag archive
/rss.xml               feed (Hugo built-in)
/sitemap.xml           sitemap (Hugo built-in)
/og/<slug>.png         generated social card per post
/404.html
```

`posts` and `projects` are two Hugo sections sharing one rendering pipeline, differing only in frontmatter schema and index layout. This is the central simplification: one prose renderer, two content shapes.

---

## 4. Visual design

Direction: **technical / engineer-coded**, with three deliberate constraints.

1. **Monospace for metadata, labels, dates, and code — never for body copy.** Long-form essays set in mono are hostile to read; that is the characteristic failure of developer blogs.
2. **Dark by default, `prefers-color-scheme`-aware, with a toggle.** A recruiter opening the site in light mode on a managed laptop still sees something composed.
3. **Restrained, not neon-terminal.** Two accent colors maximum, hairline rules, generous whitespace.

### The `$ ls` device

The writing index carries a file-listing flavour: a `$ ls` ornament in the section header, monospace date and kind columns in fixed alignment.

**Refinement over the original sketch:** the sketch showed filenames (`ai-leverage.md`) as the row label. Post *titles* must be the visible row label — they are what a reader skims and what search engines index. The listing aesthetic is carried by monospace metadata and column alignment, not by hiding titles.

```
$ ls writing/

2026-07-27   The Real Leverage Problem                    [essay]
2026-07-12   Async FastAPI Under Real Load                [note]
2026-06-30   Reconciling Postgres State With Kubernetes   [essay]
```

### Design tokens

A single `tokens.css` defines: `--bg`, `--bg-elev`, `--fg`, `--fg-muted`, `--accent`, `--accent-2`, `--border`, `--code-bg`, plus the type scale, measure, and spacing steps. Both themes are defined here and nowhere else, so a palette change is one file.

### Typography

- Text: Inter Variable, self-hosted, latin subset, `woff2`, preloaded.
- Mono: JetBrains Mono, self-hosted, latin subset.
- No external font requests. No layout shift.

### Client-side JavaScript

Exactly one inline script, roughly 15 lines with no dependencies: read the stored theme preference and set `data-theme` on `<html>` before first paint to avoid a flash, and toggle it on click. Everything else on the site is static HTML and CSS.

---

## 5. Content model

### 5.1 Posts

`content/posts/<slug>/index.md` (or `index.html`)

| Field | Type | Required | Notes |
|---|---|---|---|
| `title` | string | yes | |
| `description` | string | yes | Used for `<meta name=description>`, OG, and the index blurb |
| `date` | date | yes | |
| `lastmod` | date | no | Rendered as "updated" when present |
| `tags` | []string | no | Drives `/tags/<tag>/` |
| `kind` | `essay` \| `note` | yes | The `[essay]` / `[note]` label in listings |
| `cover` | string | no | Filename within the bundle; used for OG card and post header |
| `discuss` | []{label, href} | no | Links to the Reddit / LinkedIn / X thread after sharing |
| `draft` | bool | no | Defaults false |

`discuss` is how the site gets a comment section without hosting one: after sharing a post, the thread URL is added to frontmatter and the post links out to it.

### 5.2 Projects

`content/projects/<slug>/index.md`

| Field | Type | Required | Notes |
|---|---|---|---|
| `title` | string | yes | |
| `description` | string | yes | |
| `date` | date | yes | Ordering fallback |
| `period` | {start, end} | yes | `end` omitted means present |
| `role` | string | yes | |
| `org` | string | no | Omit where a client engagement shouldn't be named |
| `stack` | []string | yes | Rendered as monospace chips |
| `metrics` | []{label, value} | no | The numbers, called out visually |
| `links` | []{label, href} | no | Repo, docs, write-up |
| `featured` | bool | no | Surfaces on the home page |
| `weight` | int | no | Manual ordering on the index |
| `draft` | bool | no | |

Case-study body structure (convention, enforced by the archetype, not by code): **Problem → Constraints → Design → Outcome**.

### 5.3 Validation

`layouts/_partials/validate.html`, called from `baseof.html`, checks the required fields for the page's section and calls `errorf` with the offending file path when one is missing or malformed. Hugo treats an ERROR-level log as a build failure, so a bad post fails `hugo server` immediately rather than shipping broken.

CI additionally runs with warnings escalated to failures, so a silently-degraded build cannot deploy.

---

## 6. Images

Images live **inside the post's own directory** (a Hugo page bundle):

```
content/posts/leverage-problem/
├── index.md
├── hero.png
└── reconciler-states.png
```

and are referenced with ordinary markdown:

```markdown
![Reconciler state machine](reconciler-states.png "Fig 1 — states and transitions")
```

Hugo would otherwise copy the file through untouched. One image render hook, written once, makes every `![...]()` in every post automatically resized, converted to WebP, dimensioned (so the page does not reflow), lazy-loaded, and captioned:

```go-html-template
{{/* layouts/_markup/render-image.html */}}
{{- with .Page.Resources.GetMatch .Destination -}}
  {{- $img := .Resize "1200x webp q82" -}}
  <figure>
    <img src="{{ $img.RelPermalink }}" width="{{ $img.Width }}" height="{{ $img.Height }}"
         alt="{{ $.Text }}" loading="lazy" decoding="async">
    {{- with $.Title }}<figcaption>{{ . }}</figcaption>{{ end -}}
  </figure>
{{- else -}}
  <img src="{{ .Destination }}" alt="{{ .Text }}" loading="lazy">
{{- end -}}
```

The `else` branch passes external URLs through unchanged. Output format is WebP; Hugo does not encode AVIF, which is immaterial at this scale.

**Authoring flow:** drag a screenshot into the post folder, write one line of markdown.

---

## 7. Social cards

Every post and project gets a generated 1200×630 PNG at `/og/<slug>.png`, produced at build time by overlaying the title and site identity onto a dark base card using Hugo's `images.Text` filter. Inputs: `assets/og/base.png` and a TTF font resource in `assets/fonts/`.

This matters more than it looks: it is the difference between a link that reads as considered and one that reads as a naked URL when shared to LinkedIn or Reddit.

`<head>` also carries: canonical URL, description, `og:*`, `twitter:card=summary_large_image`, and JSON-LD — `BlogPosting` on posts, `Person` on the home and about pages so search engines attribute the work correctly.

---

## 8. The legitimacy mechanism

A résumé page alone does not build credibility. The résumé makes verifiable-sounding claims — *97% faster, ~1,000 concurrent users at 99.9% availability, ~$200K/month saved*. Each significant claim becomes a case study that shows the reasoning, and **the résumé line links to it.** A reader goes claim → constraints → design decision → outcome without searching.

Initial case studies, ported from the résumé:

1. Enterprise document conversion API — 8 endpoints, async processing, 15 min → under 30 s.
2. Managed analytics platform — Postgres-to-Kubernetes state reconciliation, provisioning days → minutes.
3. Promotion approval workflow — Airflow-orchestrated, 80% effort reduction.
4. Production LLM agent runtime — write-safety, canonical-hash trust verification, AST-based CI gate against stub implementations.

### Résumé handling

- `/resume/` renders `content/resume.md` as HTML — skimmable in-browser and indexable.
- `static/mohammed-khalid-shaik-resume.pdf` is the download.
- **The PDF is copied in manually.** The `job-search` repo is private and holds per-application tailored variants; syncing from it automatically would publish drafts and role-specific edits. A canonical version is copied deliberately when it changes.

---

## 9. Repository layout

```
.
├── .github/workflows/deploy.yml
├── .hugoversion                  # pinned version, read by CI
├── hugo.toml
├── Makefile                      # new, serve, build, check, verify
├── archetypes/
│   ├── posts.md
│   └── projects.md
├── assets/
│   ├── css/{tokens,base,layout,prose,code-dark,code-light}.css
│   ├── fonts/*.woff2
│   └── og/base.png
├── content/
│   ├── posts/<slug>/index.md + images
│   ├── projects/<slug>/index.md + images
│   ├── resume.md
│   └── about.md
├── layouts/
│   ├── baseof.html
│   ├── home.html
│   ├── single.html
│   ├── list.html
│   ├── term.html
│   ├── 404.html
│   ├── projects/{list,single}.html
│   ├── _markup/render-image.html
│   └── _partials/{head,nav,footer,post-row,project-card,metrics,validate,og-image,theme-init}.html
├── static/                       # CNAME, favicons, robots.txt, resume PDF
├── tools/verify/main.go          # build verification
└── go.mod
```

Roughly 8 layouts and 9 partials, each with one job. No theme directory, no submodule — third-party Hugo themes are the main source of long-term Hugo maintenance pain, and the whole point of this stack choice is low churn.

Hugo version is pinned in `.hugoversion` and consumed by both the Makefile and CI. It must be ≥ 0.146 for the flattened `layouts/` convention used above; the exact version is resolved and recorded in Phase 0.

---

## 10. Verification and reliability

The reliability argument is structural: the output is static HTML and CSS, so there is no runtime to fail. What remains is protecting against shipping *broken* static output.

| Gate | Mechanism | Catches |
|---|---|---|
| Frontmatter | `validate.html` → `errorf` | Missing title, malformed date, bad `kind` |
| Build warnings | Warnings escalated to failures in CI | Silently degraded builds |
| Route coverage | `tools/verify` (Go, stdlib only) | A content file that produced no page |
| Social cards | `tools/verify` | A post with no generated OG image |
| Links | `lychee` over `public/` | Dead internal and external links |
| Deploy path | Actions is the only writer to Pages | Hand-published drift |

`tools/verify` is a small Go program — stdlib only, no modules to update — that walks `content/` and `public/` and asserts every post and project produced an HTML page and an OG image, that no page has an empty `<title>`, and that internal links resolve to emitted files. It gives real Go in the repo and is fast enough to run on every build.

Unit tests are not warranted for a content site; there is no application logic. `tools/verify` has table-driven tests for its own path-mapping helpers, and that is the extent of the test suite. Lighthouse is a one-time manual spot check, not a CI gate.

---

## 11. Build order

Each phase ends with something deployed and working.

**Phase 0 — Prove the pipeline.**
Clear the repo of the old scaffold, preserving git history, the `LICENSE`, and the domain (`CNAME` moves to `static/CNAME` so Hugo emits it). Scaffold Hugo, pin the version, ship a single unstyled page, and confirm Actions → Pages → `moknshaik.com` works end to end.
*Done when:* a trivial page is live on the custom domain over HTTPS, deployed by CI.
*Rationale:* nothing gets built on an unproven deploy path.

**Phase 1 — Design system.**
Tokens, both themes, type scale, self-hosted fonts, dual Chroma stylesheets, prose styles, `baseof`, nav, footer, theme toggle.
*Done when:* a hand-written sample page renders correctly in both themes with no external requests.

**Phase 2 — Writing.**
Posts section, archetype, writing index with the `$ ls` treatment, post layout, image render hook, RSS, sitemap, OG card generation, `head` meta and JSON-LD.
*Done when:* a post with two pictures publishes end to end and its shared link previews correctly.
**This is the point at which publishing can begin; later phases do not block writing.**

**Phase 3 — Projects.**
Projects section, archetype, case-study layout, metrics rendering, index. Port the four case studies from §8.
*Done when:* all four are live and each states its problem, constraints, design, and outcome.

**Phase 4 — Home, résumé, about.**
Home page composition, `/resume/` with PDF download and links into case studies, `/about/`, `Person` JSON-LD. Migrate the two existing drafts: the README essay ("The Real Leverage Problem") becomes the first post; the Python production-practices piece is migrated as `index.html` to exercise the HTML content path.
*Done when:* a résumé claim can be followed to its case study in one click.

**Phase 5 — Guardrails.**
`tools/verify`, `lychee`, warnings-as-failures wired into CI. Favicons, `robots.txt`, `site.webmanifest`, 404 page.
*Done when:* CI fails on a deliberately broken link and a deliberately missing OG image.

### Steady-state authoring

```
make new "The Real Leverage Problem"   # archetype-scaffolded bundle
# write; drag images into the post folder
make serve                             # localhost, live reload, drafts visible
git push                               # CI builds, deploys, live in ~60s
```

---

## 12. Risks

| Risk | Mitigation |
|---|---|
| Go templates are unpleasant to edit | Only ~17 small files, no theme; layout churn is rare after Phase 1 |
| No typed frontmatter | `validate.html` fails the build with the offending file path |
| Dark technical aesthetic reads as generic dev-blog | Light mode is fully supported; mono restricted to metadata; two accents maximum |
| Hugo upgrades change template lookup (as 0.146 did) | Version pinned in `.hugoversion`; upgrades are a deliberate, tested change |
| Building the site displaces writing | Phase 2 makes the site publishable; Phases 3–5 are non-blocking |
| Manual résumé PDF copy goes stale | `/resume/` HTML is the canonical version; the PDF is explicitly a snapshot |

---

## 13. On-ramps

Nothing here is a dead end. Named future paths, none built now:

- **Search** → Pagefind, standalone binary, static index.
- **Interactive demos or charts** → Hugo has esbuild embedded via `js.Build`; add one script, no bundler.
- **Comments** → Giscus, or keep using `discuss` links.
- **Newsletter** → a Buttondown form in the footer partial.
- **Analytics** → GoatCounter, one script tag, no cookie banner.
- **Anything dynamic** → the content is portable markdown; move off Pages to a Go service and keep `content/` unchanged.
- **Notes / digital garden** → a third Hugo section reusing the same pipeline.
- **Own the generator** → replace Hugo with a Go program using goldmark, chroma, `html/template`, and `x/image`. The content directory ports unchanged, and the generator becomes its own case study.
