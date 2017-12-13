package envibinPut

import (
	"log"

	"github.com/hashicorp/packer/packer"
)

type PostProcessor struct{}

// Entry point for configuration parsing when we've defined
func (p *PostProcessor) Configure(raws ...interface{}) error {
	log.Println("Configuring envibin-put")
	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	ui.Message("Post processing envibin-put")
	return nil, false, nil
}
