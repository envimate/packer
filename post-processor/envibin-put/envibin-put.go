package envibinPut

import (
	"fmt"

	"errors"

	envibinConfig "bitbucket.org/envimate/config"
	"bitbucket.org/envimate/envibin-cli/domain"
	"bitbucket.org/envimate/envibin-go-client"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
)

var cfg *envibinConfig.Config

type Config struct {
	Repository string   `mapstructure:"repository"`
	Artifact   string   `mapstructure:"artifact"`
	Version    string   `mapstructure:"version"`
	Tags       []string `mapstructure:"tags"`
	Labels     []string `mapstructure:"labels"`
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, nil, raws...)
	if err != nil {
		return err
	}

	if p.config.Artifact == "" {
		return errors.New("envibin-put: artifact cannot be empty")
	}
	if p.config.Version == "" {
		return errors.New("envibin-put: version cannot be empty")
	}

	err = resolveEnvibinConfig(p.config.Repository)
	return err
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	var artifactKey = domain.NewArtifactKey(fmt.Sprintf("%s:%s", p.config.Artifact, p.config.Version))
	if artifactKey.IsInvalid() {
		return nil, false, fmt.Errorf("invalid artifact-key %s", artifactKey.String())
	}

	c, err := client.New(cfg.Repo.URL)
	if err != nil {
		return nil, false, err
	}

	compressedFile, err := compressed(artifact.Files())
	if err != nil {
		return nil, false, err
	}

	_ = c.Remove(artifactKey)
	err = c.Push(artifactKey, compressedFile, p.config.Tags, p.config.Labels)
	if err != nil {
		return nil, false, err
	}

	ui.Message("put artifact with key " + fmt.Sprintf("%s:%s", p.config.Artifact, p.config.Version))

	return nil, false, nil
}

func resolveEnvibinConfig(prefix string) error {
	if prefix == "" {
		err := envibinConfig.Init()
		if err != nil {
			return fmt.Errorf("Could not read envibin default configuration: %s", err)
		}
		cfg = envibinConfig.Default
	} else {
		c, err := envibinConfig.Get(prefix)
		if err != nil {
			return fmt.Errorf("Could not read envibin configuration for %s: %s", prefix, err)
		}
		cfg = c
	}

	return nil
}

func compressed(files []string) (string, error) {
	if len(files) == 0 {
		return "", errors.New("no files to upload")
	}
	if len(files) > 0 {
		return files[0], nil
	} else {
		return "", errors.New("not yet supported - artifact contains more then one file")
	}

	return "", nil
}
