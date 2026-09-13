# Vyper grammar

Vendored from [Olyno/tree-sitter-vyper](https://github.com/Olyno/tree-sitter-vyper)
at `97d51f8b689b6a2dc62ab31b4d0d34797aec19a6` (MIT; see LICENSE).
The C parser and scanner are compiled through cgo, like Toast's other grammars.
No additional runtime or Go module dependency is needed.

Local grammar changes:

- Add current `flag` declarations alongside legacy `enum` declarations.
- Consume newlines between interface method signatures.
- Support module dependency bindings (`initializes: module[dependency := implementation]`).
- Restrict loop annotations so their type cannot consume the following `in` expression.
- Bind `extcall` and `staticcall` to the complete method call.
- Generate ABI 14, supported by Toast's pinned go-tree-sitter runtime.

To regenerate using tree-sitter-cli **0.25.9**, run from the repository root:

```sh
tree-sitter generate --abi 14 --output internal/syntax/vyper internal/syntax/vyper/grammar.js
go test ./internal/syntax/...
```

Keep `grammar.js`, generated `grammar.json`, `node-types.json`, `parser.c`,
`scanner.c`, and the `tree_sitter` headers together. Ordinary Toast builds
do not need Node.js or the Tree-sitter CLI.

The Toast-specific query is in `../queries/vyper.scm`. It uses existing theme
keys and intentionally omits a catch-all identifier capture, which would
overlap specialized captures in Toast's renderer. The grammar is for editor
highlighting; the Vyper compiler remains the authority on valid contracts.
