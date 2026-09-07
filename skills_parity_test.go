package main

// skills_parity_test.go pins the published docs against each other: the LobeHub
// manifest (lhm.plugin.json), the skill (skills/grimoire-campaign-context) and
// the README all describe the same server, and every one of them is written by
// hand.
//
// The 2026-09-06 live MCP test found the chain had drifted in three places at
// once: the reference said batch block creates need a verification and reorder
// pass while the tool description said the opposite, the reference listed four
// wiki visibility values where the store has five, and the graph type list was
// short by one. Nothing caught it, because nothing compared the files.
//
// This is the docs-side guard. The server-side guard is
// backend/internal/mcp/tools_description_parity_test.go in the Grimoire
// monorepo, which pins the description strings against the behaviour they
// describe. When the two disagree, the server is right and these files are
// wrong: SKILL.md itself names the live tool description as the tie-breaker.
//
// Stdlib only, no fixtures, no network.

import (
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const (
	manifestPath  = "lhm.plugin.json"
	skillPath     = "skills/grimoire-campaign-context/SKILL.md"
	toolGroupPath = "skills/grimoire-campaign-context/references/tool-groups.md"
	readmePath    = "README.md"
)

// readOnlyTools is the same list backend/internal/mcp's TestReadOnlyToolSet
// pins on the server. It is duplicated here on purpose: this repo cannot import
// the backend, and a read tool quietly becoming a write tool is exactly the
// drift these two tests exist to catch.
var readOnlyTools = []string{
	"current_campaign",
	"get_campaign_bible",
	"get_campaign_context",
	"get_constitution",
	"get_entity",
	"get_entity_catalog",
	"get_entity_graph",
	"get_entity_schema",
	"get_field_options",
	"get_knowledge_graph",
	"get_narrative_state",
	"get_open_threads",
	"get_relationships",
	"get_tag_options",
	"get_thread_progressions",
	"get_wiki_page",
	"get_wiki_tree",
	"list_entities",
	"list_entity_graphs",
	"search_campaign",
	"search_wiki",
}

type manifest struct {
	Version string `json:"version"`
	Tools   []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"tools"`
}

func loadManifest(t *testing.T) manifest {
	t.Helper()

	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read %s: %v", manifestPath, err)
	}

	var m manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("parse %s: %v", manifestPath, err)
	}
	if len(m.Tools) == 0 {
		t.Fatalf("%s declares no tools", manifestPath)
	}
	return m
}

func loadFile(t *testing.T, path string) string {
	t.Helper()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(raw)
}

func toolDescription(t *testing.T, m manifest, name string) string {
	t.Helper()

	for _, tool := range m.Tools {
		if tool.Name == name {
			return tool.Description
		}
	}
	t.Fatalf("%s declares no tool named %q", manifestPath, name)
	return ""
}

// TestToolCountsAgree checks every hand-written tool count against the manifest.
func TestToolCountsAgree(t *testing.T) {
	m := loadManifest(t)
	want := len(m.Tools)

	countRe := regexp.MustCompile(`(\d+) tools`)
	for _, path := range []string{skillPath, toolGroupPath, readmePath} {
		body := loadFile(t, path)
		matches := countRe.FindAllStringSubmatch(body, -1)
		if len(matches) == 0 {
			t.Errorf("%s states no tool count; the manifest declares %d", path, want)
			continue
		}
		for _, match := range matches {
			got, err := strconv.Atoi(match[1])
			if err != nil {
				t.Fatalf("%s: %v", path, err)
			}
			if got != want {
				t.Errorf("%s says %q, manifest declares %d tools", path, match[0], want)
			}
		}
	}

	// "21 read, 28 write" must add up and must match the read-only set.
	splitRe := regexp.MustCompile(`(\d+) read, (\d+) write`)
	for _, path := range []string{toolGroupPath, readmePath} {
		body := loadFile(t, path)
		for _, match := range splitRe.FindAllStringSubmatch(body, -1) {
			read, _ := strconv.Atoi(match[1])
			write, _ := strconv.Atoi(match[2])
			if read != len(readOnlyTools) {
				t.Errorf("%s says %q, the read-only tool set holds %d", path, match[0], len(readOnlyTools))
			}
			if read+write != want {
				t.Errorf("%s says %q, which is not the %d tools the manifest declares", path, match[0], want)
			}
		}
	}
}

// TestEveryManifestToolIsDocumented checks the reference covers the surface.
func TestEveryManifestToolIsDocumented(t *testing.T) {
	m := loadManifest(t)
	body := loadFile(t, toolGroupPath)

	for _, tool := range m.Tools {
		if !strings.Contains(body, tool.Name) {
			t.Errorf("%s never mentions %q", toolGroupPath, tool.Name)
		}
	}
}

// TestBlockVisibilityValuesAgree pins the wiki visibility vocabulary. The
// manifest copies it from the server, so the server is the source here.
func TestBlockVisibilityValuesAgree(t *testing.T) {
	m := loadManifest(t)
	desc := toolDescription(t, m, "batch_create_wiki_blocks")

	valueRe := regexp.MustCompile(`(?m)^- ([a-z-]+): `)
	var values []string
	for _, match := range valueRe.FindAllStringSubmatch(desc, -1) {
		switch match[1] {
		case "inherit", "common-knowledge", "dm-secret", "player-knowledge", "system":
			values = append(values, match[1])
		}
	}
	if len(values) == 0 {
		t.Fatalf("no visibility values found in batch_create_wiki_blocks description")
	}

	groups := loadFile(t, toolGroupPath)
	paragraph := blockVisibilityParagraph(t, groups)
	readme := loadFile(t, readmePath)

	for _, value := range values {
		if !strings.Contains(paragraph, value) {
			t.Errorf("%s: the Block visibility paragraph omits %q", toolGroupPath, value)
		}
		if !strings.Contains(readme, value) {
			t.Errorf("%s: the visibility sentence omits %q", readmePath, value)
		}
	}
}

// blockVisibilityParagraph returns the "Block visibility:" paragraph of the
// tool reference, so a value mentioned elsewhere in the file cannot satisfy the
// assertion by accident.
func blockVisibilityParagraph(t *testing.T, body string) string {
	t.Helper()

	start := strings.Index(body, "Block visibility:")
	if start < 0 {
		t.Fatalf("%s has no Block visibility paragraph", toolGroupPath)
	}
	rest := body[start:]
	if end := strings.Index(rest, "\n\n"); end >= 0 {
		return rest[:end]
	}
	return rest
}

// TestBatchBlockOrderingClaimsAgree catches the contradiction the live test
// found: the reference told clients to verify order and reorder afterwards
// while the tool description said the array order is already the stored order.
func TestBatchBlockOrderingClaimsAgree(t *testing.T) {
	m := loadManifest(t)
	desc := toolDescription(t, m, "batch_create_wiki_blocks")

	if !strings.Contains(desc, "no verification or reorder pass") {
		t.Skip("batch_create_wiki_blocks no longer promises array order; nothing to contradict")
	}

	body := loadFile(t, toolGroupPath)
	for _, stale := range []string{"verify order", "positioning race"} {
		if strings.Contains(body, stale) {
			t.Errorf("%s still says %q, but the tool description promises array order", toolGroupPath, stale)
		}
	}
}

// TestGraphTypesAgree pins the knowledge graph type list.
func TestGraphTypesAgree(t *testing.T) {
	m := loadManifest(t)
	desc := toolDescription(t, m, "get_knowledge_graph")

	block := desc
	if start := strings.Index(desc, "GRAPH TYPES:"); start >= 0 {
		block = desc[start:]
	}
	typeRe := regexp.MustCompile(`(?m)^- '([a-z_]+)':`)
	matches := typeRe.FindAllStringSubmatch(block, -1)
	if len(matches) == 0 {
		t.Fatalf("no graph types found in the get_knowledge_graph description")
	}

	row := toolGroupsRow(t, "get_knowledge_graph")
	for _, match := range matches {
		if !strings.Contains(row, match[1]) {
			t.Errorf("%s: the get_knowledge_graph row omits graph type %q", toolGroupPath, match[1])
		}
	}
}

// toolGroupsRow returns the table row of the reference that names the tool.
func toolGroupsRow(t *testing.T, tool string) string {
	t.Helper()

	body := loadFile(t, toolGroupPath)
	needle := "| `" + tool + "` |"
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, needle) {
			return line
		}
	}
	t.Fatalf("%s has no table row for %q", toolGroupPath, tool)
	return ""
}

// TestRoleGatingTableCoversEveryReadTool pins the reference's role table
// against the read-only tool set. A read tool missing from the table is a tool
// a connected client is told nothing about.
func TestRoleGatingTableCoversEveryReadTool(t *testing.T) {
	body := loadFile(t, toolGroupPath)

	start := strings.Index(body, "## Role gating at a glance")
	if start < 0 {
		t.Fatalf("%s has no role gating section", toolGroupPath)
	}
	section := body[start:]

	nameRe := regexp.MustCompile("`([a-z_]+)`")
	named := map[string]bool{}
	for _, line := range strings.Split(section, "\n") {
		if !strings.HasPrefix(line, "| Player:") {
			continue
		}
		// Only the second cell lists tools; the first names the availability
		// rule and can carry portal policy names in backticks.
		cells := strings.Split(line, "|")
		if len(cells) < 3 {
			continue
		}
		for _, match := range nameRe.FindAllStringSubmatch(cells[2], -1) {
			named[match[1]] = true
		}
	}

	var missing []string
	for _, tool := range readOnlyTools {
		if !named[tool] {
			missing = append(missing, tool)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Errorf("%s: the role gating table's player rows omit %v", toolGroupPath, missing)
	}

	m := loadManifest(t)
	declared := map[string]bool{}
	for _, tool := range m.Tools {
		declared[tool.Name] = true
	}
	for name := range named {
		if !declared[name] {
			t.Errorf("%s: the role gating table names %q, which the manifest does not declare", toolGroupPath, name)
		}
	}
}
