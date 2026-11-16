package ui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestGlamourRenderTextSizing(t *testing.T) {
	t.Setenv("GLOW_TEXT_SIZING", "on")

	common := &commonModel{
		cfg: Config{
			GlamourEnabled:  true,
			GlamourStyle:    "dark",
			GlamourMaxWidth: 80,
		},
		styles: newStyles(true),
		width:  80,
		height: 24,
	}
	m := newPagerModel(common)
	m.currentDocument = markdown{Note: "test.md"}
	m.setSize(80, 24)

	out, err := glamourRender(m, "# Hello\n\n## World\n")
	if err != nil {
		t.Fatal(err)
	}
	// The dark style's H1 has a background color, which renders broken when
	// scaled, so it is left at normal size.
	if strings.Contains(out, "\x1b]66;s=3;") {
		t.Errorf("expected H1 with background to not be scaled, got %q", out)
	}
	if !strings.Contains(out, "Hello") {
		t.Errorf("expected H1 text in pager output, got %q", out)
	}
	if !strings.Contains(out, "\x1b]66;s=2;World\x1b\\") {
		t.Errorf("expected H2 at 2x scale in pager output, got %q", out)
	}
	if strings.Contains(out, "\x1b]6666;") {
		t.Errorf("expected no leftover markers in pager output, got %q", out)
	}
}

func TestRenderWithGlamourWaitsForWindowSize(t *testing.T) {
	common := &commonModel{
		cfg: Config{
			GlamourEnabled:  true,
			GlamourStyle:    "notty",
			GlamourMaxWidth: 40,
		},
		styles: newStyles(true),
	}
	m := newPagerModel(common)
	m.currentDocument = markdown{Note: "test.md"}
	body := strings.Repeat("word ", 40)

	// Before the first WindowSizeMsg there is no width to wrap at.
	if cmd := renderWithGlamour(m, body); cmd != nil {
		t.Fatal("expected no render before the viewport has a width")
	}

	m.setSize(200, 24)
	msg, ok := renderWithGlamour(m, body)().(contentRenderedMsg)
	if !ok {
		t.Fatalf("expected contentRenderedMsg, got %T", msg)
	}
	for _, line := range strings.Split(string(msg), "\n") {
		if w := lipgloss.Width(line); w > 40 {
			t.Errorf("line wider than GlamourMaxWidth (%d > 40): %q", w, line)
		}
	}
}
