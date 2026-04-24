# AI Instructions for `content` CLI

You are an AI assistant or agent with access to the `content` CLI tool via the Model Context Protocol (MCP). Your goal is to help the user manage their Markdown/MDX content collections efficiently.

## Core Concepts

- **Collections:** Content is organized into named collections (e.g., `blog`, `testimonials`). Each collection has its own directory, file format, and frontmatter requirements.
- **Frontmatter:** Files use YAML frontmatter. Common fields include `title`, `date`, `draft`, and `tags`.
- **Slugs:** Files are identified and manipulated using slugs (kebab-case by default).

## Interaction Guidelines

### 1. Content Discovery
Use `list_content` to see what exists. You can filter by collection, draft status, or sort by date/title/words.
- *Example:* "Show me my latest draft blog posts" -> `list_content(collection="blog", drafts_only=true, sort="date")`

### 2. Health & Validation
Use `get_status` for a high-level overview and `check_content` to find issues like missing fields or broken links.
- Always run `check_content` after creating or modifying files to ensure they meet the project's schema.

### 3. Workflow Patterns
- **Creating:** Use `create_content` with a title and collection. This scaffolds the file with required frontmatter.
- **Publishing:** Use `publish_content` with a slug. This automatically flips `draft: false`.
- **Tagging:** Use `get_tags` to see existing tags before suggesting or adding new ones to maintain consistency.

## Tool-Specific Tips

- **`list_content`**: Returns metadata only. To read actual content, you would typically need a filesystem tool, but `content` provides the structural metadata needed for management.
- **`create_content`**: If a collection requires specific fields not covered by the tool parameters, you may need to advise the user to edit the file manually or use a filesystem tool if available.
- **`validate_config`**: If the user is having trouble with commands, check if their `.content.yaml` is valid first.

## Best Practices
- Prefer kebab-case for slugs.
- Always check if a slug exists using `list_content` before attempting to create content with a similar title to avoid confusion.
- When listing content, mention the "Reading Time" and "Word Count" to give the user context on the length of their posts.
