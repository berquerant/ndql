# sh(script: String) -> []Node

This is one of the available generators.
It generates nodes by executing bash scripts.

Environment variables are available directly within the script.
To retrieve attribute values from a node, use the following functions:

- get NAME: Retrieves the value of the specified attribute. Returns an empty string if the attribute is not found.
- get_or NAME DEFAULT_VALUE: Retrieves the value of the specified attribute. Returns DEFAULT_VALUE if the attribute is not found.

For example, the following expression retrieves the first line of the file pointed to by the path attribute and stores it in the head attribute:

```
sh("echo head=$(head -n1 $(get path))")
```

If `@file` is specified as script, the contents of the file will be used.
