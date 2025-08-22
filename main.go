package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

const (
	DBUS_SERVICE_NAME      = "org.freedesktop.secrets"
	DBUS_INTERFACE_SERVICE = "org.freedesktop.Secret.Service"
	DBUS_PATH_SECRETS      = "/org/freedesktop/secrets"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "set", "store":
		err := actionSet(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(1)
		}
	case "get", "lookup":
		err := actionGet(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(1)
		}
	case "delete", "clear":
		err := actionDelete(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(1)
		}
	case "list":
		err := actionList(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gosecret set <key> [value]        # store a secret (prompts for value if not provided)\n")
	fmt.Fprintf(os.Stderr, "       gosecret get <key>                # retrieve a secret\n")
	fmt.Fprintf(os.Stderr, "       gosecret delete <key>             # remove a secret\n")
	fmt.Fprintf(os.Stderr, "       gosecret list [pattern]           # list stored secrets\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Aliases:\n")
	fmt.Fprintf(os.Stderr, "       store, lookup, clear are also supported for compatibility\n")
}

func readPassword() (string, error) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print("Password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return "", err
		}
		return string(password), nil
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return strings.Join(lines, "\n"), nil
	}
}

func actionSet(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("must specify a key")
	}

	key := args[0]
	var value string
	var err error

	if len(args) >= 2 {
		// Value provided as argument
		value = strings.Join(args[1:], " ")
	} else {
		// Read value from input
		value, err = readPassword()
		if err != nil {
			return fmt.Errorf("couldn't read value: %v", err)
		}
	}

	secretService, err := NewSecretService()
	if err != nil {
		return err
	}
	defer secretService.Close()

	return secretService.SetSecret(key, value)
}

func actionGet(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("must specify exactly one key")
	}

	key := args[0]

	secretService, err := NewSecretService()
	if err != nil {
		return err
	}
	defer secretService.Close()

	secret, err := secretService.GetSecret(key)
	if err != nil {
		return err
	}

	if secret == "" {
		os.Exit(1)
	}

	fmt.Print(secret)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println()
	}

	return nil
}

func actionDelete(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("must specify exactly one key")
	}

	key := args[0]

	secretService, err := NewSecretService()
	if err != nil {
		return err
	}
	defer secretService.Close()

	return secretService.DeleteSecret(key)
}

func actionList(args []string) error {
	var pattern string
	if len(args) > 0 {
		pattern = args[0]
	}

	secretService, err := NewSecretService()
	if err != nil {
		return err
	}
	defer secretService.Close()

	return secretService.ListSecrets(pattern)
}

