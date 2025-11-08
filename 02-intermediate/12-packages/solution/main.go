package main

import (
	"fmt"
	"example.com/packages/calculator"
	"example.com/packages/models"
	"example.com/packages/utils"
)

func main() {
	sum := calculator.Add(5, 3)
	fmt.Printf("5 + 3 = %d\n", sum)
	
	diff := calculator.Subtract(10, 4)
	fmt.Printf("10 - 4 = %d\n", diff)
	
	user := models.NewUser("Alice", "alice@example.com")
	fmt.Printf("User: %s\n", user.Name)
	
	reversed := utils.Reverse("hello")
	fmt.Printf("Reversed: %s\n", reversed)
}
