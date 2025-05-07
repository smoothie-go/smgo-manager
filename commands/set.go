package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

	var release *github.RepositoryRelease
	var exit bool

	for {
		var option int
		if exit {
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

	versionDir := filepath.Join(paths.GetVersionsDirectory(), release.GetTagName())
	smgoBinary := ""
	if runtime.GOOS == "windows" {
		smgoBinary = filepath.Join(versionDir, "smoothie-go.exe")
	} else {
		smgoBinary = filepath.Join(versionDir, "smoothie-go")
	}

	targetPath := filepath.Join(paths.GetManagerDirectory(), "smoothie-go")

	if _, err := os.Stat(smgoBinary); os.IsNotExist(err) {
		install.Package(release.GetTagName())
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "rmdir", targetPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			if !strings.Contains(string(output), "The system cannot find the file specified") {
				hlog.Fatal("rmdir failed: " + string(output))
			}
		}
	} else {
		if err := os.RemoveAll(targetPath); err != nil && !os.IsNotExist(err) {
			hlog.Fatal("Failed to remove: " + err.Error())
		}
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "mklink", "/J", targetPath, versionDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			hlog.Fatal("mklink failed error: " + string(output))
		}
	} else {
		err = os.Symlink(smgoBinary, targetPath)
	}

	if err != nil {
		hlog.Fatal(err.Error())
	}

	hlog.Ok("Set " + *release.TagName)
}
