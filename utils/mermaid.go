package utils

import (
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
	mermaidrender "github.com/AlexanderGrooff/mermaid-ascii/pkg/render"
)

// RenderMermaidDiagrams replaces Mermaid fenced code blocks with rendered diagrams.
func RenderMermaidDiagrams(markdown string, maxWidth int) string {
	lines := strings.SplitAfter(markdown, "\n")
	var output strings.Builder

	for i := 0; i < len(lines); {
		fence, ok := mermaidFence(lineContent(lines[i]))
		if !ok {
			output.WriteString(lines[i])
			i++
			continue
		}

		end := i + 1
		for end < len(lines) && !isClosingFence(lineContent(lines[end]), fence) {
			end++
		}
		if end == len(lines) {
			for _, line := range lines[i:] {
				output.WriteString(line)
			}
			break
		}

		fields := strings.Fields(fence.info)
		if len(fields) == 0 || !strings.EqualFold(fields[0], "mermaid") {
			for _, line := range lines[i : end+1] {
				output.WriteString(line)
			}
			i = end + 1
			continue
		}

		sourceLines := make([]string, 0, end-i-1)
		for _, line := range lines[i+1 : end] {
			line = lineContent(line)
			line = strings.TrimPrefix(line, fence.indent)
			sourceLines = append(sourceLines, line)
		}

		config := diagram.DefaultConfig()
		config.MaxWidth = maxWidth
		rendered, err := mermaidrender.RenderDiagram(strings.Join(sourceLines, "\n"), config)
		if err != nil {
			for _, line := range lines[i : end+1] {
				output.WriteString(line)
			}
			i = end + 1
			continue
		}

		eol := lineEnding(lines[i])
		if eol == "" {
			eol = "\n"
		}
		fenceChar := string(fence.char)
		fenceLength := max(fence.length, longestRun(rendered, fence.char)+1)
		fenceLine := fence.indent + strings.Repeat(fenceChar, fenceLength)
		output.WriteString(fenceLine)
		output.WriteString("text")
		output.WriteString(eol)
		for line := range strings.SplitSeq(strings.TrimRight(rendered, "\r\n"), "\n") {
			output.WriteString(fence.indent)
			output.WriteString(line)
			output.WriteString(eol)
		}
		output.WriteString(fenceLine)
		output.WriteString(lineEnding(lines[end]))
		i = end + 1
	}

	return output.String()
}

type markdownFence struct {
	char   rune
	length int
	indent string
	info   string
}

func mermaidFence(line string) (markdownFence, bool) {
	indentLength := len(line) - len(strings.TrimLeft(line, " "))
	if indentLength > 3 || indentLength == len(line) {
		return markdownFence{}, false
	}
	indent := line[:indentLength]
	line = line[indentLength:]
	char := rune(line[0])
	if char != '`' && char != '~' {
		return markdownFence{}, false
	}
	length := 0
	for length < len(line) && rune(line[length]) == char {
		length++
	}
	if length < 3 {
		return markdownFence{}, false
	}
	info := strings.TrimSpace(line[length:])
	if char == '`' && strings.ContainsRune(info, '`') {
		return markdownFence{}, false
	}
	return markdownFence{char: char, length: length, indent: indent, info: info}, true
}

func isClosingFence(line string, opening markdownFence) bool {
	indentLength := len(line) - len(strings.TrimLeft(line, " "))
	if indentLength > 3 || indentLength == len(line) {
		return false
	}
	line = line[indentLength:]
	length := 0
	for length < len(line) && rune(line[length]) == opening.char {
		length++
	}
	return length >= opening.length && strings.TrimSpace(line[length:]) == ""
}

func lineContent(line string) string {
	return strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
}

func lineEnding(line string) string {
	if strings.HasSuffix(line, "\r\n") {
		return "\r\n"
	}
	if strings.HasSuffix(line, "\n") {
		return "\n"
	}
	return ""
}

func longestRun(s string, char rune) int {
	longest, current := 0, 0
	for _, r := range s {
		if r == char {
			current++
			longest = max(longest, current)
		} else {
			current = 0
		}
	}
	return longest
}
