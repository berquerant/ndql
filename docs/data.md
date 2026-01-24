# Data

## Type

`ndql` supports the following data types (corresponding to Go types):

- Null (nil)
- Float (float64)
- Int (int64)
- Bool (bool)
- String (string)
- Time (time.Time)
- Duration (time.Duration)

## Cast

- ✅: Fully supported
- ⚠️: Supported with potential precision loss or specific format requirements
- ❌: Not supported

| From \ To | Null | Float | Int | Bool | String | Time | Duration |
|-----------|------|-------|-----|------|--------|------|----------|
| Null      | -    | ❌    | ❌  | ❌   | ❌     | ❌   | ❌       |
| Float     | ❌   | -     | ⚠️   | ✅   | ✅     | ⚠️    | ⚠️        |
| Int       | ❌   | ✅    | -   | ✅   | ✅     | ✅   | ✅       |
| Bool      | ❌   | ✅    | ✅  | -    | ✅     | ❌   | ❌       |
| String    | ❌   | ⚠️     | ⚠️   | ✅   | -      | ⚠️    | ⚠️        |
| Time      | ❌   | ⚠️     | ✅  | ❌   | ✅     | -    | ❌       |
| Duration  | ❌   | ⚠️     | ✅  | ❌   | ✅     | ❌   | -        |

Please note that the standard `CAST` is not yet implemented.
To perform type casting, use the following conversion functions instead:

- to_float(value): Converts value to Float.
- to_int(value): Converts value to Int.
- to_bool(value): Converts value to Bool.
- to_string(value): Converts value to String.
- to_time(value): Converts value to Time.
- to_duration(value): Converts value to Duration.
