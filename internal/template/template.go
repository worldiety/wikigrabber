package template

import (
	"github.com/worldiety/wikigrabber/internal/config"
	"github.com/worldiety/wikigrabber/internal/git"
	"path/filepath"
	"strings"
)

// Provide ensures the availability of the template
func Provide(cfg config.Config) (string, error) {
	sanitizedName := cfg.Template
	if strings.HasPrefix(cfg.Template, "https") {
		sanitizedName = sanitizedName[len("https://")-1:]
		myPath := filepath.Join(cfg.TmpDir, sanitizedName)
		if err := git.CloneOrPull(cfg.Template, myPath); err != nil {
			return "", err
		}
		return myPath, nil
	}

	return cfg.Template, nil
}
