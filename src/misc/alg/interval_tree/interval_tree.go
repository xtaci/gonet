package interval_tree

import (
	"math"
)

const (
	RED = iota
	BLACK
)

type Node struct {
	left   *Node
	right  *Node
	parent *Node

	low   int64 // lower-bound
	high  int64 // higher-bound
	m     int64 // max subtree upper bound
	color int

	data interface{} // associated data
}

func (n *Node) Data() interface{} {
	return n.data
}

//
type Tree struct {
	root *Node
}

func new_node(low, high int64, data interface{}, color int, left, right *Node) *Node {
	n := Node{low: low, high: high, m: high, color: color, left: left, right: right, data: data}
	return &n
}

/**
 * interval tree lookup function
 *
 * search range [low, high] for overlap, return only one element
 * use lookup & delete & insert schema to get multiple elements
 *
 * nil is returned if not found.
 */
func (t *Tree) Lookup(low, high int64) *Node {
	n := t.root
	for n != nil && (low > n.high || n.low > high) { // should search in childs
		if n.left != nil && low <= n.left.m {
			n = n.left // path choice on m.
		} else {
			n = n.right
		}
	}

	return n
}

/**
 * Insert
 * insert range [low, high] into red-black tree
 */
func (t *Tree) Insert(low, high int64, data interface{}) {
	inserted_node := new_node(low, high, data, RED, nil, nil)
	if t.root == nil {
		t.root = inserted_node
	} else {
		n := t.root
		for {
			// update 'm' for each node traversed from root
			if inserted_node.m > n.m {
				n.m = inserted_node.m
			}

			// find a proper position
			if low < n.low {
				if n.left == nil {
					n.left = inserted_node
					break
				} else {
					n = n.left
				}
			} else {
				if n.right == nil {
					n.right = inserted_node
					break
				} else {
					n = n.right
				}
			}
		}
		inserted_node.parent = n
	}

	t.insert_case1(inserted_node)
}

func (t *Tree) DeleteNode(n *Node) {
	// phase 1. fix up size
	fixup_m(n)

	// phase 2. handle red-black properties, and deletion work.
	if n.left != nil && n.right != nil {
		/* Copy fields from predecessor and then delete it instead */
		pred := maximum_node(n.left)
		n.low = pred.low
		n.high = pred.high
		n.data = pred.data
		n = pred
	}

	var child *Node
	if n.right == nil {
		child = n.left
	} else {
		child = n.right
	}

	if node_color(n) == BLACK {
		n.color = node_color(child)
		t.delete_case1(n)
	}

	t.replace_node(n, child)

	if n.parent == nil && child != nil {
		child.color = BLACK
	}
}

