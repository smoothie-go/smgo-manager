package main

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/smoothie-go/smgo-manager/ghfetch"
	"github.com/smoothie-go/smgo-manager/hlog"
	"github.com/smoothie-go/smgo-manager/install"
)

func main() {
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

	install.Package(release.GetTagName())
}
