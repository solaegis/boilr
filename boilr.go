package main

import (
	"fmt"

	"github.com/solaegis/boilr/pkg/boilr"
	"github.com/solaegis/boilr/pkg/cmd"
	"github.com/solaegis/boilr/pkg/util/exit"
	"github.com/solaegis/boilr/pkg/util/osutil"
)

func main() {
	if exists, err := osutil.DirExists(boilr.Configuration.TemplateDirPath); !exists {
		if err := osutil.CreateDirs(boilr.Configuration.TemplateDirPath); err != nil {
			exit.Error(fmt.Errorf("Tried to initialise your template directory, but it has failed: %s", err))
		}
	} else if err != nil {
		exit.Error(fmt.Errorf("Failed to init: %s", err))
	}

	cmd.Run()
}
