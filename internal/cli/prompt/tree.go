package prompt

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var ErrAborted = errors.New("aborted")

type LeafMeta struct {
	Value       string
	Description string
}

type treeNode struct {
	name     string
	children map[string]*treeNode
	leafKey  string
}

func (n *treeNode) isLeaf() bool { return n.leafKey != "" }

func (n *treeNode) childNames() []string {
	names := make([]string, 0, len(n.children))
	for name := range n.children {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func buildKeyTree(keys []string) *treeNode {
	root := &treeNode{children: make(map[string]*treeNode)}
	for _, key := range keys {
		cur := root
		parts := strings.Split(key, ".")
		for i, part := range parts {
			child, ok := cur.children[part]
			if !ok {
				child = &treeNode{name: part, children: make(map[string]*treeNode)}
				cur.children[part] = child
			}
			cur = child
			if i == len(parts)-1 {
				cur.leafKey = key
			}
		}
	}
	return root
}

var dimValue = lipgloss.NewStyle().Faint(true)

func leafLabel(name string, meta LeafMeta) string {
	value := meta.Value
	if value == "" {
		value = "(empty)"
	}
	line := dimValue.Render(value)
	if meta.Description != "" {
		line += dimValue.Render("  —  " + meta.Description)
	}
	return name + "\n  " + line
}

func renderLevel(cur *treeNode, vals map[string]LeafMeta) (names []string, labels []string) {
	names = cur.childNames()
	labels = make([]string, 0, len(names)+1)
	labels = append(labels, LabelBack)
	for _, name := range names {
		child := cur.children[name]
		if !child.isLeaf() {
			labels = append(labels, name+"/")
			continue
		}
		labels = append(labels, leafLabel(name, vals[child.leafKey]))
	}
	return names, labels
}

func BrowseKeys(
	title string,
	keys []string,
	values func() map[string]LeafMeta,
	onLeaf func(key string) error,
	startPrefix string,
	focusKey string,
) error {
	root := buildKeyTree(keys)
	stack := stackAt(root, startPrefix)
	focus := leafFocus(stack[len(stack)-1], focusKey)

	for {
		cur := stack[len(stack)-1]
		vals := values()
		if vals == nil {
			vals = map[string]LeafMeta{}
		}
		names, labels := renderLevel(cur, vals)

		idx, err := Select(title, treePath(stack), labels, focus)
		if err != nil {
			return err
		}
		if idx == 0 {
			if len(stack) == 1 {
				return ErrAborted
			}
			stack = stack[:len(stack)-1]
			focus = 0
			continue
		}

		chosen := cur.children[names[idx-1]]
		if !chosen.isLeaf() {
			stack = append(stack, chosen)
			focus = 0
			continue
		}
		if err := onLeaf(chosen.leafKey); err != nil {
			return err
		}
		focus = idx
	}
}

func TreeKey(title string, keys []string, values map[string]LeafMeta, out *string) error {
	if out == nil {
		return fmt.Errorf("key value is nil")
	}
	root := buildKeyTree(keys)
	stack := []*treeNode{root}

	for {
		cur := stack[len(stack)-1]
		names, labels := renderLevel(cur, values)

		idx, err := Select(title, treePath(stack), labels, 0)
		if err != nil {
			return err
		}
		if idx == 0 {
			if len(stack) == 1 {
				return ErrAborted
			}
			stack = stack[:len(stack)-1]
			continue
		}

		chosen := cur.children[names[idx-1]]
		if !chosen.isLeaf() {
			stack = append(stack, chosen)
			continue
		}
		*out = chosen.leafKey
		return nil
	}
}

func EnsureTreeKey(title string, keys []string, values map[string]LeafMeta, out *string) error {
	if out == nil {
		return fmt.Errorf("key value is nil")
	}
	if *out != "" {
		if !slices.Contains(keys, *out) {
			return fmt.Errorf("unknown key %q", *out)
		}
		return nil
	}
	return TreeKey(title, keys, values, out)
}

func ParentPrefix(key string) string {
	i := strings.LastIndex(key, ".")
	if i <= 0 {
		return ""
	}
	return key[:i]
}

func leafFocus(cur *treeNode, key string) int {
	if key == "" {
		return 0
	}
	for i, name := range cur.childNames() {
		child := cur.children[name]
		if child.isLeaf() && child.leafKey == key {
			return i + 1
		}
	}
	return 0
}

func stackAt(root *treeNode, prefix string) []*treeNode {
	stack := []*treeNode{root}
	if prefix == "" {
		return stack
	}
	cur := root
	for _, part := range strings.Split(prefix, ".") {
		child, ok := cur.children[part]
		if !ok || child.isLeaf() {
			break
		}
		stack = append(stack, child)
		cur = child
	}
	return stack
}

func treePath(stack []*treeNode) string {
	if len(stack) <= 1 {
		return "↑↓ move · enter · back to leave"
	}
	parts := make([]string, 0, len(stack)-1)
	for _, n := range stack[1:] {
		parts = append(parts, n.name)
	}
	return strings.Join(parts, ".") + "  ·  ↑↓ move · enter"
}
