package main

import (
	"fmt"
	"github.com/josharian/intern"

	_ "k8s.io/apimachinery/pkg/util/clock"
)

func main() {
	fmt.Println(intern.String("Hello, world!"))
}
