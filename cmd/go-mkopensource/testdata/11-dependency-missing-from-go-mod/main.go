package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/runtime"
)

func main() {
	fmt.Println("Caller is " + runtime.GetCaller())
}
