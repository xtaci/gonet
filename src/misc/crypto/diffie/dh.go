/*
Diffie–Hellman key exchange
*/
package diffie

import "math/big"
import "math/rand"
import "time"

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var dh1_gen = big.NewInt(2)
var dh1_prime, _ = big.NewInt(0).SetString("0xFFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381FFFFFFFFFFFFFFFF", 0)

//----------------------------------------------- Diffie-Hellman Key-exchange
func DHGenKey(G, P *big.Int) (*big.Int, *big.Int) {
	X := big.NewInt(0).Rand(rng, P)
	E := big.NewInt(0).Exp(G, X, P)
	return X, E
}
