# Unified Blogs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `software-development/subagent-driven-development` (recommended) or `executing-plans` to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Store every article in one Hugo `blogs` section, list all articles at `/blogs/`, and list only `case-study`-tagged articles at `/case-studies/`.

**Architecture:** All six page bundles move to `content/blogs/` and share the blogs list and single templates. A standalone `content/case-studies.md` page targets `layouts/case-studies.html`, which filters the blogs collection by the exact `case-study` tag. Build verification reads the generated JSON article index and confirms both list pages contain exactly the expected canonical `/blogs/<slug>/` links.

**Tech Stack:** Hugo 0.164.0 extended, Go standard library tests, Go HTML templates, Markdown/HTML page bundles, Make.

## Global Constraints

- The exact discriminator is the tag `case-study`.
- `/blogs/` contains all six current articles.
- `/case-studies/` contains exactly the four tagged case studies.
- Every article, including a case study, has a canonical `/blogs/<slug>/` URL.
- Keep existing article prose, slugs, résumé claims, case-study metrics, and image bundles unchanged.
- Preserve normal Hugo tag archives.
- Do not add dependencies or modify deployment infrastructure.
- Do not add redirects or aliases without separate approval.
- Do not commit, push, or deploy unless the user explicitly requests it.
- Do not edit `public/` or `resources/_gen/` manually; regenerate them with Hugo.

---

## File Structure

- `content/blogs/_index.md` — metadata for the complete blogs index.
- `content/blogs/<slug>/index.{md,html}` — all six article bundles.
- `content/case-studies.md` — standalone filtered-index page metadata; contains no article body.
- `archetypes/blogs.md` — common blog frontmatter plus optional case-study fields.
- `layouts/blogs/list.html` — all-blogs list.
- `layouts/blogs/single.html` — shared article renderer with conditional case-study metadata.
- `layouts/case-studies.html` — filtered case-study index.
- `layouts/_partials/case-study-card.html` — case-study summary card.
- `layouts/_partials/{nav,head,og-image,validate}.html` — unified section wiring and validation.
- `layouts/{home.html,home.json.json,rss.xml,404.html}` — unified site-level consumers.
- `tools/verify/main.go` and `main_test.go` — generated-list consistency checks.

### Task 1: Define generated-list verification behavior

**Files:**
- Modify: `tools/verify/main_test.go`
- Modify: `tools/verify/main.go`

**Interfaces:**
- Consumes: existing `page` values from `public/index.json`.
- Produces: `page.CaseStudy bool` and `listingProblems(name string, body []byte, expected []page) []string` for final build verification.

- [ ] **Step 1: Write failing tests for exact list membership**

Add a table-driven test to `tools/verify/main_test.go`:

```go
func TestListingProblems(t *testing.T) {
	pages := []page{
		{URL: "/blogs/ordinary/"},
		{URL: "/blogs/case-study/", CaseStudy: true},
	}

	t.Run("accepts exact membership", func(t *testing.T) {
		body := []byte(`<a href="/blogs/case-study/">Case study</a>`)
		if problems := listingProblems("case studies", body, pages[1:]); len(problems) != 0 {
			t.Fatalf("listingProblems() = %v, want no problems", problems)
		}
	})

	t.Run("reports missing and unexpected articles", func(t *testing.T) {
		body := []byte(`<a href="/blogs/ordinary/">Ordinary</a>`)
		problems := strings.Join(listingProblems("case studies", body, pages[1:]), "\n")
		for _, want := range []string{"missing /blogs/case-study/", "unexpected /blogs/ordinary/"} {
			if !strings.Contains(problems, want) {
				t.Errorf("listingProblems() = %q, want %q", problems, want)
			}
		}
	})
}
```

Also change the page-path case to the unified canonical path:

```go
{"blog", "/blogs/api/", "public/blogs/api/index.html"},
```

- [ ] **Step 2: Run the focused test and confirm RED**

Run:

```bash
go test ./tools/verify -run 'TestListingProblems|TestPagePathFor' -count=1
```

Expected: compilation fails because `page.CaseStudy` and `listingProblems` do not exist.

- [ ] **Step 3: Add the JSON field and exact-membership helper**

In `tools/verify/main.go`, extend `page`:

```go
type page struct {
	URL       string `json:"url"`
	Title     string `json:"title"`
	OG        string `json:"og"`
	Section   string `json:"section"`
	CaseStudy bool   `json:"caseStudy"`
}
```

