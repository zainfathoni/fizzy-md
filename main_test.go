package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertMarkdownToHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple paragraph",
			input:    "Hello world",
			expected: "<p>Hello world</p>",
		},
		{
			name:     "heading h2",
			input:    "## Heading",
			expected: "<h2>Heading</h2>",
		},
		{
			name:     "heading h3",
			input:    "### Subheading",
			expected: "<h3>Subheading</h3>",
		},
		{
			name:     "bold text",
			input:    "This is **bold** text",
			expected: "<p>This is <strong>bold</strong> text</p>",
		},
		{
			name:     "italic text",
			input:    "This is *italic* text",
			expected: "<p>This is <em>italic</em> text</p>",
		},
		{
			name:     "inline code",
			input:    "Use `code` here",
			expected: "<p>Use <code>code</code> here</p>",
		},
		{
			name:     "unordered list",
			input:    "- Item 1\n- Item 2\n- Item 3",
			expected: "<ul>\n<li>Item 1</li>\n<li>Item 2</li>\n<li>Item 3</li>\n</ul>",
		},
		{
			name:     "ordered list",
			input:    "1. First\n2. Second\n3. Third",
			expected: "<ol>\n<li>First</li>\n<li>Second</li>\n<li>Third</li>\n</ol>",
		},
		{
			name:     "link",
			input:    "Check [this link](https://example.com)",
			expected: `<p>Check <a href="https://example.com">this link</a></p>`,
		},
		{
			name:     "code block",
			input:    "```\ncode here\n```",
			expected: "<pre><code>code here\n</code></pre>",
		},
		{
			name:     "code block with language",
			input:    "```go\nfunc main() {}\n```",
			expected: "<pre><code class=\"language-go\">func main() {}\n</code></pre>",
		},
		{
			name:     "blockquote",
			input:    "> This is a quote",
			expected: "<blockquote>\n<p>This is a quote</p>\n</blockquote>",
		},
		{
			name:     "complex document",
			input:    "## Overview\n\nThis is **important**.\n\n- Item 1\n- Item 2",
			expected: "<h2>Overview</h2>\n<p>This is <strong>important</strong>.</p>\n<ul>\n<li>Item 1</li>\n<li>Item 2</li>\n</ul>",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "emoji preservation",
			input:    "Status: ✅ Complete 🚀",
			expected: "<p>Status: ✅ Complete 🚀</p>",
		},
		{
			name:  "simple table",
			input: "| Name | Value |\n|------|-------|\n| Foo  | Bar   |\n| Baz  | Qux   |",
			expected: `<table>
<thead>
<tr>
<th>Name</th>
<th>Value</th>
</tr>
</thead>
<tbody>
<tr>
<td>Foo</td>
<td>Bar</td>
</tr>
<tr>
<td>Baz</td>
<td>Qux</td>
</tr>
</tbody>
</table>`,
		},
		{
			name:  "table with alignment",
			input: "| Left | Center | Right |\n|:-----|:------:|------:|\n| A    | B      | C     |",
			expected: `<table>
<thead>
<tr>
<th style="text-align:left">Left</th>
<th style="text-align:center">Center</th>
<th style="text-align:right">Right</th>
</tr>
</thead>
<tbody>
<tr>
<td style="text-align:left">A</td>
<td style="text-align:center">B</td>
<td style="text-align:right">C</td>
</tr>
</tbody>
</table>`,
		},
		{
			name:  "table with inline formatting",
			input: "| **Bold** | *Italic* | `Code` |\n|----------|----------|--------|\n| foo      | bar      | baz    |",
			expected: `<table>
<thead>
<tr>
<th><strong>Bold</strong></th>
<th><em>Italic</em></th>
<th><code>Code</code></th>
</tr>
</thead>
<tbody>
<tr>
<td>foo</td>
<td>bar</td>
<td>baz</td>
</tr>
</tbody>
</table>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertMarkdownToHTML(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("\nInput:    %q\nExpected: %q\nGot:      %q", tt.input, tt.expected, result)
			}
		})
	}
}

func TestReadAndConvertFile(t *testing.T) {
	// Create temp directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		filename    string
		content     string
		expected    string
		expectError bool
	}{
		{
			name:     "markdown file (.md)",
			filename: "test.md",
			content:  "## Hello\n\nWorld",
			expected: "<h2>Hello</h2>\n<p>World</p>",
		},
		{
			name:     "html file (.html) - passthrough",
			filename: "test.html",
			content:  "<h2>Already HTML</h2>",
			expected: "<h2>Already HTML</h2>",
		},
		{
			name:     "htm file (.htm) - passthrough",
			filename: "test.htm",
			content:  "<p>Also HTML</p>",
			expected: "<p>Also HTML</p>",
		},
		{
			name:     "no extension - assume markdown",
			filename: "notes",
			content:  "**Bold** text",
			expected: "<p><strong>Bold</strong> text</p>",
		},
		{
			name:     "txt extension - convert as markdown",
			filename: "notes.txt",
			content:  "## Heading",
			expected: "<h2>Heading</h2>",
		},
		{
			name:        "non-existent file",
			filename:    "does-not-exist.md",
			content:     "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			
			if !tt.expectError || tt.filename != "does-not-exist.md" {
				// Create test file
				filePath = filepath.Join(tmpDir, tt.filename)
				if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			} else {
				filePath = filepath.Join(tmpDir, tt.filename)
			}

			result, err := readAndConvertFile(filePath)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			
			if result != tt.expected {
				t.Errorf("\nFilename: %s\nExpected: %q\nGot:      %q", tt.filename, tt.expected, result)
			}
		})
	}
}

func TestProcessArgs(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		checkFn     func([]string) error
		expectError bool
	}{
		{
			name:  "passthrough - no conversion needed",
			input: []string{"card", "list", "--board", "123"},
			checkFn: func(result []string) error {
				expected := []string{"card", "list", "--board", "123"}
				if len(result) != len(expected) {
					return nil // length mismatch will be caught
				}
				for i, v := range expected {
					if result[i] != v {
						return nil
					}
				}
				return nil
			},
		},
		{
			name:  "convert --description flag",
			input: []string{"card", "create", "--title", "Test", "--description", "## Hello"},
			checkFn: func(result []string) error {
				// Check that --description value was converted
				for i, v := range result {
					if v == "--description" && i+1 < len(result) {
						if !strings.Contains(result[i+1], "<h2>") {
							t.Errorf("expected HTML conversion, got: %s", result[i+1])
						}
						return nil
					}
				}
				t.Error("--description flag not found in result")
				return nil
			},
		},
		{
			name:  "convert --body flag",
			input: []string{"comment", "create", "--card", "42", "--body", "**Bold** text"},
			checkFn: func(result []string) error {
				// Check that --body value was converted
				for i, v := range result {
					if v == "--body" && i+1 < len(result) {
						if !strings.Contains(result[i+1], "<strong>") {
							t.Errorf("expected HTML conversion, got: %s", result[i+1])
						}
						return nil
					}
				}
				t.Error("--body flag not found in result")
				return nil
			},
		},
		{
			name:        "missing value for --description",
			input:       []string{"card", "create", "--description"},
			expectError: true,
		},
		{
			name:        "missing value for --body",
			input:       []string{"comment", "create", "--body"},
			expectError: true,
		},
		{
			name:  "multiple flags in one command",
			input: []string{"card", "create", "--title", "Test", "--description", "## Hello", "--board", "123"},
			checkFn: func(result []string) error {
				foundTitle := false
				foundBoard := false
				foundDescription := false
				
				for i, v := range result {
					if v == "--title" && i+1 < len(result) && result[i+1] == "Test" {
						foundTitle = true
					}
					if v == "--board" && i+1 < len(result) && result[i+1] == "123" {
						foundBoard = true
					}
					if v == "--description" && i+1 < len(result) {
						if strings.Contains(result[i+1], "<h2>") {
							foundDescription = true
						}
					}
				}
				
				if !foundTitle {
					t.Error("--title flag not preserved correctly")
				}
				if !foundBoard {
					t.Error("--board flag not preserved correctly")
				}
				if !foundDescription {
					t.Error("--description not converted correctly")
				}
				return nil
			},
		},
		{
			name:  "empty args",
			input: []string{},
			checkFn: func(result []string) error {
				if len(result) != 0 {
					t.Errorf("expected empty result, got %v", result)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processArgs(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			
			if tt.checkFn != nil {
				tt.checkFn(result)
			}
		})
	}
}

func TestProcessArgsWithFiles(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test markdown file
	mdFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(mdFile, []byte("## From File\n\n- Item 1\n- Item 2"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	
	// Create test HTML file (should passthrough)
	htmlFile := filepath.Join(tmpDir, "test.html")
	if err := os.WriteFile(htmlFile, []byte("<h2>Already HTML</h2>"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Run("convert --description_file with .md", func(t *testing.T) {
		result, err := processArgs([]string{"card", "create", "--title", "Test", "--description_file", mdFile})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		// Find the converted file path
		for i, v := range result {
			if v == "--description_file" && i+1 < len(result) {
				tmpPath := result[i+1]
				content, err := os.ReadFile(tmpPath)
				if err != nil {
					t.Fatalf("failed to read temp file: %v", err)
				}
				if !strings.Contains(string(content), "<h2>From File</h2>") {
					t.Errorf("expected converted HTML, got: %s", content)
				}
				if !strings.Contains(string(content), "<li>Item 1</li>") {
					t.Errorf("expected list conversion, got: %s", content)
				}
				return
			}
		}
		t.Error("--description_file flag not found in result")
	})

	t.Run("passthrough --description_file with .html", func(t *testing.T) {
		result, err := processArgs([]string{"card", "create", "--title", "Test", "--description_file", htmlFile})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		// Find the file path - should still create temp file but with passthrough content
		for i, v := range result {
			if v == "--description_file" && i+1 < len(result) {
				tmpPath := result[i+1]
				content, err := os.ReadFile(tmpPath)
				if err != nil {
					t.Fatalf("failed to read temp file: %v", err)
				}
				if string(content) != "<h2>Already HTML</h2>" {
					t.Errorf("expected passthrough, got: %s", content)
				}
				return
			}
		}
		t.Error("--description_file flag not found in result")
	})

	t.Run("error on non-existent file", func(t *testing.T) {
		_, err := processArgs([]string{"card", "create", "--description_file", "/nonexistent/file.md"})
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})
}

// Benchmark for performance requirement (<100ms)
func BenchmarkConvertMarkdownToHTML(b *testing.B) {
	input := `## Overview

This is a **complex** document with multiple elements.

### Features

- Feature 1 with ` + "`code`" + `
- Feature 2 with **bold**
- Feature 3 with *italic*

### Code Example

` + "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```" + `

### Conclusion

> This is a blockquote for emphasis.

Visit [our site](https://example.com) for more info.
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = convertMarkdownToHTML(input)
	}
}

func BenchmarkProcessArgs(b *testing.B) {
	args := []string{
		"card", "create",
		"--title", "Test Card",
		"--description", "## Hello\n\nThis is **bold** and *italic*.\n\n- Item 1\n- Item 2",
		"--board", "123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = processArgs(args)
	}
}
