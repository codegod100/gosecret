# GoSecret

A simple, easy-to-use golang implementation for managing secrets in GNOME Keyring and other freedesktop.org Secret Service compatible backends.

## Overview

GoSecret provides a simplified key/value interface for storing and retrieving secrets, making it much easier to use than traditional secret-tool while still being fully compatible with GNOME Keyring.

## Features

- **Simple key/value interface**: No complex attributes or labels needed
- **Set secrets**: Store passwords with just a key
- **Get secrets**: Retrieve passwords by key
- **Delete secrets**: Remove passwords by key  
- **List secrets**: Show all stored secrets with optional filtering

## Installation

```bash
make build
```

Or manually:
```bash
go build -o gosecret
```

## Usage

### Store a secret
```bash
# Prompt for password securely
./gosecret set mykey

# Or provide password as argument (less secure)
./gosecret set mykey mypassword

# Or pipe password from stdin
echo "mypassword" | ./gosecret set mykey
```

### Retrieve a secret
```bash
./gosecret get mykey
```

### Delete a secret
```bash
./gosecret delete mykey
```

### List all secrets
```bash
# List all secrets
./gosecret list

# Filter by pattern
./gosecret list gmail
```

## Command Line Interface

```
gosecret set <key> [value]        # store a secret (prompts for value if not provided)
gosecret get <key>                # retrieve a secret
gosecret delete <key>             # remove a secret  
gosecret list [pattern]           # list stored secrets
```

### Compatibility Aliases

For backwards compatibility, these aliases are also supported:
- `store` (same as `set`)
- `lookup` (same as `get`) 
- `clear` (same as `delete`)

## Examples

### Basic usage
```bash
# Store a password for Gmail
./gosecret set gmail.password
Password: [enter password securely]

# Retrieve it later
./gosecret get gmail.password

# List all secrets
./gosecret list

# Filter secrets containing "gmail"
./gosecret list gmail

# Delete the secret
./gosecret delete gmail.password
```

### Integration with scripts
```bash
# Store API keys
./gosecret set api.github "ghp_xxxxxxxxxxxx"
./gosecret set api.openai "sk-xxxxxxxxxxxx"

# Use in scripts
API_KEY=$(./gosecret get api.github)
curl -H "Authorization: token $API_KEY" https://api.github.com/user
```

## Dependencies

- `github.com/godbus/dbus/v5`: D-Bus communication with the secret service
- `golang.org/x/term`: Terminal password input handling

## How it Works

GoSecret communicates with the freedesktop.org Secret Service API via D-Bus, the same API used by:

- GNOME Keyring
- KDE KWallet  
- pass (with pass-secret-service)
- Other compatible secret storage backends

All secrets are stored with the `application=gosecret` attribute to keep them organized and separate from other applications.

## Advantages over secret-tool

- **Much simpler**: Just key/value pairs instead of complex attribute matching
- **No labels required**: The key serves as both identifier and label
- **Easier scripting**: Simple get/set interface perfect for automation
- **Better listing**: See all your secrets at a glance
- **Modern Go implementation**: Single binary, easy to install

## Building

```bash
make build
```

Or manually:
```bash
go mod tidy
go build -o gosecret
```

## Testing

```bash
make test
```

Or test manually:
```bash
# Store a test secret
./gosecret set test.password mysecretvalue

# Retrieve it
./gosecret get test.password  

# List it
./gosecret list test

# Clean up
./gosecret delete test.password
```