Add a package-level regular expression and helper:

```go
var blogLinkRe = regexp.MustCompile(`href=["'](/blogs/[^"'#?]+/)["']`)

func listingProblems(name string, body []byte, expected []page) []string {
	want := make(map[string]bool, len(expected))
	for _, p := range expected {
		want[p.URL] = true
	}

	got := make(map[string]bool)
	for _, match := range blogLinkRe.FindAllSubmatch(body, -1) {
		got[string(match[1])] = true
	}

	var problems []string
	for u := range want {
		if !got[u] {
			problems = append(problems, fmt.Sprintf("%s: missing %s", name, u))
		}
	}
	for u := range got {
		if !want[u] {
			problems = append(problems, fmt.Sprintf("%s: unexpected %s", name, u))
		}
	}
	return problems
}
```

Use `sort.Strings(problems)` before returning so failures are deterministic; add the `sort` import.

- [ ] **Step 4: Run focused and full Go tests**

Run:

```bash
go test ./tools/verify -run 'TestListingProblems|TestPagePathFor' -count=1
go test ./tools/... -count=1
```

Expected: both commands pass.

- [ ] **Step 5: Review checkpoint**

Inspect `git diff -- tools/verify/main.go tools/verify/main_test.go`. Do not commit.

### Task 2: Migrate all article bundles into the blogs content model

**Files:**
- Create: `archetypes/blogs.md`
- Create: `content/blogs/_index.md`
- Create: `content/case-studies.md`
- Move: `content/writing/*` to `content/blogs/`
- Move: `content/case-studies/*` to `content/blogs/`
- Remove: `archetypes/writing.md`, `archetypes/case-studies.md`

**Interfaces:**
- Consumes: the six existing article bundles and the exact `case-study` discriminator.
- Produces: one `blogs` section containing six pages, four with `tags: ["case-study"]` plus any existing tags.

- [ ] **Step 1: Add a failing unified-section test**

Add `TestArticleIndexProblemsRejectsNonBlogSections` to `tools/verify/main_test.go`:

```go
func TestArticleIndexProblemsRejectsNonBlogSections(t *testing.T) {
	pages := []page{
		{URL: "/blogs/current/", Section: "blogs"},
		{URL: "/writing/legacy/", Section: "writing"},
	}
	problems := strings.Join(articleIndexProblems(pages), "\n")
	if !strings.Contains(problems, `article section is "writing", want blogs`) {
		t.Fatalf("articleIndexProblems() = %q, want stale section problem", problems)
	}
}
```

- [ ] **Step 2: Run the focused test and confirm RED**

Run:

```bash
go test ./tools/verify -run TestArticleIndexProblemsRejectsNonBlogSections -count=1
```

Expected: compilation fails because `articleIndexProblems` does not exist.

- [ ] **Step 3: Implement the unified-section check and confirm GREEN**

Add this helper to `tools/verify/main.go`:

```go
func articleIndexProblems(pages []page) []string {
	var problems []string
	for _, p := range pages {
		if p.Section != "blogs" {
			problems = append(problems, fmt.Sprintf("%s: article section is %q, want blogs", p.URL, p.Section))
		}
	}
	return problems
}
```

Run:

```bash
go test ./tools/verify -run TestArticleIndexProblemsRejectsNonBlogSections -count=1
```

Expected: PASS.

- [ ] **Step 4: Move the bundles without changing prose or slugs**

Run these exact moves from the repository root:

```bash
mkdir -p content/blogs
git mv content/writing/_index.md content/blogs/_index.md
git mv content/case-studies/_index.md content/case-studies.md
git mv content/writing/essentials-of-production-python content/blogs/
git mv content/writing/real-leverage-problem-ai-in-software-development content/blogs/
git mv content/case-studies/document-conversion-api content/blogs/
git mv content/case-studies/llm-agent-runtime content/blogs/
git mv content/case-studies/managed-analytics-platform content/blogs/
git mv content/case-studies/promotion-approval-workflow content/blogs/
```

The source directories disappear from Git once their tracked contents have moved; do not run a separate recursive delete.

- [ ] **Step 5: Add the exact case-study tag to four bundles**

For each migrated case study, insert this field after `date` without changing any existing metadata or body:

```yaml
tags: ["case-study"]
```

Files:

