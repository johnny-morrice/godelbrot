package godelbrot

import (
	"log"
	"math/big"
)

// Parse a big.Float
func parseBig(number string) (*big.Float, error) {
	bits := digits2bits(uint(len(number)))
	f, _, err := big.ParseFloat(number, DefaultBase, bits, big.ToNearestEven)
	return f, err
}

func emitBig(b *big.Float) string {
	digits := bits2digits(b.MinPrec())
	return b.Text('e', int(digits))
}

const maxi64 = ^uint64(0) >> 1

// Returns number of bits given number of digits
func digits2bits(digits uint) uint {
	if uint64(digits) > maxi64 {
		log.Panic("int64 would overflow")
	}

	bits := big.NewInt(int64(digits))
	two := big.NewInt(2)
	ten := big.NewInt(10)
	bits.Exp(ten, bits, nil)

	count := uint(1)
	for ; bits.Cmp(two) != -1; count++ {
		bits.Div(bits, two)
	}

	return count
}

// Returns number of digits given number of bits
func bits2digits(bits uint) uint {
	if uint64(bits) > maxi64 {
		log.Panic("int64 would overflow")
	}
	digits := big.NewInt(int64(bits))
	two := big.NewInt(2)
	ten := big.NewInt(10)
	digits.Exp(two, digits, nil)

	count := uint(1)
	for ; digits.Cmp(ten) != -1; count++ {
		digits.Div(digits, ten)
	}

	return count
}
