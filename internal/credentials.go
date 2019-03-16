package internal

import (
	"encoding/base64"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/pkg/errors"
)

// DecryptEnvCredentials fetches the credentials from the ENV vars and decrypts // them
func DecryptEnvCredentials(creds string) (map[string]string, error) {
	var res = make(map[string]string)

	d, err := decrypt(creds)
	if err != nil {
		return res, errors.Wrap(err, "not authorised to retrieve this information right now")
	}

	c := strings.Split(d, "::")
	if len(c) != 2 {
		return res, errors.New("credentials split invalid")
	}

	res["username"] = c[0]
	res["password"] = c[1]

	return res, nil
}

func decrypt(s string) (string, error) {
	// Decode string
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	// Decrypt string
	sess := session.Must(session.NewSession())
	svc := kms.New(sess, aws.NewConfig().WithRegion("eu-west-2"))
	out, err := svc.Decrypt(&kms.DecryptInput{
		CiphertextBlob: decoded,
	})
	if err != nil {
		return "", err
	}

	return string(out.Plaintext), nil
}
