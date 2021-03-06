// +build ignore

package main

import (
  ".."
  "crypto/rand"
  "crypto/rsa"
  "crypto/x509"
  "encoding/base64"
  "flag"
  "log"
)

var (
  address = flag.String("address", ":8192", "what host and port to listen to")
)

func main() {
  flag.Parse()

  key, err := rsa.GenerateKey(rand.Reader, 2048)
  if err != nil {
    log.Fatalf("generating private key: %v", err)
  }

  pubKey, err := x509.MarshalPKIXPublicKey(key.Public())
  if err != nil {
    log.Fatalf("serializing private key: %v", err)
  }

  encodedPubKey := base64.StdEncoding.EncodeToString(pubKey)
  tokenPrime, err := rand.Prime(rand.Reader, 130)
  if err != nil {
    log.Fatalf("creating token: %v", err)
  }
  token := tokenPrime.Text(36)

  log.Println("Listening on " + *address)
  log.Println("Here's your public key: " + encodedPubKey)
  log.Println("Your v2 token: " + token)

  server := votifier.NewServer(key, func(vote votifier.Vote, version int) {
    log.Println("Got vote: ", vote, ", version: " , version)
  }, votifier.StaticServiceTokenIdentifier(token))
  server.ListenAndServe(*address)
}
