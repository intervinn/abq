# abq
abq is my random thought of transpiling go code into luau (a lua-based language from roblox with types) just like roblox-ts
i dont know if i can carry this idea so it can be actually used to develop games on roblox, i do not reside a large experience of programming in golang, and especially its internals 

abq yet doesnt provide a roblox api, but if it will i think of having empty body functions that transpiler will ignore

# current abq compared to roblox-ts
* go is a compiled language so it will compile your stuff blazingly fast
* problem of compiled lang is that its a complete program, unlike scripting langs like typescript or luau, meaning your experience with abq will either be more overcomplicated or more forcefully organised 
* you write lego games in C of 21th century instead of a lang that is made to fix issues of a language that runs in browser

# tests
```go
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
```

```lua
function test2()
print("hi 2")
end 
function test3()
print("hi test3")
end 
function test(num : number, s : string) 
print("hello world", num, s)
test2()
end 
test(10, "a string")
```