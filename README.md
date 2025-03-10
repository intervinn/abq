# ABQ
ABQ is a transpiler from Go to Luau heavily inspired by roblox-ts and roblox-cs.

It consists of a Luau AST which mocks the Go's AST, so that the former would easily convert to latter. Then the transformer captures Go's expressions, statements and declaration for it to be converted into Luau analogue.

## Limitations

There are some limitations ABQ will aim to implement:

* Structs are only instantiated as pointers - in Go everything is passed by value, however in Luau the tables are purely reference based.

* Channels are to be replaced with coroutine polyfills.

* Every individual package is to be compiled and packed in a single file.

## Modding
ABQ placeholds one call expression - `transform.Mod`. It's string parameter 

## Projects
`transform` - Go AST to Luau AST transformer
`luau` - Luau AST, writer and other essentials

## TODO
* Finish up Transformer and AST
* Luau libraries polyfills
* CLI
* Packer and filesystem work