package generator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PasswordGenerator(t *testing.T) {
	as := assert.New(t)

	generator, err := NewPasswordGenerator(1, 4)
	as.Nil(err)

	for {
		pw, last := generator.Generate()
		if last {
			break
		}

		fmt.Printf("new password %s\n", pw)
	}
}
