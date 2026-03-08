# Generator

A function that generates a new node from a node is called a generator.
It must return a string in one of the following formats:

- An array of JSON objects
- A single JSON object
- An "equal pair" list

The "equal pair" format is as follows:

```
key1=value11,key2=value12,...
key1=value21,key2=value22,...
...
```

This is equivalent to the following JSON structure:

```
[
  {"key1":"value11","key2":"value12",...},
  {"key1":"value21","key2":"value22",...},
  ...
]
```

Each JSON object corresponds to a single node.
Note that nodes are not required to have the same set of keys.
