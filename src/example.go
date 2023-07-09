package src

func test2() {
	print("hi 2")
}

func test3() {
	print("hi test3")
}

func test(num int, s string) {
	print("hello world", num, s)
	test2()
	test3()
}

func main() {
	test(10, "a string")
}
