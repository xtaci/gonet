package dos

const (
	RED = iota
	BLACK
)

type Node struct {
	left   *Node
	right  *Node
	parent *Node

	score int // the score
	size  int // the size of this subtree
	color int

	data interface{} // associated data
}

func (n *Node) Data() interface{} {
	return n.data
}

func (n *Node) Score() int {
	return n.score
}

//
type Tree struct {
	root *Node
}

func (t *Tree) Clear() {
	t.root = nil
}

func (t *Tree) Root() *Node {
	return t.root
}

//--------------------------------------------------------- Dos Part
func _nodesize(n *Node) int {
	if n == nil {
		return 0
	}

	return n.size
}

func lookup_node(n *Node, rank int) *Node {
	if n == nil {
		return nil // beware of nil pointer
	}

	size := _nodesize(n.left) + 1

	if rank == size {
		return n
	}

	if rank < size {
		return lookup_node(n.left, rank)
	}
	return lookup_node(n.right, rank-size)
}

func new_node(score int, data interface{}, color int, left, right *Node) *Node {
	n := Node{score: score, color: color, left: left, right: right, size: 1, data: data}
	return &n
}

//--------------------------------------------------------- Lookup by Rank
func (t *Tree) Rank(rank int) *Node {
	return lookup_node(t.root, rank)
}

//--------------------------------------------------------- Lookup by Rank
func (t *Tree) Count() int {
	if t.root != nil {
		return t.root.size
	}

	return 0
}

//--------------------------------------------------------- Lookup by score
func (t *Tree) ByScore(score int) (n *Node, rank int) {
	n = t.root

	if n == nil {
		return
	}

	rank = 0
	for n != nil {
		if score == n.score {
			rank += _nodesize(n.left) + 1
			return
		} else if score > n.score {
			n = n.left
		} else {
			rank += _nodesize(n.left) + 1
			n = n.right
		}
	}

	return
}

func (t *Tree) Insert(score int, data interface{}) {
	inserted_node := new_node(score, data, RED, nil, nil)
	if t.root == nil {
		t.root = inserted_node
	} else {
		n := t.root
		for {
			n.size++
			if score > n.score {
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
	fixup_size(n)

	// phase 2. handle red-black properties, and deletion work.
	if n.left != nil && n.right != nil {
		/* Copy fields from predecessor and then delete it instead */
		pred := maximum_node(n.left)
		n.score = pred.score
		n.size = pred.size
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

/**
 * left/right rotation call back function
 */
func rotate_left_callback(n, parent *Node) {
	parent.size = _nodesize(n)
	n.size = _nodesize(n.left) + _nodesize(n.right) + 1
}

func rotate_right_callback(n, parent *Node) {
	rotate_left_callback(n, parent)
}

func fixup_size(n *Node) {
	for n != nil {
		n.size--
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
