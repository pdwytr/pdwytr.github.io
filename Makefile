HUGO_VERSION := $(shell cat .hugoversion)
BASEURL := https://moknshaik.com/

.PHONY: help serve build new project check clean fonts chroma ogbase

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

FONTSOURCE := https://cdn.jsdelivr.net/npm
JBM_TTF_TAG := v2.304

.PHONY: fonts
fonts:
	mkdir -p assets/fonts
	curl -sSLf -o assets/fonts/inter-latin-wght-normal.woff2 \
	  "$(FONTSOURCE)/@fontsource-variable/inter@5.2.5/files/inter-latin-wght-normal.woff2"
	curl -sSLf -o assets/fonts/jetbrains-mono-latin-400-normal.woff2 \
	  "$(FONTSOURCE)/@fontsource/jetbrains-mono@5.2.5/files/jetbrains-mono-latin-400-normal.woff2"
	curl -sSLf -o assets/fonts/jetbrains-mono-latin-700-normal.woff2 \
	  "$(FONTSOURCE)/@fontsource/jetbrains-mono@5.2.5/files/jetbrains-mono-latin-700-normal.woff2"
	curl -sSLf -o /tmp/jbm.zip \
	  "https://github.com/JetBrains/JetBrainsMono/releases/download/$(JBM_TTF_TAG)/JetBrainsMono-2.304.zip"
	unzip -o -j /tmp/jbm.zip "fonts/ttf/JetBrainsMono-Bold.ttf" -d assets/fonts/
	@ls -la assets/fonts/

.PHONY: chroma
chroma:
	printf ':root, :root[data-theme="dark"] {\n' > assets/css/code-dark.css
	hugo gen chromastyles --style=github-dark >> assets/css/code-dark.css
	printf '}\n' >> assets/css/code-dark.css
	printf ':root[data-theme="light"] {\n' > assets/css/code-light.css
	hugo gen chromastyles --style=github >> assets/css/code-light.css
	printf '}\n' >> assets/css/code-light.css
	@wc -l assets/css/code-dark.css assets/css/code-light.css

.PHONY: ogbase
ogbase:
	mkdir -p assets/og
	go run ./tools/ogbase -out assets/og/base.png
	go run ./tools/ogbase -out static/og-default.png

.PHONY: verify test
test:
	go test ./tools/...

verify: build
	go run ./tools/verify -public public

check: build
	go test ./tools/...
	go run ./tools/verify -public public
	lychee --no-progress --offline --root-dir "$(PWD)/public" public
