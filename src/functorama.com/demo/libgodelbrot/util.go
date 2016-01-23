package libgodelbrot

import (
    "math/big"
)

// Parse a big.Float
func parseBig(number string) (*big.Float, error) {
    f, _, err := big.ParseFloat(number, DefaultBase, DefaultHighPrec, big.ToNearestEven)
    return f, err
}

func emitBig(b *big.Float) string {
    return b.Text('e', -1)
}