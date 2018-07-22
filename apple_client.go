package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/pkg/errors"
	"strings"
)

type AppleClient struct {
	Username string
	Password string
	Client   http.Client
}

func NewClient(u string, p string) *AppleClient {
	c := AppleClient{
		Username: u,
		Password: p,
		Client: http.Client{
			Timeout: time.Duration(5 * time.Second),
		},
	}

	return &c
}

type AppleResponse struct {
	UserInfo struct {
		FirstName string `json:"firstName"`
	} `json:"userInfo"`
	Devices []Device `json:"content"`
}

func (r AppleResponse) Get(deviceName string) (Device, error) {
	for _, d := range r.Devices {
		if d.Name == deviceName {
			return d, nil
		}
	}

	return Device{}, fmt.Errorf("cannot find device %s", deviceName)
}

func FindDevice(c *AppleClient, deviceName string, user *string, device *Device) error {

	basicAuth := func(username, password string) string {
		auth := username + ":" + password
		return base64.StdEncoding.EncodeToString([]byte(auth))
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://fmipmobile.icloud.com/fmipservice/device/%s/initClient", c.Username), nil)
	if err != nil {
		return errors.Wrap(err, "failed creating request")
	}
	req.Header.Add("X-Apple-Realm-Support", "1.0")
	req.Header.Add("X-Apple-Find-API-Ver", "3.0")
	req.Header.Add("X-Apple-AuthScheme", "UserIDGuest")
	req.Header.Add("User-Agent", "FindMyiPhone/500 CFNetwork/758.4.3 Darwin/15.5.0")
	req.Header.Add("Authorization", "Basic "+basicAuth(c.Username, c.Password))

	resp, err := c.Client.Do(req)

	if err != nil {
		return errors.Wrap(err, "failed executing request to Apple")
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("response status code was:", resp.Status)
		return errors.Errorf("response status code was: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed reading the body into []byte")
	}

	r := AppleResponse{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return errors.Wrap(err, "failed unmarshalling apple response")
	}

	*user = r.UserInfo.FirstName
	*device, err = r.Get(deviceName)
	if err != nil {
		return errors.Wrap(err, "failed extracting device")
	}

	return nil
}

func decryptEnvCredentials(creds string) (map[string]string, error) {
	var res = make(map[string]string)

	d, err := decrypt(creds)
	if err != nil {
		return res, errors.Wrap(err,"not authorised to retrieve this information right now")
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
		CiphertextBlob: []byte(decoded),
	})
	if err != nil {
		return "", err
	}

	return string(out.Plaintext), nil
}
