// keygen.go: generate an Ed25519 keypair for MCP Registry DNS verification.
//
// Usage:
//
//	go run keygen.go
//
// Outputs:
//
//	PRIVATE_KEY_HEX: 64-character hex string; passed to mcp-publisher.
//	PUBLIC_KEY_HEX:  64-character hex string; embedded in the DNS TXT record.
//	DNS_TXT_VALUE:   the literal value for a TXT record at the apex of your domain.
//
// Store the private key hex securely. The public key is published in DNS
// and is safe to expose.
//
// This file is intentionally committed to the docs repo. The keys it generates
// must NEVER be committed; capture the stdout, paste into DNS + the publisher,
// and save the private hex in a password manager.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Fprintln(os.Stderr, "keygen failed:", err)
		os.Exit(1)
	}

	privHex := hex.EncodeToString(priv.Seed())
	pubHex := hex.EncodeToString(pub)

	fmt.Println("PRIVATE_KEY_HEX:", privHex)
	fmt.Println("PUBLIC_KEY_HEX: ", pubHex)
	fmt.Println()
	fmt.Println("DNS_TXT_VALUE (paste as the value of a TXT record at the APEX of ttrpg.bot):")
	fmt.Printf("v=MCPv1; k=ed25519; p=%s\n", pubHex)
}
