# Type Conversion

String-based resolvers (header, query, path, cookie) automatically convert raw string values to the declared Go type of the struct field.

## Supported Types

| Go Type | Example Input | Result |
|---------|---------------|--------|
| `string` | `"hello"` | `"hello"` |
| `bool` | `"true"`, `"1"`, `"false"`, `"0"` | `true`, `true`, `false`, `false` |
| `int` | `"42"` | `42` |
| `int8` | `"127"` | `127` |
| `int16` | `"32000"` | `32000` |
| `int32` | `"100000"` | `100000` |
| `int64` | `"9999999999"` | `9999999999` |
| `uint` | `"42"` | `42` |
| `uint8` | `"255"` | `255` |
| `uint16` | `"65535"` | `65535` |
| `uint32` | `"100000"` | `100000` |
| `uint64` | `"9999999999"` | `9999999999` |
| `float32` | `"3.14"` | `3.14` |
| `float64` | `"3.14159265"` | `3.14159265` |
| `*int` (pointer) | `"42"` | pointer to `42` |
| `*string` (pointer) | `"hello"` | pointer to `"hello"` |

## Empty Values

When the raw string is empty (missing header, missing query param), the field receives the **zero value** of its type:

| Type | Zero Value |
|------|------------|
| `string` | `""` |
| `int` | `0` |
| `bool` | `false` |
| `float64` | `0.0` |
| `*int` | `nil` |

This means missing optional params don't cause errors — they just get their zero value. Use pointer types if you need to distinguish "missing" from "zero".

## Error Behavior

If conversion fails (e.g., `"abc"` for an `int` field), the resolver returns an error and the adapter responds with 400:

```json
{"error": "resolve query \"page\": strconv.ParseInt: parsing \"abc\": invalid syntax"}
```

## Not Yet Supported

- Slices (e.g., `?tag=a&tag=b` → `[]string{"a", "b"}`)
- Maps
- Custom types implementing `encoding.TextUnmarshaler`
- Time parsing (`time.Time`)

These are on the [roadmap](./roadmap.md).
