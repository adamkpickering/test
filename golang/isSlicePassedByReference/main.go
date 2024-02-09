package main

import "fmt"

func main() {
	var mainSlice [6]string
	mySlice := mainSlice[0:3]
	mySlice[0] = "one"
	mySlice[1] = "two"
	mySlice[2] = "three"

	fmt.Println("mainSlice before:")
	fmt.Println(mainSlice)
	fmt.Println("mySlice before:")
	fmt.Println(mySlice)

	addElements(mySlice)

	fmt.Println("mainSlice after addElements:")
	fmt.Println(mainSlice)
	fmt.Println("mySlice after addElements:")
	fmt.Println(mySlice)

	modifyFirstElement(mySlice)

	fmt.Println("mainSlice after modifyFirstElement:")
	fmt.Println(mainSlice)
	fmt.Println("mySlice after modifyFirstElement:")
	fmt.Println(mySlice)
}

// this does not change the slice outside of the function
func addElements(s []string) {
	s = append(s, "four", "five", "six")
}

func modifyFirstElement(s []string) {
	s[0] = "four"
}
