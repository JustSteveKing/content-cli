# content

A CLI tool for managing Markdown and MDX content files. Designed for projects like Astro sites where content lives in structured directories with YAML frontmatter.

## Installation

Download a pre-built binary from [GitHub Releases](https://github.com/juststeveking/content-cli/releases), or install with Go:

```bash
go install github.com/juststeveking/content-cli@latest
```

### Nix

If you use the Nix package manager, you can install `content` directly from this repository:

```bash
nix profile install github:juststeveking/content-cli
```

Or run it without installing:

```bash
nix run github:juststeveking/content-cli -- --help
```

### NixOS

To include `content` in your NixOS configuration, add this repository as a flake input:

```nix
# flake.nix
{
  inputs.content-cli.url = "github:juststeveking/content-cli";

  outputs = { self, nixpkgs, content-cli }: {
    nixosConfigurations.my-machine = nixpkgs.lib.nixosSystem {
      modules = [
        ({ pkgs, ... }: {
          environment.systemPackages = [
            content-cli.packages.${pkgs.system}.default
          ];
        })
      ];
    };
  };
}
```

## Quick start

```bash
# Initialise a config file
content init

# Create a new post
content new --collection blog "My First Post"

# List all content
content list

# Publish a draft
content publish my-first-post
```

## Configuration

Running `content init` creates a `.content.yaml` file in the current directory. You can define multiple collections — one per content type.

```yaml
default_collection: blog
collections:
  blog:
    dir: src/content/blog
    format: mdx
    required_fields: [title, date, draft]
    optional_fields: [description, tags]
    slug: kebab
    defaults:
      draft: true
      author: ""
  testimonials:
    dir: src/content/testimonials
    format: mdx
    required_fields: [name, company, quote]
    optional_fields: [role, avatar]
    slug: kebab
    defaults:
      draft: false
```

| Field | Description |
|---|---|
| `dir` | Directory where content files are stored |
| `format` | File format: `md` or `mdx` |
| `template` | Path to a template file used when scaffolding new content |
| `required_fields` | Frontmatter fields every file must have |
| `optional_fields` | Frontmatter fields that are optional |
| `slug` | Slug style for filenames: `kebab`, `snake`, or `raw` |
| `defaults.draft` | Whether new files are drafts by default |
| `defaults.author` | Default author value inserted into frontmatter |

Validate the config at any time:

```bash
content config
```

## Commands

### `content init`

Interactive wizard that creates `.content.yaml`. Supports defining multiple collections.

### `content new --collection <name> <title>`

Scaffold a new content file in the specified collection. The filename is derived from the title using the collection's slug style. The `date` field is pre-filled with today's date.

```bash
content new --collection blog "Getting Started with Astro"
# → src/content/blog/getting-started-with-astro.mdx
```

### `content list`

List all content files across all collections, with slug, title, date, draft status, word count, and reading time.

```bash
content list                        # all collections
content list --collection blog      # one collection
content list --drafts               # drafts only
content list --published            # published only
content list --sort title           # sort by title (date|title|words)
```

### `content publish <slug>`

Set `draft: false` on a content file, identified by slug. Searches across all collections.

```bash
content publish getting-started-with-astro
```

### `content edit <slug>`

Open a content file in `$EDITOR`. Falls back to `$VISUAL`, then `vi`.

```bash
content edit getting-started-with-astro
```

### `content check`

Validate content files for missing required fields, broken image references, and broken internal links. Exits 1 if any issues are found.

```bash
content check                       # all collections
content check --collection blog     # one collection
```

### `content status`

Show an aggregate health summary — total files, published, drafts, word counts, stale posts.

```bash
content status
content status --collection blog
```

### `content tags`

List all tags found in frontmatter, sorted by frequency.

```bash
content tags
content tags --collection blog
```

### `content config`

Validate `.content.yaml` against the JSON schema. Exits 1 if the config has violations.

## AI & Agent Integration (MCP)

`content` includes a built-in [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server. This allows AI agents (like Claude Desktop) to interact with your content files directly, enabling them to list, create, validate, and publish content based on your instructions.

### Starting the MCP Server

```bash
content serve
```

The server communicates over `stdio`, making it easy to integrate with local AI clients.

### Claude Desktop Integration

To use `content` with Claude Desktop, add it to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "content": {
      "command": "content",
      "args": ["serve"]
    }
  }
}
```

### Available Tools

Once connected, your AI assistant will have access to the following tools:

- `list_content`: List content across collections (with filtering and sorting).
- `get_status`: Get aggregate health stats (published vs draft, word counts, stale posts).
- `check_content`: Validate files for missing fields and broken links.
- `get_tags`: List and count tags used across your content.
- `create_content`: Scaffold a new content file with proper frontmatter.
- `publish_content`: Set `draft: false` on a file by its slug.
- `validate_config`: Verify your `.content.yaml` is valid.

## Development

```bash
make build      # build ./content
make test       # run tests
make check      # fmt + vet
make install    # install to $GOPATH/bin
make snapshot   # build all release targets locally (requires goreleaser)
```

## Releasing

Push a version tag to trigger the release workflow:

```bash
git tag v0.1.0
git push origin v0.1.0
```

GoReleaser will build binaries for Linux, macOS (Intel + Apple Silicon), and Windows, then publish them to GitHub Releases.

## License

MIT — see [LICENSE](LICENSE).
