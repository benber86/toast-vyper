# Toast with Vyper highlighting

A modified build of [Toast](https://github.com/paradise-runner/toast) with
Tree-sitter syntax highlighting for Vyper contract (`.vy`) and interface (`.vyi`)
files. Highlighting activates automatically when you open either file type.

Supports decorators, built-in types, functions, events, structs, flags, module
bindings, typed loops, `log`, `extcall`, and `staticcall`, using Toast's existing
themes. This adds syntax highlighting; it does not add a Vyper language server.

## Install

Toast embeds its grammars and highlight queries at build time. Build this
repository to use Vyper highlighting; there is no highlighter file to install
into an existing Toast binary. Upstream Homebrew packages do not include this
change.

Prerequisites: Git, Go 1.25.2 or newer, Make, and a C compiler (GCC or Clang).
On macOS, the Xcode Command Line Tools supply Make and Clang. Node.js and the
Tree-sitter CLI are not needed for an ordinary build.

```sh
gh repo clone benber86/toast-vyper
cd toast-vyper
make build
./bin/toast /path/to/contract.vy
```

The `gh` clone command uses your GitHub login, including for private repositories.
You can also clone `https://github.com/benber86/toast-vyper.git` using Git with
GitHub authentication configured.

Install the binary as `toast-vyper` to keep it alongside an existing Toast:

```sh
mkdir -p "$HOME/.local/bin"
install -m755 bin/toast "$HOME/.local/bin/toast-vyper"
export PATH="$HOME/.local/bin:$PATH"
toast-vyper /path/to/contract.vy
```

Add that PATH export to your shell configuration if needed. To use the command
`toast` instead, install to `$HOME/.local/bin/toast` and ensure that directory
comes before your existing Toast installation on PATH.

## Development

```sh
go test ./...
go vet ./...
```

The Vyper grammar, its source attribution, and regeneration instructions live
in [internal/syntax/vyper](internal/syntax/vyper/README.md). The highlighting
query is [vyper.scm](internal/syntax/queries/vyper.scm), and language registration
is in [languages.go](internal/syntax/languages.go).

Based on Toast commit `32dd375`. The implementation also evaluates Tree-sitter
query predicates so type and keyword captures match only the intended names.
Tests cover Vyper syntax, file detection, predicate filtering, and incremental
edits.

See the [upstream README](README.upstream.md) for Toast's editor features,
keybindings, and configuration. Installation links in that document refer to
the upstream build.

## License

Toast is MIT licensed; see [LICENSE](LICENSE). The vendored Vyper grammar is
also MIT licensed, with its license and provenance retained in
[internal/syntax/vyper](internal/syntax/vyper/README.md).
