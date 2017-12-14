package envibinPut

import (
	"fmt"

	"errors"

	"bitbucket.org/envimate/envibin-cli/client"
	"bitbucket.org/envimate/envibin-cli/domain"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
)

type Config struct {
	Url         string   `mapstructure:"url"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
	ArtifactKey string   `mapstructure:"artifact-key"`
	Tags        []string `mapstructure:"tags"`
	Labels      []string `mapstructure:"labels"`
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
	if p.config.Url == "" {
		p.config.Url = "http://localhost:9000"
		ui.Message("no repository url configured, using default " + p.config.Url)
	}

	var artifactKey = domain.NewArtifactKey(p.config.ArtifactKey)
	if artifactKey.IsInvalid() {
		return nil, false, fmt.Errorf("invalid artifact-key %s", artifactKey.String())
	}

	c, err := client.New(p.config.Url)
	if err != nil {
		return nil, false, err
	}

	compressedFile, err := compressed(artifact.Files())
	if err != nil {
		return nil, false, err
	}

	err = c.Push(artifactKey, compressedFile, p.config.Tags, p.config.Labels)
	if err != nil {
		return nil, false, err
	}

	ui.Message("put artifact with key " + p.config.ArtifactKey)

	return nil, false, nil
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
