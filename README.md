# t

Manage todo lists in the CLI.

## Getting started

### Installation and usage

Add an item to the general todo list.

```bash
t "Do something"
```

Add an item to today's todo list.

```bash
# --tomorrow works too
t "Do something today" --today
```

Open the TUI.

```bash
t
```

### Development and testing

#### Requirements

- [Go] 1.25+

```bash
git clone git@github.com:unfunco/t.git
cd t
```

```bash
go build
```

## License

© 2025 [Daniel Morris]\
Made available under the terms of the [MIT License].

[daniel morris]: https://unfun.co
[go]: https://go.dev
[mit license]: LICENSE.md
