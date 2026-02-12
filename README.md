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

### From Releases

Download the latest binary from [Releases](https://github.com/daniel-butler/prompt-sanitizer/releases).

### From Source

```bash
git clone https://github.com/daniel-butler/prompt-sanitizer
cd prompt-sanitizer
go build -o prompt-sanitizer ./cmd/prompt-sanitizer
sudo mv prompt-sanitizer /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/daniel-butler/prompt-sanitizer/cmd/prompt-sanitizer@latest
```

## Usage

### Wrap Content from Stdin

```bash
echo "untrusted data" | prompt-sanitizer --source "Web Search"
```

### Wrap a File

```bash
prompt-sanitizer --source "email" --file message.txt
```

### Wrap Command Output

```bash
prompt-sanitizer --source "curl" -- curl https://example.com
```

### Check Version

```bash
prompt-sanitizer --version
```

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

## Research & References

This tool implements the "spotlighting" defense pattern for prompt injection, informed by the following research:

### General Prompt Injection

- [OWASP LLM Top 10](https://owasp.org/www-project-top-10-for-large-language-model-applications/) - Industry standard for LLM security risks
- [OWASP Prompt Injection Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/LLM_Prompt_Injection_Prevention_Cheat_Sheet.html) - Prevention techniques
- [Simon Willison's Prompt Injection Explainer](https://simonwillison.net/2023/Apr/14/worst-that-can-happen/) - Excellent introduction to the threat
- [Greshake et al. "Not what you've signed up for"](https://arxiv.org/abs/2302.12173) - Academic paper on indirect prompt injection

### Defense Techniques

- [Spotlighting (Microsoft)](https://arxiv.org/abs/2403.14720) - The technique this tool implements: marking external data with delimiters
- [Instruction Hierarchy](https://arxiv.org/abs/2404.13208) - Training LLMs to prioritize system instructions
- [Anthropic Prompt Injection Prevention](https://www.anthropic.com/research/prompt-injection-prevention) - Structural separation guidelines

### Attack Patterns (Informing Test Suite)

- [Perez & Ribeiro "Ignore This Title and HackAPrompt"](https://arxiv.org/abs/2311.16119) - Competition findings on injection attacks
- [HackAPrompt Competition](https://www.aicrowd.com/challenges/hackaprompt-2023) - Real-world attack dataset
- [Anthropic Many-Shot Jailbreaking](https://www.anthropic.com/research/many-shot-jailbreaking) - Analysis of jailbreak techniques
- [ChatGPT System Prompt Extraction](https://github.com/LouisShark/chatgpt_system_prompt) - Prompt leakage techniques
- [Indirect Prompt Injection via PDF](https://kai-greshake.de/posts/inject-my-pdf/) - Document-based injection
- [Data Exfiltration via Email](https://embracethered.com/blog/posts/2023/google-bard-data-exfiltration/) - Tool use exploitation

### Obfuscation & Evasion

- [Typoglycemia in NLP](https://arxiv.org/abs/1905.11268) - Character perturbation attacks
- [AI Prompt Injection Techniques](https://www.mdsec.co.uk/2023/04/ai-prompt-injection/) - Obfuscation methods
- [LLM Agent Prompt Injection](https://labs.withsecure.com/publications/llm-agent-prompt-injection) - Encoded payload attacks

## Test Coverage

The test suite includes 160+ tests covering:

- **Marker manipulation attacks** - Null bytes, zero-width chars, BOM injection
- **Unicode confusion** - Cyrillic/Greek homoglyphs, RTL override, bidi attacks
- **Control characters** - ANSI escapes, terminal injection
- **Real LLM jailbreaks** - DAN, developer mode, instruction override
- **Resource exhaustion** - Large inputs, zalgo text
- **Fuzzing** - Randomized input testing

Run tests:
```bash
go test -v ./...
go test -race ./...           # Race detection
go test -fuzz=Fuzz ./pkg/wrapper/  # Fuzzing
```

## Development

### Running Tests

```bash
go test ./...
```

### Building with Version

```bash
VERSION=2024.01.15-1
go build -ldflags="-X main.Version=${VERSION}" -o prompt-sanitizer ./cmd/prompt-sanitizer
```

## Project Structure

```
prompt-sanitizer/
├── .github/
│   └── workflows/
│       └── release.yml       # CI/CD with auto-releases
├── cmd/
│   └── prompt-sanitizer/
│       ├── main.go
│       └── main_test.go
├── pkg/
│   └── wrapper/
│       ├── wrapper.go
│       ├── wrapper_test.go
│       └── adversarial_test.go  # Security-focused tests
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
2. Include adversarial test cases for security features
3. Keep commits small and focused
4. Update README for new features

## Acknowledgments

- Built following security research from OWASP, Microsoft, Anthropic, and the academic community
- Inspired by the need to safely integrate untrusted content with LLM agents
- Test suite informed by real-world prompt injection competitions and research
