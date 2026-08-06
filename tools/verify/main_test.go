package main

import (
	"strings"
	"testing"
)

func TestHomepageProblemsRequiresCompleteBlogsList(t *testing.T) {
	html := []byte(`<nav></nav><main></main>`)
	problems := strings.Join(homepageProblems(html), "\n")

	for _, want := range []string{
		"blogs list missing",
		"blogs navigation missing",
		"case-studies navigation missing",
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
		<nav><a href="/blogs/">blogs</a><a href="/case-studies/">case-studies</a></nav>
		<h2>$ ls blogs/</h2>
		<button role="switch" aria-checked="true" aria-label="Switch to light theme"></button>
	</main>`)

	if problems := homepageProblems(html); len(problems) != 0 {
		t.Fatalf("homepageProblems() = %v, want no problems", problems)
	}
}

func TestHomepageProblemsAcceptsHugoMinifiedNavigationLinks(t *testing.T) {
	html := []byte(`<main>
		<a href="https://github.com/pdwytr" aria-label="GitHub"></a>
		<a href="https://x.com/pdwytrfa" aria-label="X (Twitter)"></a>
		<a href="https://www.linkedin.com/in/pdwytr/" aria-label="LinkedIn"></a>
		<nav><a href=/blogs/>blogs</a><a href=/case-studies/>case-studies</a></nav>
		<h2>$ ls blogs/</h2>
		<button role="switch" aria-checked="true" aria-label="Switch to light theme"></button>
	</main>`)

	if problems := homepageProblems(html); len(problems) != 0 {
		t.Fatalf("homepageProblems() = %v, want no problems", problems)
	}
}

func TestHasExactLink(t *testing.T) {
	for _, tc := range []struct {
		name string
		html string
		want bool
	}{
		{"double quoted", `<a href="/blogs/">blogs</a>`, true},
		{"single quoted", `<a href='/blogs/'>blogs</a>`, true},
		{"unquoted", `<a href=/blogs/>blogs</a>`, true},
		{"other attributes", `<a class=nav href=/blogs/ aria-current=page>blogs</a>`, true},
		{"wrong URL", `<a href=/work/>blogs</a>`, false},
		{"URL substring", `<a href=/blogs/archive/>blogs</a>`, false},
		{"wrong text", `<a href=/blogs/>all blogs</a>`, false},
		{"attribute substring", `<a data-href=/blogs/>blogs</a>`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasExactLink(tc.html, "/blogs/", "blogs"); got != tc.want {
				t.Errorf("hasExactLink(%q) = %t, want %t", tc.html, got, tc.want)
			}
		})
	}
}

func TestHomepageProblemsRejectsWrongCaseStudiesURL(t *testing.T) {
	html := []byte(`<main>
		<a href="https://github.com/pdwytr" aria-label="GitHub"></a>
		<a href="https://x.com/pdwytrfa" aria-label="X (Twitter)"></a>
		<a href="https://www.linkedin.com/in/pdwytr/" aria-label="LinkedIn"></a>
		<nav><a href="/blogs/">blogs</a><a href="/work/">case-studies</a></nav>
		<h2>$ ls blogs/</h2>
		<button role="switch" aria-checked="true" aria-label="Switch to light theme"></button>
	</main>`)

	problems := strings.Join(homepageProblems(html), "\n")
	if !strings.Contains(problems, "case-studies navigation missing") {
		t.Fatalf("homepageProblems() = %q, want case-studies navigation problem", problems)
	}
}

func TestHomepageProblemsRejectsSocialURLsOutsideAccessibleAnchors(t *testing.T) {
	html := []byte(`<main>
		<script type="application/ld+json">{"sameAs":["https://github.com/pdwytr","https://x.com/pdwytrfa","https://www.linkedin.com/in/pdwytr/"]}</script>
		<nav><a href="/blogs/">blogs</a><a href="/case-studies/">case-studies</a></nav>
		<h2>$ ls blogs/</h2>
		<button role="switch" aria-checked="true" aria-label="Dark theme"></button>
	</main>`)
	problems := strings.Join(homepageProblems(html), "\n")

	for _, want := range []string{"GitHub social link missing", "X (Twitter) social link missing", "LinkedIn social link missing"} {
		if !strings.Contains(problems, want) {
			t.Errorf("homepageProblems() = %q, want problem containing %q", problems, want)
		}
	}
}

func TestArticleIndexProblemsRejectsNonBlogSections(t *testing.T) {
	pages := []page{
		{URL: "/blogs/current/", Section: "blogs"},
		{URL: "/writing/legacy/", Section: "writing"},
	}

	problems := strings.Join(articleIndexProblems(pages), "\n")
	if want := `article section is "writing", want blogs`; !strings.Contains(problems, want) {
		t.Fatalf("articleIndexProblems() = %q, want problem containing %q", problems, want)
	}
}

func TestListingProblems(t *testing.T) {
	pages := []page{
		{URL: "/blogs/ordinary/"},
		{URL: "/blogs/case-study/", CaseStudy: true},
	}

	t.Run("exact membership", func(t *testing.T) {
		html := []byte(`<a href="/blogs/case-study/">case study</a>`)
		if problems := listingProblems("case studies", html, pages[1:]); len(problems) != 0 {
			t.Fatalf("listingProblems() = %v, want no problems", problems)
		}
	})

	t.Run("minified unquoted membership", func(t *testing.T) {
		html := []byte(`<a href=/blogs/case-study/>case study</a>`)
		if problems := listingProblems("case studies", html, pages[1:]); len(problems) != 0 {
			t.Fatalf("listingProblems() = %v, want no problems", problems)
		}
	})

	t.Run("missing and unexpected", func(t *testing.T) {
		html := []byte(`<a href="/blogs/ordinary/">ordinary</a>`)
		problems := strings.Join(listingProblems("case studies", html, pages[1:]), "\n")
		for _, want := range []string{"missing /blogs/case-study/", "unexpected /blogs/ordinary/"} {
			if !strings.Contains(problems, want) {
				t.Errorf("listingProblems() = %q, want problem containing %q", problems, want)
			}
		}
	})
}

func TestPagePathFor(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"post", "/archive/my-post/", "public/archive/my-post/index.html"},
		{"blog", "/blogs/api/", "public/blogs/api/index.html"},
		{"no trailing slash", "/archive/my-post", "public/archive/my-post/index.html"},
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
