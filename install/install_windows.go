//go:build windows

package install

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/smoothie-go/smgo-manager/arch"
	"github.com/smoothie-go/smgo-manager/download"
	"github.com/smoothie-go/smgo-manager/ghfetch"
	"github.com/smoothie-go/smgo-manager/hlog"
	"github.com/smoothie-go/smgo-manager/paths"
)

var links = map[string]string{
	"ffmpeg":   "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n7.1-latest-win64-gpl-7.1.zip",
	"vsbundle": "https://github.com/smoothie-go/VSBundler/releases/download/Nightly_2025.08.14_02-55/VapourSynth.zip",
}

func Package(tag string) {
	dlDir := paths.GetDownloadsDirectory()
	tmpDir := paths.GetTempDirectory()

	pkgs := []pkg{
		{"vsbundle", "VapourSynth.7z", arch.Unzip},
		{"ffmpeg", "ffmpeg.zip", arch.Unzip},
	}

	for _, p := range pkgs {
		src := filepath.Join(dlDir, p.filename)
		download.Download(links[p.name], src)
		drop := filepath.Join(tmpDir, p.name)
		if err := installPackage(src, drop, p.extractFn); err != nil {
			hlog.Error(fmt.Sprintf("%s: %v", p.name, err))
		}
	}

	smgoURL, err := ghfetch.FetchSmoothieGoFromGitTag(tag)
	if err != nil {
		hlog.Fatal(err.Error())
	}
	smgoSrc := filepath.Join(dlDir, "smoothie-go.zip")
	download.Download(smgoURL, smgoSrc)
	installPackage(smgoSrc, filepath.Join(tmpDir, "smoothie-go"), arch.Unzip)

	verDir := filepath.Join(paths.GetVersionsDirectory(), tag)
	if err := os.MkdirAll(verDir, 0755); err != nil {
		hlog.Fatal(err.Error())
	}

	copyLayout(tmpDir, verDir)

	/*err = os.RemoveAll(tmpDir)
	if err != nil {
		hlog.Error(err.Error())
		return
	}

	err = os.RemoveAll(dlDir)
	if err != nil {
		hlog.Error(err.Error())
		return
	}*/
}

func copyLayout(srcRoot, destRoot string) {
	copyDir(filepath.Join(srcRoot, "vsbundle", "VapourSynth"), destRoot)
	ffbin := filepath.Join(srcRoot, "ffmpeg", "ffmpeg-n7.1-latest-win64-gpl-7.1", "bin")
	for _, exe := range []string{"ffplay.exe", "ffmpeg.exe", "ffprobe.exe"} {
		p := filepath.Join(ffbin, exe)
		d := filepath.Join(destRoot, exe)
		copyFile(p, d)
		os.Chmod(d, 0755)
	}

	bin := fmt.Sprintf("smoothie-go-%s-%s", runtime.GOOS, runtime.GOARCH)
	copyFile(filepath.Join(srcRoot, "smoothie-go", bin), filepath.Join(destRoot, "smoothie-go.exe"))
}
