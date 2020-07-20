package generator

import (
	"bytes"
	"errors"
	"math"
	"sort"
	"strings"
)

var (
	NumberPool            = "0123456789"
	LetterInLowCasePool   = "abcdefghijklmnopqrstuvwxyz"
	LetterInUpperCasePool = strings.ToUpper(LetterInLowCasePool)
	DefaultPool           = NumberPool + LetterInLowCasePool + LetterInUpperCasePool
)

func NewPasswordGenerator(max, min int) (*PasswordGenerator, error) {
	return NewPasswordGeneratorWithPool(max, min, DefaultPool)
}

func NewNumberPasswordGenerator(max, min int) (*PasswordGenerator, error) {
	return NewPasswordGeneratorWithPool(max, min, NumberPool)
}

func NewPasswordGeneratorWithPool(max, min int, pool string) (*PasswordGenerator, error) {
	if max < 0 || min < 0 || max < min || pool == "" {
		return nil, errors.New("invalid args")
	}

	distinctPool := distinctString(pool)

	return &PasswordGenerator{
		pool:        distinctPool,
		max:         max,
		min:         min,
		current:     int64(math.Pow(float64(len(distinctPool)), float64(min-1))),
		maxPassword: int64(math.Pow(float64(len(distinctPool)), float64(max))),
	}, nil
}

type PasswordGenerator struct {
	pool        []byte
	max         int
	min         int
	current     int64
	maxPassword int64
}

func (pg *PasswordGenerator) Generate() (string, bool) {
	if pg.current < pg.maxPassword {
		pw, err := pg.decimalToAny(pg.current, int64(len(pg.pool)))
		if err != nil {
			return "", true
		}
		pg.current++
		return pw, false
	}

	return "", true
}

// 10 进制转 N 进制
func (pg *PasswordGenerator) decimalToAny(num int64, base int64) (string, error) {
	if int(base) > len(pg.pool) {
		return "", errors.New("base too big")
	}
	resByte := make([]byte, 0, 0)
	for num != 0 {
		b := pg.pool[num%base]
		resByte = append([]byte{b}, resByte...)
		num = num / base
	}
	return string(resByte), nil
}

func distinctString(s string) []byte {
	set := make(map[byte]bool, len(s))
	for i := 0; i < len(s); i++ {
		set[s[i]] = true
	}

	buf := bytes.Buffer{}
	for key := range set {
		buf.WriteByte(key)
	}

	unsorted := buf.Bytes()
	sort.Slice(unsorted, func(i, j int) bool {
		return unsorted[i] < unsorted[j]
	})
	return unsorted
}
