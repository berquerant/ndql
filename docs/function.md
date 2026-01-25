# Functions

## grep(pattern: String, template: String) -> []Node

This is one of the available generators.
It greps the file pointed to by the path attribute using a specified pattern, then applies the captured strings to a template.

For example, the following expression roughly extracts Go function definitions and stores the function names in the func attribute:

``` sql
grep("func (?P<name>[^(]+)", "func=$name")
```

## tmpl(template: String) -> []Node

This is one of the available generators.
It generates nodes using [text/template](https://pkg.go.dev/text/template).
The current node is passed as the data for the template.

Additionally, the following functions are predefined:

- env: Wrapper for [os.Getenv](https://pkg.go.dev/os#Getenv).
- envor: Similar to [os.Getenv](https://pkg.go.dev/os#Getenv), but allows a default value as the second argument. It returns the default value if os.Getenv returns an empty string.

For example, the following expression sets the type attribute to "dir" if the is_dir attribute is true, and "file" otherwise:

``` sql
tmpl("type={{if .is_dir}}dir{{else}}file{{end}}")'
```

If `@file` is specified as template, the contents of the file will be used.

## sh(script: String) -> []Node

This is one of the available generators.
It generates nodes by executing bash scripts.

Environment variables are available directly within the script.
To retrieve attribute values from a node, use the following functions:

- get NAME: Retrieves the value of the specified attribute. Returns an empty string if the attribute is not found.
- get_or NAME DEFAULT_VALUE: Retrieves the value of the specified attribute. Returns DEFAULT_VALUE if the attribute is not found.

For example, the following expression retrieves the first line of the file pointed to by the path attribute and stores it in the head attribute:

``` sql
sh("echo head=$(head -n1 $(get path))")
```

If `@file` is specified as script, the contents of the file will be used.

## lua(script: String, entrypoint: String) -> []Node

This is one of the available generators.
It generates nodes by executing Lua scripts.

The entrypoint must specify a function predefined within the script
This function must accept exactly one argument and return a string.
The first argument is the current node, passed as a Lua table.
A global table `E` is predefined, containing environment variables equivalent to [os.Environ](https://pkg.go.dev/os#Environ).

For example, the following expression calculates the logarithm of the size attribute and stores the result in the lsize attribute:

``` sql
lua("function f(n) return \"lsize=\" .. tostring(math.log(n.size, 10)) end", "f")
```

If `@file` is specified as script, the contents of the file will be used.

## expr(expression: String) -> []Node

This is one of the available generators.
It generates nodes using [CEL](https://cel.dev/overview/cel-overview).

The following variables are predefined:

- e: Environment variables, equivalent to [os.Environ](https://pkg.go.dev/os#Environ).
- n: The current node.

For example, the following expression determines if the size attribute is less than 1000 and stores the result in the small attribute:

``` sql
expr("\"small=\" + string(n.size < 1000)")
```

If `@file` is specified as expression, the contents of the file will be used.

## to_int(value) -> Int

See [data.md](./data.md).

## to_float(value) -> Float

See [data.md](./data.md).

## to_bool(value) -> Bool

See [data.md](./data.md).

## to_string(value) -> String

See [data.md](./data.md).

## to_time(value) -> Time

See [data.md](./data.md).

## to_duration(value) -> Duration

See [data.md](./data.md).

## len(value: String) -> Int

The number of characters in a String.

## size(value: String) -> Int

The number of bytes in a String.

## format(format: String, args...) -> String

[fmt.Sprintf](https://pkg.go.dev/fmt#Sprintf).

## strtotime(string: String, format: String) -> Time

[time.Parse](https://pkg.go.dev/time#Parse).

## timeformat(t: Time, format: String) -> String

[time.Fomat](https://pkg.go.dev/time#Time.Format).

## dir(path: String) -> String

[filepath.Dir](https://pkg.go.dev/path/filepath#Dir).

## basename(path: String) -> String

[filepath.Base](https://pkg.go.dev/path/filepath#Base).

## extension(path: String) -> String

[filepath.Ext](https://pkg.go.dev/path/filepath#Ext).

## abspath(path: String) -> String

[filepath.Abs](https://pkg.go.dev/path/filepath#Abs).

## relpath(path: String, base: String) -> String

[filepath.Rel](https://pkg.go.dev/path/filepath#Rel)

## inverse(value: Float | Int) -> Float

Calculate inverse of the value.

## inverse(value: String) -> String

Reverse the String.

## env(name: String) -> String

[os.Getenv](https://pkg.go.dev/os#Getenv).

## envor(name: String, default: String) -> String

[os.Getenv](https://pkg.go.dev/os#Getenv), returns default if empty.
