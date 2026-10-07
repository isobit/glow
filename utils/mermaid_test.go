package utils

import (
	"strings"
	"testing"
)

func TestRenderMermaidDiagrams(t *testing.T) {
	markdown := "# Flow\n\n```mermaid\ngraph LR\nA[Start] --> B[Finish]\n```\n\nDone.\n"

	got := RenderMermaidDiagrams(markdown, 80)
	if strings.Contains(got, "graph LR") || strings.Contains(got, "A[Start]") {
		t.Fatalf("Mermaid source remains in output:\n%s", got)
	}
	for _, want := range []string{"```text", "Start", "Finish", "Done."} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q:\n%s", want, got)
		}
	}
}

func TestRenderMermaidDiagramsLeavesOtherFencesAlone(t *testing.T) {
	markdown := "```go\npackage main\n```\n"
	if got := RenderMermaidDiagrams(markdown, 80); got != markdown {
		t.Fatalf("non-Mermaid code block changed:\n%s", got)
	}

	markdown = "````\n```mermaid\ngraph LR\nA --> B\n```\n````\n"
	if got := RenderMermaidDiagrams(markdown, 80); got != markdown {
		t.Fatalf("nested Mermaid-like fence was transformed:\n%s", got)
	}
}

func TestRenderMermaidDiagramsLeavesUnclosedFenceAlone(t *testing.T) {
	markdown := "```mermaid\ngraph LR\nA --> B\n"
	if got := RenderMermaidDiagrams(markdown, 80); got != markdown {
		t.Fatalf("unclosed Mermaid fence changed:\n%s", got)
	}
}
