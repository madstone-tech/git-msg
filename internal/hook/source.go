package hook

// SourceType represents the $2 argument passed to prepare-commit-msg.
type SourceType string

const (
	SourceNormal   SourceType = ""         // plain git commit
	SourceMessage  SourceType = "message"  // git commit -m "..."
	SourceTemplate SourceType = "template" // commit.template used
	SourceMerge    SourceType = "merge"    // merge commit
	SourceSquash   SourceType = "squash"   // squash commit
	SourceCommit   SourceType = "commit"   // amend
)

// ShouldGenerate returns true only for source types where generation is implemented.
// All other sources exit 0 without generating.
func ShouldGenerate(s SourceType) bool {
	switch s {
	case SourceNormal, SourceMessage:
		return true
	default:
		return false
	}
}
