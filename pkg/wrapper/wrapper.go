package wrapper

import "fmt"

// WrapContent wraps untrusted content with safety markers for LLM consumption
func WrapContent(content, source string) string {
	return fmt.Sprintf(`<<<EXTERNAL_UNTRUSTED_CONTENT>>>
Source: %s
---
%s
<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>`, source, content)
}
