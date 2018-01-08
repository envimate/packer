package envibin

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/envimate/config"
)

// Lookup lookup
func Lookup(repo, base, tag string) (string, error) {
	// get envibin config
	baseURL, username, password, err := fromConfig(repo)
	if err != nil {
		return "", fmt.Errorf("could not read envibin configuration: %v", err)
	}
	url := fmt.Sprintf("%s/%s/presigned/versions/%s", baseURL, base, tag)

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

func fromConfig(prefix string) (baseURL, username, password string, err error) {
	cfg, err := config.Get(prefix)
	if err != nil {
		return "", "", "", err
	}

	baseURL = cfg.Repo.URL
	username = cfg.Repo.Username
	password = cfg.Repo.Password
	return
}
