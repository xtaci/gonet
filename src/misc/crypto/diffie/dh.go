/***********************************************************
 *
 * Diffie–Hellman key exchange
 *
 */
package diffie

import "math/big"
import "math/rand"
import "time"

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var DH1BASE = big.NewInt(3)
var DH1PRIME, _ = big.NewInt(0).SetString("0xFFFFFFFB", 0)

//----------------------------------------------- Diffie-Hellman Key-exchange
func DHGenKey(G, P *big.Int) (*big.Int, *big.Int) {
	X := big.NewInt(0).Rand(rng, P)
	E := big.NewInt(0).Exp(G, X, P)
	return X, E
}
