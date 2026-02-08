package models

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestPageBaseStyles(t *testing.T) {
	var w strings.Builder

	err := WritePageBaseStyles(&w, &Config{
		Fonts: map[string]*Font{
			"primary":  {Family: "Inter"},
			"headings": {Family: "Montserrat"},
		},
		Colors: map[string]string{
			"primary":   "#1a73e8",
			"secondary": "#e8f0fe",
		},
	})

	if err != nil {
		t.Fatalf("Error writing page base styles: %v", err)
	}

	expected := map[string]bool{
		`--font-primary: "Inter", sans-serif;`:       true,
		`--font-headings: "Montserrat", sans-serif;`: true,
		`--color-primary: #1a73e8;`:                  true,
		`--color-secondary: #e8f0fe;`:                true,
	}

	lines := strings.Split(w.String(), "\n")
	t.Log("lines", w.String())

	if lines[0] != ":root {" {
		t.Errorf("Expected :root {, got %s", lines[0])
	}

	if lines[len(lines)-2] != "}" {
		t.Errorf("Expected }, got %s", lines[len(lines)-2])
	}

	for _, line := range lines[1 : len(lines)-2] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !expected[line] {
			t.Errorf("Unexpected line in output: %s", line)
		} else {
			delete(expected, line)
		}
	}

	if len(expected) > 0 {
		t.Errorf("Missing expected lines: %v", slices.Collect(maps.Keys(expected)))
	}

}
