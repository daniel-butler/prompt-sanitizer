package wrapper

import (
	"testing"
)

func TestWrapContent_Basic(t *testing.T) {
	content := "Hello, world!"
	source := "Test Source"
	
	result := WrapContent(content, source)
	
	expected := `<<<EXTERNAL_UNTRUSTED_CONTENT>>>
Source: Test Source
---
Hello, world!
<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>`
	
	if result != expected {
		t.Errorf("WrapContent() mismatch\nGot:\n%s\n\nExpected:\n%s", result, expected)
	}
}

func TestWrapContent_EmptyContent(t *testing.T) {
	content := ""
	source := "Empty Source"
	
	result := WrapContent(content, source)
	
	expected := `<<<EXTERNAL_UNTRUSTED_CONTENT>>>
Source: Empty Source
---

<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>`
	
	if result != expected {
		t.Errorf("WrapContent() with empty content mismatch\nGot:\n%s\n\nExpected:\n%s", result, expected)
	}
}

func TestWrapContent_MultilineContent(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3"
	source := "Multiline Source"
	
	result := WrapContent(content, source)
	
	expected := `<<<EXTERNAL_UNTRUSTED_CONTENT>>>
Source: Multiline Source
---
Line 1
Line 2
Line 3
<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>`
	
	if result != expected {
		t.Errorf("WrapContent() with multiline content mismatch\nGot:\n%s\n\nExpected:\n%s", result, expected)
	}
}

func TestWrapContent_SpecialCharacters(t *testing.T) {
	content := "Content with <script>alert('xss')</script> and 特殊字符"
	source := "Web Scrape"
	
	result := WrapContent(content, source)
	
	expected := `<<<EXTERNAL_UNTRUSTED_CONTENT>>>
Source: Web Scrape
---
Content with <script>alert('xss')</script> and 特殊字符
<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>`
	
	if result != expected {
		t.Errorf("WrapContent() with special characters mismatch\nGot:\n%s\n\nExpected:\n%s", result, expected)
	}
}
