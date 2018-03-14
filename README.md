iam-get-ssh-keys
================

Retrieve SSH keys for an user in AWS IAM

Usage
-----

Get active SSH keys:

```
iam-get-ssh-keys USERNAME
```

Integrate with OpenSSH server:

```
AuthorizedKeysCommand /usr/bin/iam-get-ssh-keys %u
AuthorizedKeysCommandUser nobody
```

Don't forget to add the following policy to the role attached to your EC2 instance:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "iam:ListSSHPublicKeys",
                "iam:GetSSHPublicKey"
            ],
            "Resource": "*"
        }
    ]
}
```

Important
---------

This is a proof of concept. Not ready for production use.
