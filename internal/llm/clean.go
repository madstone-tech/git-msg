package llm

import (
	"regexp"
	"strings"
)

var (
	// fenceRegex matches markdown code fences and captures the content inside.
	fenceRegex = regexp.MustCompile("(?s)```(?:[a-z]*\\n)?(.*?)\\n?```")

	// fillerRegex matches common conversational prefixes LLMs use.
	fillerRegex = regexp.MustCompile(`(?i)^(?:here is (?:the|your)? (?:generated )?commit message|suggested commit message|the commit message is|commit message|generated message):?\s*`)
)

// CleanResponse removes common LLM conversational filler and markdown fencing
// to extract the actual commit message.
func CleanResponse(input string) string {
	output := strings.TrimSpace(input)

	// 1. Remove markdown fences if present.
	if matches := fenceRegex.FindStringSubmatch(output); len(matches) > 1 {
		output = matches[1]
	}

	// 2. Remove common conversational fillers.
	output = fillerRegex.ReplaceAllString(output, "")

	return strings.TrimSpace(output)
}