- `content/blogs/document-conversion-api/index.md`
- `content/blogs/llm-agent-runtime/index.md`
- `content/blogs/managed-analytics-platform/index.md`
- `content/blogs/promotion-approval-workflow/index.md`

- [ ] **Step 6: Update the blogs index and standalone filtered page**

Replace the moved `content/blogs/_index.md` frontmatter with:

```yaml
---
title: Blogs
description: Writing on backend systems, APIs, and where AI actually pays off in software work.
---
```

Replace the moved `content/case-studies.md` frontmatter with:

```yaml
---
title: Case Studies
description: Systems I built, what constrained them, and what changed as a result.
layout: case-studies
---
```

- [ ] **Step 7: Replace both archetypes with one blogs archetype**

Create `archetypes/blogs.md`:

```yaml
---
title: "{{ replace .File.ContentBaseName "-" " " | title }}"
description: ""
date: {{ .Date }}
tags: []
period:
  start: ""
  end: ""
role: ""
org: ""
stack: []
metrics: []
links: []
featured: false
weight: 100
draft: true
---

Opening paragraph — state the claim before the setup.
```

Delete the obsolete writing and case-study archetypes after the new file exists.

- [ ] **Step 8: Verify the source content model**

Run this read-only assertion:

```bash
python3 - <<'PY'
from pathlib import Path

pages = [p for p in Path('content/blogs').glob('*/index.*') if p.suffix in {'.md', '.html'}]
case_studies = [p for p in pages if 'case-study' in p.read_text()]
assert len(pages) == 6, len(pages)
assert len(case_studies) == 4, len(case_studies)
print('source model: 6 blogs, 4 case studies')
PY
```

Expected output: `source model: 6 blogs, 4 case studies`.

- [ ] **Step 9: Review checkpoint**

Inspect `git status --short` and `git diff -- content archetypes`. Do not commit.

### Task 3: Build the unified blogs and filtered case-studies renderers

**Files:**
- Create: `layouts/blogs/list.html`
- Create: `layouts/blogs/single.html`
- Create: `layouts/case-studies.html`
- Modify: `layouts/_partials/case-study-card.html`
- Remove: `layouts/writing/list.html`, `layouts/writing/single.html`
- Remove: `layouts/case-studies/list.html`, `layouts/case-studies/single.html`

**Interfaces:**
- Consumes: `site.RegularPages` from section `blogs`, `.Params.tags`, and existing `metrics.html`, `post-row.html`, and `case-study-card.html` partials.
- Produces: `/blogs/`, `/blogs/<slug>/`, and `/case-studies/` HTML.

- [ ] **Step 1: Create the all-blogs list template**

Create `layouts/blogs/list.html`:

```go-html-template
{{ define "main" }}
  <header class="section-head">
    <h1><span class="prompt">$</span> ls blogs/</h1>
    <p class="section-head__desc">{{ .Params.description }}</p>
  </header>

  {{ range .RegularPagesRecursive.GroupByDate "2006" }}
    <section class="year">
      <h2 class="year__label">{{ .Key }}</h2>
      <ul class="rows">
        {{ range .Pages }}{{ partial "post-row.html" . }}{{ end }}
      </ul>
    </section>
  {{ end }}
{{ end }}
```

- [ ] **Step 2: Create the shared article template**

Create `layouts/blogs/single.html`. Define classification once:

```go-html-template
{{ $tags := .Params.tags | default (slice) }}
{{ $isCaseStudy := in $tags "case-study" }}
```

Render the existing date, reading time, updated date, and tags for every article. Inside the header, conditionally render the existing case-study metadata:

```go-html-template
{{ if $isCaseStudy }}
  <p class="post__meta">
    {{ .Params.role }}{{ with .Params.org }} · {{ . }}{{ end }}
    <span class="sep">·</span>{{ .Params.period.start }}–{{ .Params.period.end | default "present" }}
  </p>
  <p class="post__desc">{{ .Params.description }}</p>
  {{ partial "metrics.html" . }}
  <ul class="stack">{{ range .Params.stack }}<li>{{ . }}</li>{{ end }}</ul>
  {{ with .Params.links }}
    <p class="post__links">{{ range . }}<a href="{{ .href }}">{{ .label }} →</a> {{ end }}</p>
  {{ end }}
{{ end }}
```

Keep `.Content` and the existing optional `discuss` block for all blogs. Do not hide `#case-study` from the tag list.

- [ ] **Step 3: Create the filtered page template using Hugo's page-layout lookup**

