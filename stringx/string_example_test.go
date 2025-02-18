package stringx

import "fmt"

func ExampleUnsafeToBytes() {
	str := "hello"
	val := UnsafeToBytes(str)
	fmt.Println(len(val))
	// Output:
	// 5
}

func ExampleUnsafeToString() {
	val := UnsafeToString([]byte("hello"))
	fmt.Println(val)
	// Output:
	// hello
}
