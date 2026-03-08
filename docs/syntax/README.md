# Syntax

`ndql` uses a SQL-based syntax.

## Implementation Status

- Statements: Currently, only the SELECT statement is implemented.
- Clauses: FROM and WHERE clauses are available. Other clauses (e.g., GROUP BY, ORDER BY, JOIN) are not yet supported.
- Operators, Functions: Some operators and functions are not yet implemented. Even if implemented, the behavior may differ from standard SQL specifications.

## Operators

- `AND`
- `OR`
- `XOR`
- `+` (binary)
- `-` (binary)
- `*`
- `/`
- `%`
- `<<`
- `>>`
- `<`
- `<=`
- `=`
- `<>`
- `>=`
- `>`
- `CASE`
- `IS NULL`
- `IS TRUE`
- `IS FALSE`
- `REGEXP`
- `LIKE`
- `BETWEEN`
- `-` (unary)
- `~`

# Children

- [functions](./functions/README.md)
- [generator](./generator/README.md)