Hugo 0.164 selects the root custom layout `layouts/case-studies.html` for this page's `layout: case-studies` value. Create it with:

```go-html-template
{{ define "main" }}
  <header class="section-head">
    <h1><span class="prompt">$</span> ls case-studies/</h1>
    <p class="section-head__desc">{{ .Params.description }}</p>
  </header>
  {{ $blogs := where site.RegularPages "Section" "blogs" }}
  {{ $caseStudies := where $blogs "Params.tags" "intersect" (slice "case-study") }}
  <div class="cards">
    {{ range $caseStudies.ByWeight }}{{ partial "case-study-card.html" . }}{{ end }}
  </div>
{{ end }}
```

- [ ] **Step 4: Remove obsolete section templates**

Delete the writing and section-based case-study templates only after all three new templates exist.

- [ ] **Step 5: Build and inspect the three routes**

Run:

```bash
make build
python3 - <<'PY'
from pathlib import Path

blogs = Path('public/blogs/index.html').read_text()
case_studies = Path('public/case-studies/index.html').read_text()
assert Path('public/blogs/document-conversion-api/index.html').is_file()
assert '/blogs/essentials-of-production-python/' in blogs
assert '/blogs/document-conversion-api/' in blogs
assert '/blogs/document-conversion-api/' in case_studies
assert '/blogs/essentials-of-production-python/' not in case_studies
print('routes and filtered membership are correct')
PY
```

Expected: every command exits 0. If Hugo does not select `layouts/case-studies.html`, inspect `hugo --printPathWarnings --printUnusedTemplates` before changing the approved content model.

- [ ] **Step 6: Review checkpoint**

Inspect `git diff -- layouts content`. Do not commit.

### Task 4: Rewire every site-level consumer to blogs

**Files:**
- Modify: `Makefile`
- Modify: `content/resume.md`
- Modify: `layouts/home.html`
- Modify: `layouts/home.json.json`
- Modify: `layouts/rss.xml`
- Modify: `layouts/404.html`
- Modify: `layouts/_partials/nav.html`
- Modify: `layouts/_partials/head.html`
- Modify: `layouts/_partials/og-image.html`
- Modify: `layouts/_partials/validate.html`

**Interfaces:**
- Consumes: one `blogs` section and the exact `case-study` tag.
- Produces: unified navigation, homepage, metadata, feed, validation, JSON index, résumé links, and `make blog` scaffolding.

- [ ] **Step 1: Update homepage contract tests and verify RED**

In `tools/verify/main_test.go`, change valid homepage fixtures to include:

```html
<nav><a href="/blogs/">blogs</a><a href="/case-studies/">case-studies</a></nav>
<h2>$ ls blogs/</h2>
```

Change expected problems to `blogs list missing`, `blogs navigation missing`, and `case-studies navigation missing`.

Run:

```bash
go test ./tools/verify -run TestHomepageProblems -count=1
```

Expected: FAIL because `homepageProblems` still requires writing.

- [ ] **Step 2: Update homepage and navigation**

In `layouts/home.html`, query section `blogs` and change the heading to `$ ls blogs/`.

In `layouts/_partials/nav.html`, use exactly:

```go-html-template
{{ range slice (slice "blogs" "/blogs/") (slice "case-studies" "/case-studies/") (slice "résumé" "/resume/") (slice "about" "/about/") }}
```

Update `layouts/404.html` from `/writing/` to `/blogs/`.

Update `homepageProblems` requirements to exact anchors for `/blogs/` and `/case-studies/` and the `$ ls blogs/` heading.

- [ ] **Step 3: Update feed, JSON, social cards, and structured data**

Use section `blogs` in:

- `layouts/rss.xml`
- `layouts/home.json.json`
- `layouts/_partials/og-image.html`
- `layouts/_partials/head.html`

In `layouts/home.json.json`, emit the classification field:

```go-html-template
"caseStudy" (in (.Params.tags | default (slice)) "case-study")
```

The complete appended dictionary remains `url`, `title`, `og`, `section`, and `caseStudy`.

- [ ] **Step 4: Make validation conditional on the tag**

In `layouts/_partials/validate.html`, validate title, description, and date for every page in section `blogs`. Then use:

```go-html-template
{{- $tags := .Params.tags | default (slice) -}}
{{- if in $tags "case-study" -}}
  {{- range slice "role" "period" "stack" -}}
    {{- if not (index $.Params .) -}}
      {{- errorf "frontmatter: missing '%s' in %s" . $f -}}
    {{- end -}}
  {{- end -}}
{{- end -}}
```

