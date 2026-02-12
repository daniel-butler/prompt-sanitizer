# Prompt Sanitizer

A Go CLI tool for wrapping untrusted external content with safety markers for LLM consumption. This helps protect against prompt injection attacks by clearly marking content from external sources.

## What It Does

Prompt Sanitizer wraps untrusted content with clear boundary markers that help LLMs distinguish between instructions and external data:

```
<<<EXTERNAL_UNTRUSTED_CONTENT>>>
Source: {source_label}
---
{content}
<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>
```

## Installation

### From Source

```bash
git clone https://github.com/openclaw/prompt-sanitizer
cd prompt-sanitizer
go build -o prompt-sanitizer ./cmd/prompt-sanitizer
sudo mv prompt-sanitizer /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/openclaw/prompt-sanitizer/cmd/prompt-sanitizer@latest
```

## Usage

### Wrap Content from Stdin

```bash
echo "untrusted data" | prompt-sanitizer --source "Web Search"
```

**Output:**
```
<<<EXTERNAL_UNTRUSTED_CONTENT>>>
Source: Web Search
---
untrusted data

<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>
```

### Wrap a File

```bash
prompt-sanitizer --source "email" --file message.txt
```

### Wrap Command Output

```bash
prompt-sanitizer --source "curl" -- curl https://example.com
```

The `--` separator is required before the command to prevent flag confusion.

## Examples

### Safe Web Scraping for LLM Context

```bash
prompt-sanitizer --source "Wikipedia API" -- curl -s "https://en.wikipedia.org/api/rest_v1/page/summary/Go_(programming_language)"
```

### Processing Email Content

```bash
cat email.txt | prompt-sanitizer --source "User Email" > safe-email.txt
```

### Wrapping Multiple Sources

```bash
# Search results
prompt-sanitizer --source "Search: Golang best practices" -- curl -s "https://api.search.example/q=golang"

# Documentation
prompt-sanitizer --source "Official Docs" --file golang-docs.md
```

## Edge Cases Handled

- **Empty input**: Properly wraps even when content is empty
- **Very long content**: Handles large files (tested with 1MB+)
- **Binary data**: Safely wraps binary content with null bytes
- **Special characters**: Preserves Unicode and special characters
- **Multiline content**: Maintains formatting and line breaks
- **Command failures**: Returns appropriate error codes

## Use Cases

1. **LLM Context Building**: Safely include web search results, API responses, or user-provided content in LLM prompts
2. **Email Processing**: Wrap email content before feeding to LLM-based email assistants
3. **Documentation Indexing**: Mark external documentation clearly when building RAG systems
4. **API Response Handling**: Wrap third-party API responses before LLM processing
5. **User Input Sanitization**: Clearly mark user-provided content in chatbot contexts

## Security Considerations

This tool provides **defense in depth** for prompt injection attacks:

- **Clear Boundaries**: Explicit markers help LLMs distinguish instructions from data
- **Source Attribution**: Labels help track where content originated
- **No Escape Issues**: The wrapper format is simple and doesn't require escaping

**Note:** This is a mitigation layer, not a complete solution. Always:
- Validate LLM outputs
- Use appropriate system prompts
- Apply principle of least privilege
- Monitor for unusual behavior

## Development

### Running Tests

```bash
go test ./...
```

### Running Tests with Verbose Output

```bash
go test ./... -v
```

### Building

```bash
go build -o prompt-sanitizer ./cmd/prompt-sanitizer
```

## Project Structure

```
prompt-sanitizer/
├── cmd/
│   └── prompt-sanitizer/    # Main CLI application
│       ├── main.go
│       └── main_test.go
├── pkg/
│   └── wrapper/              # Core wrapping logic
│       ├── wrapper.go
│       └── wrapper_test.go
├── go.mod
└── README.md
```

## Dependencies

None! This project uses only the Go standard library.

## License

MIT

## Contributing

Contributions welcome! Please:

1. Write tests first (TDD approach)
2. Keep commits small and focused
3. Follow Go best practices
4. Update README for new features

## FAQ

**Q: Why not use JSON or XML for wrapping?**  
A: Simple text markers are easier for LLMs to recognize and harder to accidentally escape.

**Q: Can this prevent all prompt injection attacks?**  
A: No security measure is perfect. This adds a layer of defense but should be combined with other security practices.

**Q: What if the content already contains the markers?**  
A: The LLM should still be able to distinguish the outermost wrapper as the boundary. For extra safety, you could pre-process content to escape existing markers.

**Q: Does this work with all LLMs?**  
A: The effectiveness depends on the LLM's training and how you phrase your system prompt. Most modern LLMs can be instructed to respect these boundaries.

## Acknowledgments

Built with test-driven development (TDD) following Go best practices.
