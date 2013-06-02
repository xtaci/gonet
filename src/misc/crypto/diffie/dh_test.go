package diffie

import (
	"fmt"
	"math/big"
	"testing"
)

func TestDH(t *testing.T) {
	X1, E1 := DHGenKey(DH1BASE, DH1PRIME)
	X2, E2 := DHGenKey(DH1BASE, DH1PRIME)

	fmt.Println("Secret 1:", X1, E1)
	fmt.Println("Secret 2:", X2, E2)

	KEY1 := big.NewInt(0).Exp(E2, X1, DH1PRIME)
	KEY2 := big.NewInt(0).Exp(E1, X2, DH1PRIME)

	fmt.Println("KEY1:", KEY1)
	fmt.Println("KEY2:", KEY2)

	if KEY1.Cmp(KEY2) != 0 {
		t.Error("Diffie-Hellman failed")
	}
}

func BenchmarkDH(b *testing.B) {
	for i := 0; i < b.N; i++ {
		X1, E1 := DHGenKey(DH1BASE, DH1PRIME)
		X2, E2 := DHGenKey(DH1BASE, DH1PRIME)

		big.NewInt(0).Exp(E2, X1, DH1PRIME)
		big.NewInt(0).Exp(E1, X2, DH1PRIME)
	}
}