func Max(a, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

func M(n *Node) int64 {
	if n == nil {
		return int64(math.MinInt64)
	}

	return n.m
}

/**
 * fix 'm' value caused by rotation
 */
func rotate_left_callback(n, parent *Node) {
	// parent inherit max m value
	parent.m = n.m
	// update node 'm' value by it's children.
	n.m = Max(n.high, Max(M(n.left), M(n.right)))
}

func rotate_right_callback(n, parent *Node) {
	rotate_left_callback(n, parent)
}

/**
 * fix up 'm' value caued by deletion
 */
func fixup_m(n *Node) {
	m := M(n)
	m_new := Max(M(n.left), M(n.right))

	// if current 'm' is not decided by n, just return.
	if m == m_new {
		return
	}

	for n.parent != nil {
		n.parent.m = Max(n.parent.high, Max(m_new, M(sibling(n))))

		if M(n.parent) > m {
			break // since node n does not affect
			//  the result anymore, we break.
		}
		n = n.parent
	}
}

//--------------------------------------------------------- Tree part
func grandparent(n *Node) *Node {
	return n.parent.parent
}

func sibling(n *Node) *Node {
	if n == n.parent.left {
		return n.parent.right
	}
	return n.parent.left
}

func uncle(n *Node) *Node {
	return sibling(n.parent)
}

func node_color(n *Node) int {
	if n == nil {
		return BLACK
	}
	return n.color
}

func (t *Tree) rotate_left(n *Node) {
	r := n.right
	t.replace_node(n, r)
	n.right = r.left
	if r.left != nil {
		r.left.parent = n
	}
	r.left = n
	n.parent = r

	rotate_left_callback(n, r)
}

func (t *Tree) rotate_right(n *Node) {
	L := n.left
	t.replace_node(n, L)
	n.left = L.right
	if L.right != nil {
		L.right.parent = n
	}
	L.right = n
	n.parent = L

	rotate_right_callback(n, L)
}

func (t *Tree) replace_node(oldn, newn *Node) {
	if oldn.parent == nil {
		t.root = newn
	} else {
		if oldn == oldn.parent.left {
			oldn.parent.left = newn
		} else {
			oldn.parent.right = newn
		}
	}
	if newn != nil {
		newn.parent = oldn.parent
	}
}

func (t *Tree) insert_case1(n *Node) {
	if n.parent == nil {
		n.color = BLACK
	} else {
		t.insert_case2(n)
	}
}

func (t *Tree) insert_case2(n *Node) {
	if node_color(n.parent) == BLACK {
		return /* Tree is still valid */
	} else {
		t.insert_case3(n)
	}
}

func (t *Tree) insert_case3(n *Node) {
	if node_color(uncle(n)) == RED {
		n.parent.color = BLACK
		uncle(n).color = BLACK
		grandparent(n).color = RED
		t.insert_case1(grandparent(n))
	} else {
		t.insert_case4(n)
	}
}

func (t *Tree) insert_case4(n *Node) {
	if n == n.parent.right && n.parent == grandparent(n).left {
		t.rotate_left(n.parent)
		n = n.left
	} else if n == n.parent.left && n.parent == grandparent(n).right {
		t.rotate_right(n.parent)
		n = n.right
	}
	t.insert_case5(n)
}

func (t *Tree) insert_case5(n *Node) {
	n.parent.color = BLACK
	grandparent(n).color = RED
	if n == n.parent.left && n.parent == grandparent(n).left {
		t.rotate_right(grandparent(n))
	} else {
		t.rotate_left(grandparent(n))
	}
}

func maximum_node(n *Node) *Node {
	for n.right != nil {
		n = n.right
	}
	return n
}

func (t *Tree) delete_case1(n *Node) {
	if n.parent == nil {
		return
	} else {
		t.delete_case2(n)
	}
}

func (t *Tree) delete_case2(n *Node) {
	if node_color(sibling(n)) == RED {
		n.parent.color = RED
		sibling(n).color = BLACK
		if n == n.parent.left {
			t.rotate_left(n.parent)
		} else {
			t.rotate_right(n.parent)
		}
	}
	t.delete_case3(n)
}

func (t *Tree) delete_case3(n *Node) {
	if node_color(n.parent) == BLACK &&
		node_color(sibling(n)) == BLACK &&
		node_color(sibling(n).left) == BLACK &&
		node_color(sibling(n).right) == BLACK {
		sibling(n).color = RED
		t.delete_case1(n.parent)
	} else {
		t.delete_case4(n)
	}
}

func (t *Tree) delete_case4(n *Node) {
	if node_color(n.parent) == RED &&
		node_color(sibling(n)) == BLACK &&
		node_color(sibling(n).left) == BLACK &&
		node_color(sibling(n).right) == BLACK {
		sibling(n).color = RED
		n.parent.color = BLACK
	} else {
		t.delete_case5(n)
	}
}

func (t *Tree) delete_case5(n *Node) {
	if n == n.parent.left &&
		node_color(sibling(n)) == BLACK &&
		node_color(sibling(n).left) == RED &&
		node_color(sibling(n).right) == BLACK {
		sibling(n).color = RED
		sibling(n).left.color = BLACK
		t.rotate_right(sibling(n))
	} else if n == n.parent.right &&
		node_color(sibling(n)) == BLACK &&
		node_color(sibling(n).right) == RED &&
		node_color(sibling(n).left) == BLACK {
		sibling(n).color = RED
		sibling(n).right.color = BLACK
		t.rotate_left(sibling(n))
	}
	t.delete_case6(n)
}

func (t *Tree) delete_case6(n *Node) {
	sibling(n).color = node_color(n.parent)
	n.parent.color = BLACK
	if n == n.parent.left {
		sibling(n).right.color = BLACK
		t.rotate_left(n.parent)
	} else {
		sibling(n).left.color = BLACK
		t.rotate_right(n.parent)
	}
}
