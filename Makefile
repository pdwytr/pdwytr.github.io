HUGO_VERSION := $(shell cat .hugoversion)
BASEURL := https://moknshaik.com/

.PHONY: help serve build new project check clean

help:
	@echo "make serve              live preview with drafts  -> http://localhost:1313"
	@echo "make build              production build          -> public/"
	@echo "make new S=my-slug      scaffold a post           -> content/writing/my-slug/"
	@echo "make project S=my-slug  scaffold a case study      -> content/projects/my-slug/"
	@echo "make check              build + verify + link check"
	@echo "make clean              remove build output"
	@echo ""
	@echo "pinned hugo: $(HUGO_VERSION)"

serve:
	hugo server -D --disableFastRender

build:
	hugo --gc --minify --baseURL "$(BASEURL)"

new:
	@test -n "$(S)" || (echo 'usage: make new S=my-post-slug'; exit 1)
	hugo new content "writing/$(S)/index.md"
	@echo "edit the title in content/writing/$(S)/index.md"

project:
	@test -n "$(S)" || (echo 'usage: make project S=my-project-slug'; exit 1)
	hugo new content "projects/$(S)/index.md"

clean:
	rm -rf public resources/_gen
