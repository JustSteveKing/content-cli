package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/juststeveking/content-cli/internal/content"
	"github.com/juststeveking/content-cli/internal/schema"
	"github.com/juststeveking/content-cli/internal/slug"
	tmpl "github.com/juststeveking/content-cli/internal/template"
	"github.com/juststeveking/content-cli/internal/validate"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start an MCP server for agent/AI tool integration",
	Long: `Starts a Model Context Protocol (MCP) server over stdio.
Register with Claude Desktop or any MCP-compatible client to manage
content directly from your AI assistant.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := server.NewMCPServer(
			"content",
			Version,
			server.WithToolCapabilities(true),
		)

		addListContentTool(s)
		addGetStatusTool(s)
		addCheckContentTool(s)
		addGetTagsTool(s)
		addCreateContentTool(s)
		addPublishContentTool(s)
		addValidateConfigTool(s)

		return server.ServeStdio(s)
	},
}

func addListContentTool(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("list_content",
			mcp.WithDescription("List content files across collections. Returns slug, title, date, draft status, word count, reading time, and collection for each file."),
			mcp.WithString("collection", mcp.Description("Filter to a specific collection name (optional)")),
			mcp.WithBoolean("drafts_only", mcp.Description("Return only draft files")),
			mcp.WithBoolean("published_only", mcp.Description("Return only published files")),
			mcp.WithString("sort", mcp.Description("Sort field: date, title, or words (default: date)")),
		),
		handleListContent,
	)
}

func handleListContent(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	col := mcp.ParseString(req, "collection", "")
	draftsOnly := mcp.ParseBoolean(req, "drafts_only", false)
	publishedOnly := mcp.ParseBoolean(req, "published_only", false)
	sortField := mcp.ParseString(req, "sort", "date")

	files, err := scanCollections(col)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	type item struct {
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Date        string `json:"date"`
		Draft       bool   `json:"draft"`
		Words       int    `json:"words"`
		ReadingTime int    `json:"reading_time"`
		Collection  string `json:"collection"`
	}

	content.SortBy(files, sortField)

	result := make([]item, 0, len(files))
	for _, f := range files {
		draft := f.FM.GetBool("draft")
		if draftsOnly && !draft {
			continue
		}
		if publishedOnly && draft {
			continue
		}
		result = append(result, item{
			Slug:        f.Slug(),
			Title:       f.FM.GetString("title"),
			Date:        f.FM.GetString("date"),
			Draft:       draft,
			Words:       f.WordCount(),
			ReadingTime: f.ReadingTime(),
			Collection:  f.Collection,
		})
	}

	out, err := mcp.NewToolResultJSON(result)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return out, nil
}

func addGetStatusTool(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_status",
			mcp.WithDescription("Get aggregate content health stats: total files, published, drafts, word counts, and stale posts (not updated in 30+ days)."),
			mcp.WithString("collection", mcp.Description("Limit stats to a specific collection (optional)")),
		),
		handleGetStatus,
	)
}

func handleGetStatus(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	col := mcp.ParseString(req, "collection", "")

	files, err := scanCollections(col)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	type collectionStats struct {
		Total      int `json:"total"`
		Published  int `json:"published"`
		Drafts     int `json:"drafts"`
		TotalWords int `json:"total_words"`
		AvgWords   int `json:"avg_words"`
		Stale      int `json:"stale"`
	}

	byCollection := map[string]*collectionStats{}
	for _, f := range files {
		if _, ok := byCollection[f.Collection]; !ok {
			byCollection[f.Collection] = &collectionStats{}
		}
		s := byCollection[f.Collection]
		s.Total++
		if f.FM.GetBool("draft") {
			s.Drafts++
		} else {
			s.Published++
		}
		s.TotalWords += f.WordCount()
		if f.IsStale() {
			s.Stale++
		}
	}
	for _, s := range byCollection {
		if s.Total > 0 {
			s.AvgWords = s.TotalWords / s.Total
		}
	}

	totals := &collectionStats{}
	for _, s := range byCollection {
		totals.Total += s.Total
		totals.Published += s.Published
		totals.Drafts += s.Drafts
		totals.TotalWords += s.TotalWords
		totals.Stale += s.Stale
	}
	if totals.Total > 0 {
		totals.AvgWords = totals.TotalWords / totals.Total
	}

	result := map[string]any{
		"collections": byCollection,
		"totals":      totals,
	}
	out, err := mcp.NewToolResultJSON(result)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return out, nil
}

func addCheckContentTool(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("check_content",
			mcp.WithDescription("Validate content files for missing required fields, broken image references, and broken internal links."),
			mcp.WithString("collection", mcp.Description("Limit checks to a specific collection (optional)")),
		),
		handleCheckContent,
	)
}

func handleCheckContent(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	col := mcp.ParseString(req, "collection", "")

	files, err := scanCollections(col)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var allIssues []validate.Issue
	for _, f := range files {
		c := cfg.Collections[f.Collection]
		allIssues = append(allIssues, validate.RequiredFields(f, c.RequiredFields)...)
		allIssues = append(allIssues, validate.ImagesExist(f)...)
		allIssues = append(allIssues, validate.InternalLinks(f, c.Dir)...)
	}

	result := map[string]any{
		"ok":      len(allIssues) == 0,
		"checked": len(files),
		"issues":  allIssues,
	}
	out, err := mcp.NewToolResultJSON(result)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return out, nil
}

func addGetTagsTool(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_tags",
			mcp.WithDescription("List all tags found in frontmatter, sorted by frequency."),
			mcp.WithString("collection", mcp.Description("Filter to a specific collection (optional)")),
		),
		handleGetTags,
	)
}

func handleGetTags(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	col := mcp.ParseString(req, "collection", "")

	files, err := scanCollections(col)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	freq := map[string]int{}
	for _, f := range files {
		for _, tag := range extractTags(f.FM) {
			freq[tag]++
		}
	}

	type entry struct {
		Tag   string `json:"tag"`
		Count int    `json:"count"`
	}
	entries := make([]entry, 0, len(freq))
	for tag, count := range freq {
		entries = append(entries, entry{tag, count})
	}

	out, err := mcp.NewToolResultJSON(entries)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return out, nil
}

func addCreateContentTool(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("create_content",
			mcp.WithDescription("Scaffold a new content file in the specified collection. Returns the file path and slug of the created file."),
			mcp.WithString("collection",
				mcp.Description("The collection to create content in (required)"),
				mcp.Required(),
			),
			mcp.WithString("title",
				mcp.Description("Title of the new content (required)"),
				mcp.Required(),
			),
		),
		handleCreateContent,
	)
}

func handleCreateContent(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	colName := mcp.ParseString(req, "collection", "")
	title := mcp.ParseString(req, "title", "")

	if colName == "" {
		return mcp.NewToolResultError("collection is required"), nil
	}
	if title == "" {
		return mcp.NewToolResultError("title is required"), nil
	}

	col, ok := cfg.Collections[colName]
	if !ok {
		return mcp.NewToolResultError(fmt.Sprintf("unknown collection %q", colName)), nil
	}

	s := slug.From(title, col.Slug)
	filename := fmt.Sprintf("%s.%s", s, col.Format)
	destPath := filepath.Join(col.Dir, filename)

	fm := content.Frontmatter{
		"title": title,
		"date":  time.Now().Format("2006-01-02"),
		"draft": col.Defaults.Draft,
	}
	if col.Defaults.Author != "" {
		fm["author"] = col.Defaults.Author
	}
	for _, field := range col.RequiredFields {
		if _, exists := fm[field]; !exists {
			fm[field] = ""
		}
	}

	body := tmpl.Read(col.Template)

	if err := os.MkdirAll(col.Dir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("creating dir: %v", err)), nil
	}
	encoded := content.Encode(fm, body)
	if err := os.WriteFile(destPath, encoded, 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("writing file: %v", err)), nil
	}

	result := map[string]string{
		"path":       destPath,
		"slug":       s,
		"collection": colName,
	}
	out, err := mcp.NewToolResultJSON(result)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return out, nil
}

func addPublishContentTool(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("publish_content",
			mcp.WithDescription("Publish a draft content file by setting draft: false. Searches across all collections by slug."),
			mcp.WithString("slug",
				mcp.Description("The slug of the content file to publish (required)"),
				mcp.Required(),
			),
		),
		handlePublishContent,
	)
}

func handlePublishContent(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slugVal := mcp.ParseString(req, "slug", "")
	if slugVal == "" {
		return mcp.NewToolResultError("slug is required"), nil
	}

	files, err := scanCollections("")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	f, err := content.FindBySlug(files, slugVal)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if !f.FM.GetBool("draft") {
		return mcp.NewToolResultText(fmt.Sprintf("%s is already published", f.Path)), nil
	}

	f.FM["draft"] = false
	encoded := content.Encode(f.FM, f.Body)
	if err := os.WriteFile(f.Path, encoded, 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("writing %s: %v", f.Path, err)), nil
	}

	result := map[string]string{"path": f.Path, "slug": slugVal}
	out, err := mcp.NewToolResultJSON(result)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return out, nil
}

func addValidateConfigTool(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("validate_config",
			mcp.WithDescription("Validate the .content.yaml config file against the JSON schema. Returns any validation errors found."),
		),
		handleValidateConfig,
	)
}

func handleValidateConfig(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := os.ReadFile(".content.yaml")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("reading .content.yaml: %v", err)), nil
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("parsing .content.yaml: %v", err)), nil
	}

	errs, err := schema.ValidateConfig(raw)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("schema error: %v", err)), nil
	}

	result := map[string]any{
		"ok":     len(errs) == 0,
		"errors": errs,
	}
	out, err := mcp.NewToolResultJSON(result)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return out, nil
}
