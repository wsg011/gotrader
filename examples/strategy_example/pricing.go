package main

import "fmt"

func (ms MockStrategy) Pricing() {
	fmt.Println("Pricing")
}
