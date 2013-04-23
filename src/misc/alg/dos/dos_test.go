package dos

import "testing"
import "fmt"

func TestDos(t *testing.T) {
	tree := Tree{}

	for i := 0; i < 100; i++ {
		tree.Insert(i, 100-i)
	}

	print_helper(tree.Root(), 0)

	if tree.Rank(1).DATA.(int) != 1 {
		t.Error("dynamic order stat failed")
	}

	// delete 50 elements
	fmt.Println("delete 50 elements")
	for i := 1; i <= 50; i++ {
		if tree.Rank(1) != nil {
			tree.DeleteNode(tree.Rank(1))
		}
	}

	print_helper(tree.Root(), 0)

	fmt.Println("add another 50 elements")
	// adding other 50 elements
	for i := 1000; i < 1051; i++ {
		tree.Insert(i, i*2)
	}

	print_helper(tree.Root(), 0)

	if tree.Rank(1).DATA.(int) != 2100 {
		t.Error("dynamic order stat failed")
	}

}

const INDENT_STEP = 4

func print_helper(n *node, indent int) {
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
		fmt.Printf("[score:%v DATA:%v]\n", n.score, n.DATA)
	} else {
		fmt.Printf("*[score:%v DATA:%v]\n", n.score, n.DATA)
	}

	if n.left != nil {
		print_helper(n.left, indent+INDENT_STEP)
	}
}