- [ ] **Step 5: Update résumé links and publishing command**

Change all four résumé links to `/blogs/<existing-slug>/`.

Replace the separate `new` and `case-study` publishing targets with one target:

```make
.PHONY: help serve build blog check clean fonts chroma ogbase

blog:
	@test -n "$(S)" || (echo 'usage: make blog S=my-blog-slug'; exit 1)
	hugo new content "blogs/$(S)/index.md"
	@echo "edit the title in content/blogs/$(S)/index.md"
```

The help text must document `make blog S=my-slug` and `content/blogs/my-slug/`.

- [ ] **Step 6: Run focused tests and build**

Run:

```bash
go test ./tools/verify -run 'TestHomepageProblems|TestArticleIndexProblems' -count=1
make build
```

Expected: both commands pass.

- [ ] **Step 7: Review checkpoint**

Inspect `git diff -- Makefile content/resume.md layouts tools/verify`. Do not commit.

### Task 5: Integrate generated-site verification and finish migration

**Files:**
- Modify: `tools/verify/main.go`
- Modify: `tools/verify/main_test.go`
- Regenerate: `public/**`, `resources/_gen/**`

**Interfaces:**
- Consumes: `public/index.json`, `public/blogs/index.html`, and `public/case-studies/index.html`.
- Produces: a failing verification result for wrong list membership or non-blogs article sections.

- [ ] **Step 1: Add list-page integration to `run`**

After decoding `idx`, reject stale article sections:

```go
problems = append(problems, articleIndexProblems(idx.Pages)...)
```

Read `blogs/index.html` and compare it with all `idx.Pages` through `listingProblems("blogs index", body, idx.Pages)`.

Build the case-study subset:

```go
var caseStudies []page
for _, p := range idx.Pages {
	if p.CaseStudy {
		caseStudies = append(caseStudies, p)
	}
}
```

Read `case-studies/index.html` and call `listingProblems("case-studies index", body, caseStudies)`.

Replace the old writing-index read and deprecated writing-classification checks. Keep title, page-file, social-card, CNAME, and required-output checks intact.

- [ ] **Step 2: Run all Go tests**

Run:

```bash
make test
```

Expected: all Go tests pass.

- [ ] **Step 3: Regenerate from a clean destination with approval**

Because stale generated `/writing/` and section-based case-study detail files survive an ordinary Hugo build, obtain approval for the repository's existing destructive cleanup target, then run:

```bash
make clean
make build
```

Expected: only current Hugo output remains. Do not manually delete individual generated files.

- [ ] **Step 4: Run full checks**

Run:

```bash
make check
```

Expected: Hugo build, Go tests, generated-site verifier, and offline link checker all exit 0.

- [ ] **Step 5: Run the required standalone production build**

Run after the checked source is unchanged:

```bash
make build
```

Expected: exit 0.

- [ ] **Step 6: Verify exact published membership**

Run:

```bash
python3 - <<'PY'
import json
from pathlib import Path

pages = json.loads(Path('public/index.json').read_text())['pages']
assert len(pages) == 6, len(pages)
assert {p['section'] for p in pages} == {'blogs'}
assert sum(bool(p['caseStudy']) for p in pages) == 4
assert all(p['url'].startswith('/blogs/') for p in pages)
print('6 blogs; 4 case studies; all canonical URLs under /blogs/')
PY
```

Expected output: `6 blogs; 4 case studies; all canonical URLs under /blogs/`.

- [ ] **Step 7: Verify obsolete active taxonomy paths are gone**

Use `search_files(target='files')` to confirm there are no files under `content/writing/`, `layouts/writing/`, `content/case-studies/`, or `layouts/case-studies/`. Then use `search_files(target='content')`, excluding historical docs and generated output, for the regex `Section" "(writing|case-studies)"|/writing/|/case-studies/[a-z-]+/`.

Expected: no active-source matches. `content/case-studies.md`, `layouts/case-studies.html`, and the `/case-studies/` index URL are intentional and therefore not rejected.

- [ ] **Step 8: Final diff and status review**

Run:

```bash
git status --short
git diff --check
git diff --stat
```

Confirm that only the approved migration, its spec/plan, tests, and regenerated artifacts changed. Do not commit.
