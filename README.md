# Drop

A peer-to-peer file transfer tool for your local network, written in Go. Think AirDrop, but from the terminal.

## How it works

1. The **receiver** starts listening and broadcasts its presence on the local network via UDP.
2. The **sender** discovers available receivers, picks one, and initiates a TCP connection.
3. The receiver sees the incoming filename and can **accept or reject** the transfer before any file data is sent.
4. On acceptance, the file is streamed directly over TCP.

## Install

**Homebrew:**

```
brew install arjunkomath/tap/drop
```

**From source:**

```
go install github.com/arjunkomath/go-drop@latest
```

Or download compiled binaries from [releases](https://github.com/arjunkomath/go-drop/releases).

## Usage

**Receive a file:**

```
drop receive
```

When a sender connects, you'll be prompted to accept or reject the incoming file.

**Send a file:**

```
drop send <file>
```

Select a discovered device from the list and press Enter to start the transfer.
