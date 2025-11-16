# t

Manage your todo lists in the CLI.

## Getting started

### Installation and usage

#### macOS

```bash
brew tap unfunco/tap
brew install t
```

Add an item to the general todo list:

```bash
t "Do something"
```

Add an item to today's todo list:

```bash
t "Do something today" --today
```

Add an item to tomorrow's todo list:

```bash
t "Do something tomorrow" --tomorrow
```

Open the TUI:

```bash
t
```

### Development and testing

#### Requirements

- [Go] 1.25+

Clone the repository and navigate to the `t` directory.

```bash
git clone git@github.com:unfunco/t.git
cd t
```

```bash
go build
```

Play with some test data:

```bash
mkdir testdata/t
cp testdata/*.json testdata/t/
XDG_DATA_HOME="$PWD/testdata" XDG_CONFIG_HOME="$PWD/testdata" ./t
```

## License

© 2025 [Daniel Morris]\
Made available under the terms of the [MIT License].

[daniel morris]: https://unfun.co
[go]: https://go.dev
[mit license]: LICENSE.md
