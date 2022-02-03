package main

import (
	"fmt"

	"github.com/josharian/intern"
	_ "github.com/stretchr/testify/assert"
)

func main() {
	fmt.Println(intern.String("Hello, world!"))
}
