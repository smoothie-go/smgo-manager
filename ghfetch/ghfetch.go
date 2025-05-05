package ghfetch

import (
	"context"
	"errors"
	"runtime"
	"sort"

	"github.com/google/go-github/github"
	"github.com/smoothie-go/smgo-manager/hlog"
)

var client *github.Client = nil

func Init() {
	if client != nil {
		hlog.Fatal("Github client already initiated")
	}

	client = github.NewClient(nil)

}

func FetchReleases() ([]*github.RepositoryRelease, error) {
	ctx := context.Background()

	releases, _, err := client.Repositories.ListReleases(ctx, "smoothie-go", "smoothie-go", nil)
	if err != nil {
		hlog.Error("Cannot check for updates: Unable to request releases\n" + err.Error())
		return nil, errors.Unwrap(err)
	}

	sort.Slice(releases, func(i, j int) bool {
		t1 := releases[i].PublishedAt
		t2 := releases[j].PublishedAt
		return t1.Time.After(t2.Time)
	})

	return releases, nil
}

func FetchSmoothieGoEx(release *github.RepositoryRelease, os string, arch string, archive string) (string, error) {
	ctx := context.Background()
	assets, _, err := client.Repositories.ListReleaseAssets(ctx, "smoothie-go", "smoothie-go", *release.ID, nil)

	if err != nil {
		hlog.Error("Cannot check for updates: Unable to fetch assets\n" + err.Error())
		return "", errors.Unwrap(err)
	}

	for _, asset := range assets {
		if *asset.Name == "smoothie-go"+"-"+os+"-"+arch+archive {
			return asset.GetBrowserDownloadURL(), nil
		}
	}
	return "", errors.New("Unable to find smoothie-go for " + os + " " + arch + " " + archive)
}

func FetchSmoothieGo(release *github.RepositoryRelease) (string, error) {
	archive := ""
	switch runtime.GOOS {
	case "windows":
		archive = ".zip"
	default:
		archive = ".tar.gz"
	}
	return FetchSmoothieGoEx(release, runtime.GOOS, runtime.GOARCH, archive)
}

func FetchLatestNonPreRelease() (*github.RepositoryRelease, error) {
	releases, err := FetchReleases()
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if release.Prerelease != nil && *release.Prerelease {
			continue
		}
		return release, nil
	}

	return nil, errors.New("no non-prerelease versions found")
}

func FetchLatestRelease() (*github.RepositoryRelease, error) {
	releases, err := FetchReleases()

	return releases[0], err
}
