// Package menu defines the transport-neutral browse menu model.
package menu

import "yukumo-script/cli/streams"

// Node is a directory when Children is non-empty, otherwise a runnable leaf.
type Node struct {
	Name     string
	Short    string
	Children []Node
	Run      func(s streams.Streams) error
}
