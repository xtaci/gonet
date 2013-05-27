package interval_tree

import (
	"fmt"
	"testing"
)

func TestIntervalTree(t *testing.T) {
	tree := Tree{}

	for i := int64(0); i < 100; i++ {
		tree.Insert(i, i+10, string([]byte{byte(i%26 + 65)}))
	}

	print_helper(tree.root, 0)

	fmt.Printf("searching for %v:%v\n", 10, 10)
	node := tree.Lookup(10, 10)
	fmt.Println(node)
}

const INDENT_STEP = 4

func print_helper(n *Node, indent int) {
	if n == nil {
		fmt.Printf("<empty tree>")
		return
	}
	if n.right != nil {
		print_helper(n.right, indent+INDENT_STEP)
	}
	for i := 0; i < indent; i++ {
		fmt.Printf(" ")
	}
	if n.color == BLACK {
		fmt.Printf("[%v %v m->%v Value: %v]\n", n.low, n.high, n.m, n.Data())
	} else {
		fmt.Printf("*[%v %v m->%v Value %v]\n", n.low, n.high, n.m, n.Data())
	}

	if n.left != nil {
		print_helper(n.left, indent+INDENT_STEP)
	}
}
