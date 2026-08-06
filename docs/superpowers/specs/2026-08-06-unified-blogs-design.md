# Unified Blogs and Case Studies Design

**Date:** 2026-08-06
**Status:** Approved

## Objective

Replace the separate writing and case-study content sections with one blog content model. Every published article appears in the blogs index. Articles tagged `case-study` also appear in a dedicated case-studies index.

The resulting public structure is:

- `/blogs/` — every published article, reverse chronological
- `/blogs/<slug>/` — the canonical URL for every article, including case studies
- `/case-studies/` — only articles whose `tags` contain `case-study`
- `/tags/<tag>/` — Hugo's normal tag archives

Navigation contains `blogs` and `case-studies`. The old `/writing/` section is removed. Case studies no longer have canonical detail pages below `/case-studies/`.

## Content Model

All article bundles live below `content/blogs/` and use the blogs archetype.

Every article supports:

- `title`
- `description`
- `date`
- `lastmod` (optional)
- `tags` (optional for ordinary blogs)
- `draft`

A case study is an ordinary blog whose `tags` list contains the exact tag `case-study`. Case studies additionally require:

- `period`
- `role`
- `stack`

They may also define:

- `org`
- `metrics`
- `links`
- `weight`
- `featured`

Tag membership is the only classification mechanism. There is no separate case-study section, boolean, or duplicate content record.

## Rendering

### Blogs index

`/blogs/` lists every regular page in the `blogs` section. Case studies are not excluded or visually duplicated; each article appears once in the complete list.

### Case-studies index

`/case-studies/` is a standalone Hugo page defined by `content/case-studies.md`, not a second article section. Its dedicated template selects regular pages from the `blogs` section whose tags contain `case-study`. It retains the current case-study card presentation, including role, period, metrics, and stack.

### Article pages

All detail pages use the blogs single template. The template checks whether `case-study` is present in the page's tags:

- Ordinary blogs render date, reading time, tags, and prose.
- Case studies additionally render role, organization, period, metrics, stack, and external links.

The case-study tag remains visible with the article's other tags.

## Data Flow

1. Hugo discovers all article bundles in `content/blogs/`.
2. The blogs index and homepage query the `blogs` section.
3. The case-studies index queries the same pages and filters by the `case-study` tag.
4. RSS, the JSON page index, social-card generation, and frontmatter validation consume the unified `blogs` section.
5. Résumé case-study links point to the articles' canonical `/blogs/<slug>/` URLs.

## Migration

- Move existing `content/writing/*` bundles into `content/blogs/` without changing their slugs or article content.
- Move existing case-study bundles into `content/blogs/`.
- Add `case-study` to every migrated case study's tags.
- Replace the writing and case-study archetypes with one blogs archetype containing the common fields and optional case-study fields.
- Replace section-specific templates with blogs templates and the filtered case-studies list template.
- Rename the publishing command to scaffold content below `content/blogs/`.
- Update internal links and generated-output checks to use `/blogs/` canonical article URLs.
- Remove obsolete writing and case-study section identifiers and paths from active source files.

Historical design and implementation documents remain historical records and are not runtime source of truth.

## Commands

- Scaffold: `make blog S=my-blog-slug`
- Develop: `make serve`
- Unit tests: `make test`
- Full checks: `make check`
- Production build: `make build`

## Testing Strategy

Tests must be written or updated before production changes and observed failing for the missing unified behavior.

Automated verification must prove:

- The homepage links to `/blogs/` and `/case-studies/`.
- The generated JSON index contains only `blogs` as the article section.
- `/blogs/` includes ordinary blogs and case studies.
- `/case-studies/` includes tagged case studies and excludes ordinary blogs.
- Case-study metadata validation is applied only when the `case-study` tag is present.
- Every résumé case-study link resolves to a generated `/blogs/<slug>/` page.
- RSS and social cards still include all articles.
- No active source path or section identifier refers to the removed `writing` or separate `case-studies` content sections.

Final verification uses `make test`, `make check`, and `make build`.

## Boundaries

### Always

- Keep every article in one `blogs` section.
- Use the exact tag `case-study` as the case-study discriminator.
- Keep existing article slugs and content intact during migration.
- Preserve ordinary Hugo tag archives.
- Keep case-study metadata and card presentation.

### Ask first

- Add redirects or aliases for old public URLs.
- Change article prose, résumé claims, or case-study metrics.
- Add dependencies or modify deployment infrastructure.

### Never

- Duplicate an article between blogs and case studies.
- Introduce a second case-study boolean or classification field.
- Commit, push, or deploy without explicit user instruction.
- Modify generated output manually; regenerate it through Hugo.

## Success Criteria

- All six current articles are stored below `content/blogs/`.
- `/blogs/` lists all six articles.
- `/case-studies/` lists exactly the four migrated case studies.
- Each case study's canonical detail URL is `/blogs/<slug>/`.
- Ordinary blogs do not appear in `/case-studies/`.
- Case-study detail pages retain their specialized metadata.
- Navigation, homepage, résumé, RSS, JSON index, social cards, and validation use the unified model.
- `make test`, `make check`, and `make build` exit successfully.
