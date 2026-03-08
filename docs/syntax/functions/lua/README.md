# lua(script: String, entrypoint: String) -> []Node

This is one of the available generators.
It generates nodes by executing Lua scripts.

The entrypoint must specify a function predefined within the script
This function must accept exactly one argument and return a string.
The first argument is the current node, passed as a Lua table.
A global table `E` is predefined, containing environment variables equivalent to [os.Environ](https://pkg.go.dev/os#Environ).

For example, the following expression calculates the logarithm of the size attribute and stores the result in the lsize attribute:

```
lua("function f(n) return \"lsize=\" .. tostring(math.log(n.size, 10)) end", "f")
```

If `@file` is specified as script, the contents of the file will be used.
