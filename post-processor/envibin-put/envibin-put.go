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
	ConfigPrefix string   `mapstructure:"config-prefix"`
	ArtifactKey  string   `mapstructure:"artifact-key"`
	Tags         []string `mapstructure:"tags"`
	Labels       []string `mapstructure:"labels"`
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, nil, raws...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	err := resolveEnvibinConfig(p.config.ConfigPrefix)
	if err != nil {
		return nil, false, err
	}

	var artifactKey = domain.NewArtifactKey(p.config.ArtifactKey)
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

	ui.Message("put artifact with key " + p.config.ArtifactKey)

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
