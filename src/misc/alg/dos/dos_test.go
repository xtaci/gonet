package dos

import "testing"
import "fmt"

func Benchmark(b *testing.B) {
	tree := Tree{}

	for i := 0; i < b.N; i++ {
		tree.Insert(int32(i), int32(i))
	}

	fmt.Println(tree.Count())
	for i := 0; i < b.N; i++ {
		n := tree.Rank(101)
		if n != nil {
			tree.DeleteNode(n)
		}
	}
	fmt.Println(tree.Count())
}

func TestDos(t *testing.T) {
	tree := Tree{}

	for i := 0; i < 100; i++ {
		tree.Insert(int32(i), int32(100-i))
	}

	for i := 0; i < 100; i++ {
		n, rank := tree.Locate(int32(i), int32(100-i))
		fmt.Println(n, rank)
	}

	Print_helper(tree.Root(), 0)

	if tree.Rank(1).id != 1 {
		t.Error("dynamic order stat failed")
	}

	for i := 0; i < 100; i++ {
		rank, _ := tree.Locate(int32(i), int32(100-i))
		fmt.Printf("score %v, rank %v \n", i, rank)
	}

	for i := 1; i < 50; i++ {
		n := tree.Rank(1)
		if n != nil {
			tree.DeleteNode(n)
		}
	}
	Print_helper(tree.Root(), 0)

	for i := 0; i < 100; i++ {
		rank, n := tree.Locate(int32(i), int32(100-i))
		if rank != -1 {
			fmt.Printf("score %v, id %v rank %v \n", n.Score(), n.Id(), rank)
		}
	}
	fmt.Println(tree.Count())
}
