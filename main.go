package main

import (
	"flag"
	"fmt"
	"log/syslog"
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
	logger, err := syslog.New(syslog.LOG_NOTICE, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Can't connect to syslog: %s\n", err)
		os.Exit(1)
		return
	}
	defer logger.Close()

	flag.Parse()
	userName := flag.Arg(0)

	logger.Notice(fmt.Sprintf("Getting SSH keys for user %s in AWS IAM", userName))

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := iam.New(sess)

	keys, err := readSSHKeys(client, &userName)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to get SSH keys for user %s in AWS IAM: %s", userName, err))
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
		return
	}

	for _, key := range keys {
		fmt.Println(*key)
	}
}
