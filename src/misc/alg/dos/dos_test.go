package dos

import "testing"
import "fmt"

func TestDos(t *testing.T) {
	tree := Tree{}

	for i:=0;i<100;i++ {
		tree.Insert(i, 100-i)
	}

	for i:=1;i<=100;i++ {
		if tree.Rank(i) != nil {
			fmt.Println(tree.Rank(i).DATA)
		}
	}

	if tree.Rank(1).DATA.(int) != 1 {
		t.Error("dynamic order stat failed")
	}

	// delete 50 elements
	fmt.Println("delete 50 elements")
	for i:=1;i<=50;i++ {
		if tree.Rank(1) != nil {
			tree.DeleteNode(tree.Rank(1))
		}
	}

	for i:=1;i<=50;i++ {
		if tree.Rank(i) != nil {
			fmt.Println(tree.Rank(i).DATA)
		}
	}

	// adding other 50 elements
	for i:=1000;i<1051;i++ {
		tree.Insert(i, i*2)
	}

	if tree.Rank(1).DATA.(int) != 2100 {
		t.Error("dynamic order stat failed")
	}
}
