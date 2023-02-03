package main

import (
	"fmt"
	"github.com/josharian/intern"

	_ "k8s.io/client-go/rest"
)

func main() {
	fmt.Println(intern.String("Hello, world!"))
}
