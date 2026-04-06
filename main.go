package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// envRefPattern matches both ${VAR} and $VAR references.
var envRefPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

func main() {
	var inputPath string
	var outputPath string
	var compact bool

	flag.StringVar(&inputPath, "i", "", "Input .env file path (required)")
	flag.StringVar(&outputPath, "o", "", "Output file path (optional)")
	flag.BoolVar(&compact, "c", false, "Compact output: strip comments and blank/whitespace lines")
	flag.Usage = usage
	flag.Parse()

	if inputPath == "" {
		usage()
		fmt.Fprintln(os.Stderr, "\nerror: -i input file is required")
		os.Exit(2)
	}

	if outputPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			exitWithError(err)
		}

		outputPath = filepath.Join(cwd, defaultOutputName(filepath.Base(inputPath)))
	}

	inBytes, err := os.ReadFile(inputPath)
	if err != nil {
		exitWithError(fmt.Errorf("read input file: %w", err))
	}

	processed, err := processEnvContent(string(inBytes), compact)
	if err != nil {
		exitWithError(err)
	}

	if err := os.WriteFile(outputPath, []byte(processed), 0o644); err != nil {
		exitWithError(fmt.Errorf("write output file: %w", err))
	}

	fmt.Printf("processed: %s -> %s\n", inputPath, outputPath)

}

// usage prints CLI help text with behavior and examples.
func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), `envexpand - .env variable expansion tool

Expands environment references in .env-like files, with support for:
  - ${VAR}
  - $VAR

Resolution order:
  1) variables already parsed earlier in the same input file
  2) operating-system environment variables

Unresolved variables are left unchanged.
Quoted values are expanded and quotes are preserved.

Usage:
  envexpand -i <input-file> [-o <output-file>] [-c]

Flags:
`)

	flag.PrintDefaults()

	fmt.Fprint(flag.CommandLine.Output(), `
Examples:
  envexpand -i api/.env.example
    # writes to ./<base>_out<ext> in current working directory

  envexpand -i api/.env.example -o api/.env.processed

  envexpand -i api/.env.example -c
    # compact output: strips comments and blank/whitespace-only lines
`)
}

// defaultOutputName returns <base>_out<ext> for an input filename.
func defaultOutputName(inputBase string) string {
	ext := filepath.Ext(inputBase)
	base := strings.TrimSuffix(inputBase, ext)
	return base + "_out" + ext
}

// processEnvContent parses, expands, and renders .env-like content.
//
// Expansion behavior:
//   - supports ${VAR} and $VAR
//   - resolves from already parsed keys first, then OS env
//   - preserves unresolved placeholders
//   - preserves quotes when values were originally quoted
func processEnvContent(content string, compact bool) (string, error) {
	hadTrailingNewline := strings.HasSuffix(content, "\n")
	lines := strings.Split(content, "\n")

	parsedVars := map[string]string{}
	cache := map[string]string{}
	outLines := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !compact {
				outLines = append(outLines, line)
			}
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			if !compact {
				outLines = append(outLines, line)
			}
			continue
		}

		eqIdx := strings.Index(line, "=")
		if eqIdx < 0 {
			if compact {
				outLines = append(outLines, strings.TrimSpace(line))
			} else {
				outLines = append(outLines, line)
			}
			continue
		}

		left := line[:eqIdx]
		right := line[eqIdx+1:]

		key, ok := parseKey(left)
		if !ok {
			if compact {
				outLines = append(outLines, strings.TrimSpace(line))
			} else {
				outLines = append(outLines, line)
			}
			continue
		}

		leadSpaces, rawValue, trailSpaces := splitOuterSpaces(right)
		valueBody, quoteChar, quoted := unquote(rawValue)

		// Use a fresh visiting map for each top-level value expansion.
		expandedValue := expandValue(valueBody, func(name string) (string, bool) {
			return resolveVar(name, parsedVars, cache, map[string]bool{})
		})

		parsedVars[key] = expandedValue
		cache[key] = expandedValue

		rendered := expandedValue
		if quoted {
			rendered = string(quoteChar) + rendered + string(quoteChar)
		}

		if compact {
			outLines = append(outLines, key+"="+rendered)
			continue
		}

		outLines = append(outLines, left+"="+leadSpaces+rendered+trailSpaces)
	}

	out := strings.Join(outLines, "\n")
	if hadTrailingNewline && !strings.HasSuffix(out, "\n") {
		out += "\n"
	}

	return out, nil
}

// parseKey extracts a normalized key from the left side of KEY=VALUE.
// It accepts optional "export " prefixes.
func parseKey(left string) (string, bool) {
	keyPart := strings.TrimSpace(left)
	if strings.HasPrefix(keyPart, "export ") {
		keyPart = strings.TrimSpace(strings.TrimPrefix(keyPart, "export "))
	}

	if keyPart == "" {
		return "", false
	}

	return keyPart, true
}

// splitOuterSpaces separates leading/trailing whitespace from content.
func splitOuterSpaces(s string) (leading string, middle string, trailing string) {
	start := 0
	for start < len(s) && isSpace(s[start]) {
		start++
	}

	end := len(s)
	for end > start && isSpace(s[end-1]) {
		end--
	}

	return s[:start], s[start:end], s[end:]
}

// isSpace reports whether b is considered outer-space padding in values.
func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r'
}

// unquote removes matching single or double quotes from both ends.
// It returns the unquoted value, quote character, and whether unquoting happened.
func unquote(v string) (string, byte, bool) {
	if len(v) >= 2 {
		first := v[0]
		last := v[len(v)-1]
		if (first == '"' || first == '\'') && first == last {
			return v[1 : len(v)-1], first, true
		}
	}

	return v, 0, false
}

// resolveVar resolves a variable by checking cached/parsed values first,
// then falling back to OS environment variables.
//
// visiting guards recursive expansion from infinite cycles.
func resolveVar(name string, parsedVars map[string]string, cache map[string]string, visiting map[string]bool) (string, bool) {
	if name == "" {
		return "", false
	}

	if cached, ok := cache[name]; ok {
		return cached, true
	}

	if visiting[name] {
		// Circular reference detected in current expansion path.
		return "", false
	}

	if v, ok := parsedVars[name]; ok {
		visiting[name] = true
		expanded := expandValue(v, func(inner string) (string, bool) {
			return resolveVar(inner, parsedVars, cache, visiting)
		})
		delete(visiting, name)
		cache[name] = expanded
		return expanded, true
	}

	if v, ok := os.LookupEnv(name); ok {
		expanded := expandValue(v, func(inner string) (string, bool) {
			if inner == name {
				// Avoid immediate self-reference loops from OS values.
				return "", false
			}
			return resolveVar(inner, parsedVars, cache, map[string]bool{name: true})
		})
		return expanded, true
	}

	return "", false
}

// expandValue replaces variable references using resolver.
// Unknown names are kept unchanged.
func expandValue(input string, resolver func(string) (string, bool)) string {
	return envRefPattern.ReplaceAllStringFunc(input, func(token string) string {
		name := ""
		if strings.HasPrefix(token, "${") && strings.HasSuffix(token, "}") {
			name = token[2 : len(token)-1]
		} else if strings.HasPrefix(token, "$") {
			name = token[1:]
		}

		if name == "" {
			return token
		}

		if resolved, ok := resolver(name); ok {
			return resolved
		}

		return token
	})
}

// exitWithError prints a user-facing error and exits non-zero.
func exitWithError(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	}

	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)

}
