package common

// Ref identifies a source location for a matched entry.
type Ref struct {
	PackageID   string
	PackageName string
	Filename    string
	Line        int
	Column      int
}
