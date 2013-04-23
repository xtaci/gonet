package dos

import "testing"
import "fmt"

func TestDos(t *testing.T) {
	tree := Tree{}

	data := make([]int, 100)
	for i:=0;i<100;i++ {
		data[i] = 100-i
		tree.Insert(i, data[i])
	}

	for i:=1;i<=100;i++ {
		if tree.Rank(i) != nil {
			fmt.Println(tree.Rank(i).DATA)
		}
	}
}
