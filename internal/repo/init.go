package repo

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
)

const (
	ScriptsPath = "../../scripts"
	ConfigPath  = "../../config"
)

func Initialize(repoRoot string) error {
	scriptBox := packr.New(ScriptsPath, ScriptsPath)
	configBox := packr.New(ConfigPath, ConfigPath)

	var walkFn = func(s string, file packd.File) error {
		p := filepath.Join(repoRoot, s)
		dir := filepath.Dir(p)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				return err
			}
		}
		return ioutil.WriteFile(p, []byte(file.String()), 0777)
	}

	if err := scriptBox.Walk(walkFn); err != nil {
		return err
	}

	if err := configBox.Walk(walkFn); err != nil {
		return err
	}
	return nil
}
