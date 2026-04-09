package common

import "unicode"

// RefKind describes the kind of source entry a ref points at.
type RefKind string

const (
	RefKindPackage      RefKind = "package"
	RefKindFile         RefKind = "file"
	RefKindType         RefKind = "type"
	RefKindFunction     RefKind = "function"
	RefKindVariable     RefKind = "variable"
	RefKindFunctionCall RefKind = "functioncall"
	RefKindDependency   RefKind = "dependency"
)

// Ref identifies a source location for a matched entry.
type Ref struct {
	PackageID   string
	PackageName string
	Filename    string
	Line        int
	Column      int
	Kind        RefKind
	Match       string
}

// IsExportedName reports whether name is an exported Go identifier
// (starts with an uppercase Unicode letter).
func IsExportedName(name string) bool {
	if len(name) == 0 {
		return false
	}
	return unicode.IsUpper([]rune(name)[0])
}
