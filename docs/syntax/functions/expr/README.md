# expr(expression: String) -> []Node

This is one of the available generators.
It generates nodes using [CEL](https://cel.dev/overview/cel-overview).

The following variables are predefined:

- e: Environment variables, equivalent to [os.Environ](https://pkg.go.dev/os#Environ).
- n: The current node.

For example, the following expression determines if the size attribute is less than 1000 and stores the result in the small attribute:

```
expr("\"small=\" + string(n.size < 1000)")
```

If `@file` is specified as expression, the contents of the file will be used.
