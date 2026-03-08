# grep(pattern: String, template: String) -> []Node

This is one of the available generators.
It greps the file pointed to by the path attribute using a specified pattern, then applies the captured strings to a template.

For example, the following expression roughly extracts Go function definitions and stores the function names in the func attribute:

```
grep("func (?P<name>[^(]+)", "func=$name")
```
