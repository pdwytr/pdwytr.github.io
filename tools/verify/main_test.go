package main

import "testing"

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
