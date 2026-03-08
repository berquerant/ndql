# tmpl(template: String) -> []Node

This is one of the available generators.
It generates nodes using [text/template](https://pkg.go.dev/text/template).
The current node is passed as the data for the template.

Additionally, the following functions are predefined:

- env: Wrapper for [os.Getenv](https://pkg.go.dev/os#Getenv).
- envor: Similar to [os.Getenv](https://pkg.go.dev/os#Getenv), but allows a default value as the second argument. It returns the default value if os.Getenv returns an empty string.

For example, the following expression sets the type attribute to "dir" if the is_dir attribute is true, and "file" otherwise:

```
tmpl("type={{if .is_dir}}dir{{else}}file{{end}}")'
```

If `@file` is specified as template, the contents of the file will be used.
