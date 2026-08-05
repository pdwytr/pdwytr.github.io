package main

import (
	"strings"
	"testing"
)

func TestHomepageProblemsRequiresCompleteWritingList(t *testing.T) {
	html := []byte(`<nav><a href="/projects/">projects</a></nav><main><h2>$ ls projects/ --featured</h2></main>`)
	problems := strings.Join(homepageProblems(html), "\n")

	for _, want := range []string{
		"writing list missing",
		"featured projects still present",
		"case-studies navigation missing",
		"projects navigation still present",
		"GitHub social link missing",
		"X (Twitter) social link missing",
		"LinkedIn social link missing",
		"theme switch missing",
	} {
		if !strings.Contains(problems, want) {
			t.Errorf("homepageProblems() = %q, want problem containing %q", problems, want)
		}
	}
}

func TestHomepageProblemsAcceptsFocusedHomepageControls(t *testing.T) {
	html := []byte(`<main>
		<a href="https://github.com/pdwytr" aria-label="GitHub"></a>
		<a href="https://x.com/pdwytrfa" aria-label="X (Twitter)"></a>
		<a href="https://www.linkedin.com/in/pdwytr/" aria-label="LinkedIn"></a>
		<nav><a href="/projects/">case-studies</a></nav>
		<h2>$ ls writing/</h2>
		<button role="switch" aria-checked="true" aria-label="Switch to light theme"></button>
	</main>`)

	if problems := homepageProblems(html); len(problems) != 0 {
		t.Fatalf("homepageProblems() = %v, want no problems", problems)
	}
}

func TestHomepageProblemsRejectsSocialURLsOutsideAccessibleAnchors(t *testing.T) {
	html := []byte(`<main>
		<script type="application/ld+json">{"sameAs":["https://github.com/pdwytr","https://x.com/pdwytrfa","https://www.linkedin.com/in/pdwytr/"]}</script>
		<nav><a href="/projects/">case-studies</a></nav>
		<h2>$ ls writing/</h2>
		<button role="switch" aria-checked="true" aria-label="Dark theme"></button>
	</main>`)
	problems := strings.Join(homepageProblems(html), "\n")

	for _, want := range []string{"GitHub social link missing", "X (Twitter) social link missing", "LinkedIn social link missing"} {
		if !strings.Contains(problems, want) {
			t.Errorf("homepageProblems() = %q, want problem containing %q", problems, want)
		}
	}
}

func TestDeprecatedClassificationProblems(t *testing.T) {
	if problems := deprecatedClassificationProblems([]byte(`<span>[essay]</span>`)); len(problems) != 1 {
		t.Fatalf("deprecatedClassificationProblems() = %v, want one problem", problems)
	}
	if problems := deprecatedClassificationProblems([]byte(`<article>Writing without a classification.</article>`)); len(problems) != 0 {
		t.Fatalf("deprecatedClassificationProblems() = %v, want no problems", problems)
	}
}

func TestPagePathFor(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"post", "/writing/my-post/", "public/writing/my-post/index.html"},
		{"project", "/projects/api/", "public/projects/api/index.html"},
		{"no trailing slash", "/writing/my-post", "public/writing/my-post/index.html"},
		{"root", "/", "public/index.html"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := pagePathFor("public", tc.url); got != tc.want {
				t.Errorf("pagePathFor(%q) = %q, want %q", tc.url, got, tc.want)
			}
		})
	}
}

func TestAssetPathFor(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"absolute", "https://moknshaik.com/og/my-post.png", "public/og/my-post.png"},
		{"root relative", "/og-default.png", "public/og-default.png"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := assetPathFor("public", tc.raw)
			if err != nil {
				t.Fatalf("assetPathFor(%q): %v", tc.raw, err)
			}
			if got != tc.want {
				t.Errorf("assetPathFor(%q) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}
