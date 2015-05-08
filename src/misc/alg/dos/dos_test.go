package dos

import "testing"

func Benchmark(b *testing.B) {
	tree := Tree{}

	for i := 0; i < b.N; i++ {
		tree.Insert(100, int32(i))
	}

	for i := 0; i < b.N; i++ {
		tree.Locate(100, int32(i))
	}
}

func TestDos(t *testing.T) {
	tree := Tree{}

	N := 500
	t.Log("testing tree.Insert()")
	for i := 0; i < N; i++ {
		tree.Insert(int32(i), int32(N-i))
	}
	t.Log("Count:", tree.Count())

	t.Log("testing tree.Locate()")
	for i := 0; i < N; i++ {
		rank, node := tree.Locate(int32(i), int32(N-i))
		if node != nil {
			t.Log("id:", N-i, "score:", i, "rank:", rank, "ids", node.ids)
		}
	}
	t.Log("Count:", tree.Count())

	Print_helper(tree.Root(), 0)

	t.Log("testing tree.Rank()")
	for i := 1; i <= tree.Count()+1; i++ {
		id, node := tree.Rank(i)
		if node != nil {
			t.Log("rank:", i, "id", id)
		}
	}
	t.Log("Count:", tree.Count())

	t.Log("testing tree.Delete()")
	cnt := tree.Count()
	for i := 0; i < cnt; i++ {
		id, n := tree.Rank(1)
		if n != nil {
			t.Log("deleting id", id)
			tree.Delete(id, n)
		}
	}
	t.Log("Count:", tree.Count())
	Print_helper(tree.Root(), 0)

	t.Log("testing tree.Locate()")
	for i := 0; i < N; i++ {
		rank, n := tree.Locate(int32(i), int32(20-i))
		if rank != -1 {
			t.Logf("score %v, ids %v rank %v \n", n.Score(), n.Ids(), rank)
		}
	}
	t.Log("Count:", tree.Count())
}
