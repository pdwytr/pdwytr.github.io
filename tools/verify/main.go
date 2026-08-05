// Command verify checks a built Hugo site for structural defects that a
// successful build can still hide: a page listed in the index with no HTML on
// disk, an empty <title>, a missing social card, or a dropped CNAME.
//
// It reads public/index.json (emitted by layouts/home.json.json) rather than
// parsing content frontmatter, so Hugo remains the single source of truth for
// what should exist. Link checking is lychee's job, not this program's.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type page struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	OG      string `json:"og"`
	Section string `json:"section"`
}

type index struct {
	Count int    `json:"count"`
	Pages []page `json:"pages"`
}

var titleRe = regexp.MustCompile(`(?is)<title>(.*?)</title>`)

func main() {
	root := flag.String("public", "public", "path to the built site")
	flag.Parse()

	problems, err := run(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify: %v\n", err)
		os.Exit(2)
	}
	for _, p := range problems {
		fmt.Fprintf(os.Stderr, "verify: %s\n", p)
	}
	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "verify: %d problem(s)\n", len(problems))
		os.Exit(1)
	}
	fmt.Println("verify: ok")
}

func run(root string) ([]string, error) {
	raw, err := os.ReadFile(filepath.Join(root, "index.json"))
	if err != nil {
		return nil, fmt.Errorf("read page index: %w", err)
	}
	var idx index
	if err := json.Unmarshal(raw, &idx); err != nil {
		return nil, fmt.Errorf("parse page index: %w", err)
	}
	if idx.Count == 0 {
		return nil, fmt.Errorf("page index is empty; expected at least one page")
	}
	if idx.Count != len(idx.Pages) {
		return nil, fmt.Errorf("page index count %d disagrees with %d entries", idx.Count, len(idx.Pages))
	}

	var problems []string

	for _, p := range idx.Pages {
		html := pagePathFor(root, p.URL)
		body, err := os.ReadFile(html)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: no page on disk at %s", p.URL, html))
			continue
		}
		m := titleRe.FindSubmatch(body)
		if m == nil || strings.TrimSpace(string(m[1])) == "" {
			problems = append(problems, fmt.Sprintf("%s: empty or missing <title>", p.URL))
		}
		if strings.TrimSpace(p.Title) == "" {
			problems = append(problems, fmt.Sprintf("%s: empty title in page index", p.URL))
		}
		if p.URL == "/" {
			problems = append(problems, homepageProblems(body)...)
		}
		if p.Section == "writing" {
			problems = append(problems, deprecatedClassificationProblems(body)...)
		}

		if p.OG == "" {
			problems = append(problems, fmt.Sprintf("%s: no social card declared", p.URL))
			continue
		}
		og, err := assetPathFor(root, p.OG)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: unparseable social card url %q: %v", p.URL, p.OG, err))
			continue
		}
		if st, err := os.Stat(og); err != nil || st.Size() == 0 {
			problems = append(problems, fmt.Sprintf("%s: social card missing or empty at %s", p.URL, og))
		}
	}

	writingIndex, err := os.ReadFile(filepath.Join(root, "writing", "index.html"))
	if err != nil {
		problems = append(problems, "writing index missing from build output")
	} else {
		problems = append(problems, deprecatedClassificationProblems(writingIndex)...)
	}

	cname, err := os.ReadFile(filepath.Join(root, "CNAME"))
	if err != nil {
		problems = append(problems, "CNAME missing from build output; the custom domain would break")
	} else if got := strings.TrimSpace(string(cname)); got != "moknshaik.com" {
		problems = append(problems, fmt.Sprintf("CNAME is %q, want \"moknshaik.com\"", got))
	}

	for _, required := range []string{"404.html", "rss.xml", "sitemap.xml"} {
		if _, err := os.Stat(filepath.Join(root, required)); err != nil {
			problems = append(problems, fmt.Sprintf("%s missing from build output", required))
		}
	}

	return problems, nil
}

func homepageProblems(body []byte) []string {
	html := string(body)
	var problems []string
	required := []struct {
		text    string
		problem string
	}{
		{"$ ls projects/ --featured", "homepage: featured projects command missing"},
		{`role="switch"`, "homepage: theme switch missing"},
		{`aria-checked=`, "homepage: theme switch state missing"},
	}
	for _, requirement := range required {
		if !strings.Contains(html, requirement.text) {
			problems = append(problems, requirement.problem)
		}
	}
	for _, social := range []struct {
		url   string
		label string
	}{
		{"https://github.com/pdwytr", "GitHub"},
		{"https://x.com/pdwytrfa", "X (Twitter)"},
		{"https://www.linkedin.com/in/pdwytr/", "LinkedIn"},
	} {
		if !hasSocialLink(html, social.url, social.label) {
			problems = append(problems, fmt.Sprintf("homepage: %s social link missing", social.label))
		}
	}
	if strings.Contains(html, "$ ls -t writing/") {
		problems = append(problems, "homepage: recent writing block still present")
	}
	return problems
}

func hasSocialLink(html, linkURL, label string) bool {
	quotedURL := regexp.QuoteMeta(linkURL)
	quotedLabel := regexp.QuoteMeta(label)
	pattern := `(?is)<a\b[^>]*\bhref=(?:"` + quotedURL + `"|` + quotedURL + `)(?:\s|>)[^>]*\baria-label=(?:"` + quotedLabel + `"|` + quotedLabel + `)(?:\s|>)`
	return regexp.MustCompile(pattern).MatchString(html)
}

func deprecatedClassificationProblems(body []byte) []string {
	html := string(body)
	if strings.Contains(html, "[essay]") || strings.Contains(html, "[note]") || strings.Contains(html, "row__kind") {
		return []string{"writing: deprecated essay/note classification still rendered"}
	}
	return nil
}

// pagePathFor maps a site-relative page URL to its file on disk.
func pagePathFor(root, u string) string {
	clean := strings.Trim(u, "/")
	if clean == "" {
		return filepath.Join(root, "index.html")
	}
	return filepath.Join(root, filepath.FromSlash(clean), "index.html")
}

// assetPathFor maps an absolute or root-relative asset URL to its file on disk.
func assetPathFor(root, raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	p := strings.TrimPrefix(parsed.Path, "/")
	if p == "" {
		return "", fmt.Errorf("no path component")
	}
	return filepath.Join(root, filepath.FromSlash(p)), nil
}
