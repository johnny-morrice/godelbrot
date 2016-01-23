package libgodelbrot

import (
    "math/big"
    "log"
)

// Parse a big.Float
func parseBig(number string) (*big.Float, error) {
    f, _, err := big.ParseFloat(number, DefaultBase, DefaultHighPrec, big.ToNearestEven)
    return f, err
}

func emitBig(b *big.Float) string {
    return b.Text('e', bitDecAcc(b.MinPrec()))
}

func bitDecAcc(bits uint) int {
    const maxi64 = ^uint64(0) >> 1
    const topi = (^uint(0) >> 1) - 1
    if uint64(bits) > maxi64 {
        log.Panic("int64 would overflow")
    }
    digits := big.NewInt(int64(bits))
    two := big.NewInt(2)
    ten := big.NewInt(10)
    digits.Exp(two, digits, nil)

    var count int
    for count = int(1); digits.Cmp(ten) != -1; count++ {
        digits.Div(digits, ten)
        // Check for overflow
        // TODO analytic solution before loop
        if uint(count) == topi {
            log.Panic("int would overflow")
        }
    }

    return count
}