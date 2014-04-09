package dos

import (
	"fmt"
)

const (
	RED = iota
	BLACK
)

type Node struct {
	left   *Node
	right  *Node
	parent *Node

	size  int // the size of this subtree
	color int

	score int32 // the score
	id    int32 // associated id
}

func (n *Node) Id() int32 {
	return n.id
}

func (n *Node) Score() int32 {
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

func new_node(score int32, id int32, color int, left, right *Node) *Node {
	n := Node{score: score, color: color, left: left, right: right, size: 1, id: id}
	return &n
}

//--------------------------------------------------------- Lookup by Rank
// READ-LOCK
func (t *Tree) Rank(rank int) *Node {
	return lookup_node(t.root, rank)
}

//--------------------------------------------------------- Lookup by Rank
// READ-LOCK
func (t *Tree) Count() int {
	if t.root != nil {
		return _nodesize(t.root.left) + _nodesize(t.root.right) + 1
	}

	return 0
}

//--------------------------------------------------------- Lookup by score
func (t *Tree) _lookup_score(score int32) (rank int, n *Node) {
	n = t.root

	if n == nil {
		return -1, nil
	}

	base := 0
	for n != nil {
		if score == n.score {
			rank = base + _nodesize(n.left) + 1
			return rank, n
		} else if score > n.score {
			n = n.left
		} else {
			base += _nodesize(n.left) + 1
			n = n.right
		}
	}

	return -1, nil
}

//---------------------------------------------------------- locate a score & id
// WRITE-LOCK
func (t *Tree) Locate(score int32, id int32) (int, *Node) {
	tmplist := make([]int32, 0, 64)

	defer func() {
		for i := range tmplist {
			t.Insert(score, tmplist[i])
		}
	}()

	for {
		rank, node := t._lookup_score(score)

		if node == nil { // no such score exists
			return -1, nil
		}

		if node.id == id { // found matched id
			return rank, node
		} else {
			// temporary delete
			tmplist = append(tmplist, node.id)
			t.DeleteNode(node)
		}
	}
}

//---------------------------------------------------------- Insert an element
// WRITE-LOCK
func (t *Tree) Insert(score int32, id int32) {
	inserted_node := new_node(score, id, RED, nil, nil)
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

//---------------------------------------------------------- Delete an element
// WRITE-LOCK
func (t *Tree) DeleteNode(n *Node) {
	// handle red-black properties, and deletion work.
	if n.left != nil && n.right != nil {
		/* Copy fields from predecessor and then delete it instead */
		pred := maximum_node(n.left)
		// copy score, id
		n.score = pred.score
		n.id = pred.id
		// deal with predecessor after.
		n = pred
	}

	// fixup from maximum_node
	fixup_size(n)

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

//---------------------------------------------------------- tree print
const INDENT_STEP = 4

func Print_helper(n *Node, indent int) {
	if n == nil {
		fmt.Printf("<empty tree>")
		return
	}
	if n.right != nil {
		Print_helper(n.right, indent+INDENT_STEP)
	}
	for i := 0; i < indent; i++ {
		fmt.Printf(" ")
	}
	if n.color == BLACK {
		fmt.Printf("[score:%v size:%v id:%v]\n", n.score, n.size, n.id)
	} else {
		fmt.Printf("*[score:%v size:%v id:%v]\n", n.score, n.size, n.id)
	}

	if n.left != nil {
		Print_helper(n.left, indent+INDENT_STEP)
	}
}
