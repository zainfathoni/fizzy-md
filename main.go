package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var version = "dev" // Set by GoReleaser at build time

// Markdown to HTML converter using goldmark
var md goldmark.Markdown

func init() {
	// Configure goldmark with table support and safe HTML rendering
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.Table, // Enable Markdown table support
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // Allow raw HTML passthrough
		),
	)
}

// convertMarkdownToHTML converts Markdown text to HTML
func convertMarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", fmt.Errorf("markdown conversion failed: %w", err)
	}
	
	// Trim trailing newline (goldmark adds it)
	html := strings.TrimSuffix(buf.String(), "\n")
	return html, nil
}

// readAndConvertFile reads a file and converts it to HTML if it's Markdown
func readAndConvertFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	
	// If .html file, return as-is
	if ext == ".html" || ext == ".htm" {
		return string(content), nil
	}
	
	// If .md file or no extension, convert Markdown to HTML
	if ext == ".md" || ext == "" {
		return convertMarkdownToHTML(string(content))
	}
	
	// Unknown extension - assume Markdown for safety (agent-friendly default)
	return convertMarkdownToHTML(string(content))
}

// processArgs intercepts and modifies arguments for Markdown conversion
func processArgs(args []string) ([]string, error) {
	result := make([]string, 0, len(args))
	i := 0
	
	for i < len(args) {
		arg := args[i]
		
		// Check for flags that need conversion
		if arg == "--description" || arg == "--body" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag %s requires a value", arg)
			}
			
			// Convert inline Markdown text to HTML
			markdown := args[i+1]
			html, err := convertMarkdownToHTML(markdown)
			if err != nil {
				return nil, fmt.Errorf("converting %s: %w", arg, err)
			}
			
			result = append(result, arg, html)
			i += 2
			continue
		}
		
		if arg == "--description_file" || arg == "--body_file" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag %s requires a file path", arg)
			}
			
			// Read file and convert to HTML
			filePath := args[i+1]
			html, err := readAndConvertFile(filePath)
			if err != nil {
				return nil, err
			}
			
			// Create temp file with HTML content
			tmpFile, err := os.CreateTemp("", "fizzy-md-*.html")
			if err != nil {
				return nil, fmt.Errorf("failed to create temp file: %w", err)
			}
			defer tmpFile.Close()
			
			if _, err := tmpFile.WriteString(html); err != nil {
				return nil, fmt.Errorf("failed to write temp file: %w", err)
			}
			
			// Replace flag with temp file path
			result = append(result, arg, tmpFile.Name())
			i += 2
			continue
		}
		
		// Pass through all other arguments unchanged
		result = append(result, arg)
		i++
	}
	
	return result, nil
}

func main() {
	// Get original args (skip program name)
	args := os.Args[1:]
	
	// Handle --version flag
	if len(args) == 1 && (args[0] == "--version" || args[0] == "-v") {
		fmt.Printf("fizzy-md version %s\n", version)
		os.Exit(0)
	}
	
	// Handle stdin mode: if no args and stdin is piped, convert and output
	if len(args) == 0 {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Stdin is piped, not a terminal
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(os.Stdin); err != nil {
				fmt.Fprintf(os.Stderr, "fizzy-md error: failed to read stdin: %v\n", err)
				os.Exit(1)
			}
			
			html, err := convertMarkdownToHTML(buf.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "fizzy-md error: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Print(html)
			os.Exit(0)
		}
	}
	
	// Process args for Markdown conversion
	processedArgs, err := processArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fizzy-md error: %v\n", err)
		os.Exit(1)
	}
	
	// Find fizzy executable
	fizzyPath, err := exec.LookPath("fizzy")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fizzy-md error: fizzy command not found in PATH\n")
		fmt.Fprintf(os.Stderr, "Please install fizzy-cli: https://github.com/robzolkos/fizzy-cli\n")
		os.Exit(1)
	}
	
	// Execute real fizzy with processed args
	cmd := exec.Command(fizzyPath, processedArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "fizzy-md error: %v\n", err)
		os.Exit(1)
	}
}
