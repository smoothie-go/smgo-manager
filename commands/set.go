package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/smoothie-go/smgo-manager/ghfetch"
	"github.com/smoothie-go/smgo-manager/hlog"
	"github.com/smoothie-go/smgo-manager/install"
	"github.com/smoothie-go/smgo-manager/paths"
)

func Set() {
	ghfetch.Init()
	releases, err := ghfetch.FetchReleases()

	if err != nil {
		hlog.Error(err.Error())
	}

	for i := len(releases) - 1; i >= 0; i-- {
		fmt.Printf("(%d) smoothie-go %s\n", i+1, releases[i].GetName())
	}

	var release *github.RepositoryRelease = nil

	var exit bool

	for {
		var option int
		if exit == true {
			break
		}
		fmt.Printf("Select a version: ")
		fmt.Scanf("%d", &option)

		if option > len(releases) || option < 1 {
			hlog.Error("Please enter a valid release")
		} else {
			release = releases[option-1]
			exit = true
		}
	}

	smgo_script_path := filepath.Join(paths.GetVersionsDirectory(), release.GetTagName(), "smoothie-go")
	if _, err := os.Stat(smgo_script_path); os.IsNotExist(err) {
		install.Package(release.GetTagName())
	}

	err = os.Remove(filepath.Join(paths.GetManagerDirectory(), "smoothie-go"))
	if os.IsNotExist(err) {
	} else if err != nil {
		hlog.Fatal(err.Error())
	}

	err = os.Symlink(smgo_script_path, filepath.Join(paths.GetManagerDirectory(), "smoothie-go"))
	if err != nil {
		hlog.Fatal(err.Error())
	}

	hlog.Ok("Set " + *release.TagName)
}
