# go-fast Documentation

## Table of Contents

- [Overview](./overview.md) — What go-fast is and why it exists
- [Getting Started](./getting-started.md) — Install, write your first handler, run it
- [Resolvers](./resolvers/README.md) — The tag-based parameter resolution system
  - [Body](./resolvers/body.md)
  - [Header](./resolvers/header.md)
  - [Query](./resolvers/query.md)
  - [Path](./resolvers/path.md)
  - [Cookie](./resolvers/cookie.md)
- [Adapter](./adapter.md) — How `Adapt()` wires everything together
- [Type Conversion](./type-conversion.md) — Automatic string-to-type conversion
- [DX Comparison](./dx-comparison.md) — go-fast vs Gin vs Fiber side-by-side
- [Architecture](./architecture.md) — Internal design: analyzer, metadata, resolvers, adapter
- [Roadmap](./roadmap.md) — What's coming next
