package apple

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client represents an HTTP client with Apple credentials injected
type Client struct {
	Username string
	Password string
	Client   http.Client
}

// NewClient returns a pointer to a new Client with the username/password
// injected and a default timeout of 5 seconds
func NewClient(u string, p string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := Client{
		Username: u,
		Password: p,
		Client: http.Client{
			Transport: tr,
			Timeout:   5 * time.Second,
		},
	}

	return &c
}

// Response struct is used to unmarshal the raw response from iCloud
// services
type Response struct {
	UserInfo struct {
		FirstName string `json:"firstName"`
	} `json:"userInfo"`
	Devices []Device `json:"content"`
}

// Get on the Response receiver picks the device matching the requested
// deviceName, error if the device does not exist
func (r Response) Get(deviceName string) (Device, error) {
	for _, d := range r.Devices {
		if d.Name == deviceName {
			return d, nil
		}
	}

	return Device{}, fmt.Errorf("cannot find device %s", deviceName)
}

// FindDevice calls iCloud fmipmobile service to retrieve the list of associated
// devices
func FindDevice(c *Client, deviceName string, user *string, device *Device) error {

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

	r := Response{}
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
