package repo

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
)

const (
	ScriptsPath = "../../scripts"
	ConfigPath  = "../../config"
)

func Initialize(repoRoot string) error {
	scriptBox := packr.NewBox(ScriptsPath)
	configBox := packr.NewBox(ConfigPath)

	var walkFn = func(s string, file packd.File) error {
		p := filepath.Join(repoRoot, s)
		dir := filepath.Dir(p)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}
		return ioutil.WriteFile(p, []byte(file.String()), 0644)
	}

	if err := scriptBox.Walk(walkFn); err != nil {
		return err
	}

	if err := configBox.Walk(walkFn); err != nil {
		return err
	}
	return nil
}
