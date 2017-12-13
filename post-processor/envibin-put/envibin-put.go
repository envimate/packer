package envibinPut

import (
	"fmt"

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

	err = c.Push(artifactKey, compressed(artifact.Files()), p.config.Tags, p.config.Labels)
	if err != nil {
		return nil, false, err
	}

	ui.Message("put artifact with key " + p.config.ArtifactKey)

	return nil, false, nil
}

func compressed(files []string) string {
	if len(files) > 0 {
		return files[0]
	}

	return ""
}
