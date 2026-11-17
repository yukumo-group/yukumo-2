package port

// ArtifactStore holds synthesised output and its cache index.
//
// Path is display-only: a filesystem adapter returns an absolute path, while
// memory and JS adapters return the logical key (or an empty string). Path's
// result must never be passed back as a key.
type ArtifactStore interface {
	Store
	Path(key string) string
}

// OutputWriter writes an explicit user-supplied path (-o). Hosts that have no
// filesystem (browser) leave it nil.
type OutputWriter interface {
	Write(path string, data []byte) error
}
