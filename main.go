package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func readSSHKeys(client *iam.IAM, userName *string) ([]*string, error) {
	publicKeys, err := client.ListSSHPublicKeys(&iam.ListSSHPublicKeysInput{UserName: userName})
	if err != nil {
		return nil, fmt.Errorf("Failed to list SSH keys for user %s: %s", *userName, err)
	}

	results := make([]*string, 0, len(publicKeys.SSHPublicKeys))

	for _, publicKey := range publicKeys.SSHPublicKeys {
		encoding := "SSH"
		publicKeyDetails, err := client.GetSSHPublicKey(&iam.GetSSHPublicKeyInput{
			UserName:       userName,
			SSHPublicKeyId: publicKey.SSHPublicKeyId,
			Encoding:       &encoding,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to get public key with ID %s: %s", *publicKey.SSHPublicKeyId, err)
		}

		// Ignore inactive keys
		if *publicKeyDetails.SSHPublicKey.Status != "Active" {
			break
		}

		formattedKey := *publicKeyDetails.SSHPublicKey.SSHPublicKeyBody + " " + *publicKeyDetails.SSHPublicKey.SSHPublicKeyId
		results = append(results, &formattedKey)
	}

	return results, err
}

func main() {
	flag.Parse()
	userName := flag.Arg(0)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := iam.New(sess)

	keys, err := readSSHKeys(client, &userName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
		return
	}

	for _, key := range keys {
		fmt.Println(*key)
	}
}
