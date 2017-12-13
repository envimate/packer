package envibin

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

// Lookup lookup
func Lookup(repo, base, tag string) (string, error) {
	// get envibin config
	baseURL, username, password, err := fromConfig("envibin-" + repo)
	if err != nil {
		return "", fmt.Errorf("could not read envibin configuration: %v", err)
	}
	url := fmt.Sprintf("%s/%s/presigned/%s", baseURL, base, tag)

	// request to envibin
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("unable to create envibin request: %v", err)
	}
	req.SetBasicAuth(username, password)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", fmt.Errorf("error when sending request to envibin: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("envibin responded with %v, aborting", resp.StatusCode)
	}

	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unexpected response from envibin: %v", err)
	}

	return string(body), nil
}

func fromConfig(filename string) (baseURL, username, password string, err error) {
	v := viper.New()
	v.SetEnvPrefix("envi")
	v.AutomaticEnv()
	v.SetConfigName(filename) // name of config file (without extension)
	v.AddConfigPath("/etc/envibin/")
	v.AddConfigPath("$HOME/.envibin")
	v.AddConfigPath(".")
	err = v.ReadInConfig()
	if err != nil {
		return
	}

	baseURL = v.GetString("url")
	username = v.GetString("username")
	password = v.GetString("password")
	return
}
