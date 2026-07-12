package label_test

import (
	"bufio"
	"bytes"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/yaien/cultural/internal/application/label"
)

func TestPageBaseStyles(t *testing.T) {
	var b bytes.Buffer

	err := label.WritePageBaseStyles(&b, &label.Config{
		Fonts: map[string]*label.Font{
			"primary":  {Family: "Inter"},
			"headings": {Family: "Montserrat"},
		},
		Colors: []*label.Color{
			{Tag: "primary", Value: "#1a73e8"},
			{Tag: "secondary", Value: "#e8f0fe"},
		},
	})

	if err != nil {
		t.Fatalf("Error writing page base styles: %v", err)
	}

	expected := map[string]bool{
		`--primary-font-family: "Inter", sans-serif;`:       true,
		`--headings-font-family: "Montserrat", sans-serif;`: true,
		`--primary-color: #1a73e8;`:                         true,
		`--secondary-color: #e8f0fe;`:                       true,
	}

	first := ":root {"
	last := "}"

	s := bufio.NewScanner(&b)

	for s.Scan() {
		if err := s.Err(); err != nil {
			t.Fatalf("Error scanning output: %v", err)
		}

		line := strings.TrimSpace(s.Text())

		t.Log("line > ", line)

		if line == "" {
			continue
		}

		if first != "" {
			if line != first {
				t.Errorf("Expected first line to be '%s', got '%s'", first, line)
			}
			first = ""
			continue
		}

		if len(expected) > 0 {
			if !expected[line] {
				t.Errorf("Unexpected line in output: %s", line)
			} else {
				delete(expected, line)
			}

			continue
		}

		if last != "" {
			if line != last {
				t.Errorf("Expected last line to be '%s', got '%s'", last, line)
			}
			last = ""
			continue
		}

	}

	if len(expected) > 0 {
		t.Errorf("Missing expected lines: %v", slices.Collect(maps.Keys(expected)))
	}

}

func TestGetPageInMap(t *testing.T) {
	tests := []struct {
		name   string
		pages  []string
		params []string
		path   string
		found  bool
		page   string
	}{
		{
			name:   "index",
			pages:  []string{"index", "products", "product"},
			path:   "/",
			found:  true,
			page:   "index",
			params: []string{},
		},
		{
			name:   "product",
			pages:  []string{"index", "products", "product"},
			path:   "/product/wallpapers",
			found:  true,
			page:   "product",
			params: []string{"wallpapers"},
		},
		{
			name:   "products",
			pages:  []string{"index", "products", "product"},
			path:   "/products",
			found:  true,
			page:   "products",
			params: []string{},
		},
		{
			name:  "cars",
			pages: []string{"index", "products", "product"},
			path:  "/cars",
			found: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			pages := make(map[string]*label.Page)
			for _, page := range test.pages {
				pages[page] = &label.Page{Name: page}
			}

			page, params, found := label.GetPageInMap(pages, test.path)

			if found != test.found {
				t.Fatalf("Expected found to be %v, got %v", test.found, found)
			}

			if !test.found {
				return
			}

			if page == nil {
				t.Fatalf("Expected page to be found, got nil")
			}

			if page.Name != test.page {
				t.Fatalf("Expected page name to be %s, got %s", test.page, page.Name)
			}

			if !slices.Equal(params, test.params) {
				t.Errorf("Expected params to be %v, got %v", test.params, params)
			}
		})
	}

}
