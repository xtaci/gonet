package dos

import "testing"
import "fmt"

func Benchmark(b *testing.B) {
	tree := Tree{}
	for i := 0; i < b.N; i++ {
		tree.Insert(i, i)
	}

	for i := 0; i < b.N; i++ {
		n := tree.Rank(i)
		if n != nil {
			tree.DeleteNode(n)
		}
	}
}

func TestDos(t *testing.T) {
	tree := Tree{}

	for i := 0; i < 100; i++ {
		tree.Insert(i, 100-i)
	}

	print_helper(tree.Root(), 0)

	if tree.Rank(1).Data().(int) != 1 {
		t.Error("dynamic order stat failed")
	}

	for i := 0; i < 100; i++ {
		node, rank := tree.ByScore(100 - i)
		fmt.Printf("score %v, rank %v, node %v\n", 100-i, rank, node)
	}

	_, rank := tree.ByScore(97)
	if rank != 3 {
		t.Error("get by score failed %v", rank)
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

	if tree.Rank(1).Data().(int) != 2100 {
		t.Error("dynamic order stat failed")
	}

	tree.Clear()
	print_helper(tree.Root(), 0)
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
		fmt.Printf("[score:%v size:%v]\n", n.score, n.size)
	} else {
		fmt.Printf("*[score:%v size:%v]\n", n.score, n.size)
	}

	if n.left != nil {
		print_helper(n.left, indent+INDENT_STEP)
	}
}
