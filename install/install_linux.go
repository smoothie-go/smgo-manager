//go:build linux

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
	"vapoursynth":  "https://github.com/Stefan-Olt/vs-plugin-build/releases/download/vapoursynth-fa1f579a63b4fdef5f436c086875fe9c4879b5d2/vapoursynth-build-linux-x86_64.tar.gz",
	"fmtc":         "https://github.com/Stefan-Olt/vs-plugin-build/releases/download/vsplugin%2Ffmtconv%2Fgit-18a9cecb%2Flinux-glibc-x86_64%2F2024-10-10T14.44.42%2B00.00Z/fmtconv-git-18a9cecb-linux-glibc-x86_64.zip",
	"bs":           "https://github.com/Stefan-Olt/vs-plugin-build/releases/download/vsplugin%2Fcom.vapoursynth.bestsource%2FR8%2Flinux-glibc-x86_64%2F2024-12-12T18.37.59%2B00.00Z/BestSource-R8-linux-glibc-x86_64.zip",
	"svp1":         "https://github.com/smoothie-go/smoothie-go/raw/refs/heads/master/resources/vapoursynth/libsvpflow1.so",
	"svp2":         "https://github.com/smoothie-go/smoothie-go/raw/refs/heads/master/resources/vapoursynth/libsvpflow2.so",
	"frameblender": "https://github.com/couleurm/vs-frameblender/releases/download/1.2/vs-frameblender-1.2.so",
	"librife":      "https://github.com/styler00dollar/VapourSynth-RIFE-ncnn-Vulkan/releases/download/r9_mod_v32/librife_linux_x86-64.so",
	"ffmpeg":       "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n7.1-latest-linux64-gpl-7.1.tar.xz",
	"python":       "https://github.com/smoothie-go/pyenv-build/releases/download/py3104-ubuntu2204/python-3.10.4.tar.gz",
}

func Package(tag string) {
	dlDir := paths.GetDownloadsDirectory()
	tmpDir := filepath.Join(paths.GetTempDirectory(), tag)

	pkgs := []pkg{
		{"python", "python.tar.gz", arch.UntarGz},
		{"vapoursynth", "vapoursynth.tar.gz", arch.UntarGz},
		{"fmtc", "fmtconv.zip", arch.Unzip},
		{"bs", "bestsource.zip", arch.Unzip},
		{"svp1", "libsvpflow1.so", nil},
		{"svp2", "libsvpflow2.so", nil},
		{"frameblender", "frameblender.so", nil},
		{"librife", "librife.so", nil},
		{"ffmpeg", "ffmpeg.tar.xz", arch.UntarXz},
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
	smgoSrc := filepath.Join(dlDir, "smoothie-go.tar.gz")
	download.Download(smgoURL, smgoSrc)
	installPackage(smgoSrc, filepath.Join(tmpDir, "smoothie-go"), arch.UntarGz)

	verDir := filepath.Join(paths.GetVersionsDirectory(), tag)
	if err := os.MkdirAll(verDir, 0755); err != nil {
		hlog.Fatal(err.Error())
	}

	copyLayout(tmpDir, verDir)
	generateLaunchScript(verDir)

	err = os.RemoveAll(tmpDir)
	if err != nil {
		hlog.Fatal(err.Error())
	}

	err = os.RemoveAll(dlDir)
	if err != nil {
		hlog.Fatal(err.Error())
	}
}

func copyLayout(srcRoot, destRoot string) {
	copyDir(filepath.Join(srcRoot, "python", "3.10.4", "bin"), destRoot)
	copyDir(filepath.Join(srcRoot, "python", "3.10.4", "lib"), filepath.Join(destRoot, "lib"))
	copyDir(filepath.Join(srcRoot, "vapoursynth", "workspace", "lib"), filepath.Join(destRoot, "lib"))
	copyFile(filepath.Join(srcRoot, "vapoursynth", "workspace", "bin", "vspipe"), filepath.Join(destRoot, "vspipe"))

	ffbin := filepath.Join(srcRoot, "ffmpeg", "ffmpeg-n7.1-latest-linux64-gpl-7.1", "bin")
	for _, exe := range []string{"ffplay", "ffmpeg", "ffprobe"} {
		p := filepath.Join(ffbin, exe)
		d := filepath.Join(destRoot, exe)
		copyFile(p, d)
		os.Chmod(d, 0755)
	}

	for key, fn := range map[string]string{
		"bs":           "bestsource.so",
		"fmtc":         "libfmtconv.so",
		"frameblender": "frameblender.so",
		"svp1":         "libsvpflow1.so",
		"svp2":         "libsvpflow2.so",
		"librife":      "librife.so",
	} {
		src := filepath.Join(srcRoot, key, fn)
		dst := filepath.Join(destRoot, "lib", "vapoursynth", fn)
		copyFile(src, dst)
	}

	bin := fmt.Sprintf("smoothie-go-%s-%s", runtime.GOOS, runtime.GOARCH)
	copyFile(filepath.Join(srcRoot, "smoothie-go", bin), filepath.Join(destRoot, bin))
}

func generateLaunchScript(verDir string) {
	libDir := filepath.Join(verDir, "lib")
	py := filepath.Join(libDir, "python3.10", "site-packages")
	conf := filepath.Join(verDir, "vapoursynth.conf")
	bin := filepath.Join(verDir, fmt.Sprintf("smoothie-go-%s-%s", runtime.GOOS, runtime.GOARCH))

	script := fmt.Sprintf(`#!/bin/bash
exec env LD_LIBRARY_PATH="%s" PATH="%s" PYTHONPATH="%s" PYTHONHOME="%s" VAPOURSYNTH_CONF_PATH="%s" "%s" "$@"
`, libDir, verDir, py, verDir, conf, bin)

	if err := os.WriteFile(filepath.Join(verDir, "smoothie-go"), []byte(script), 0755); err != nil {
		hlog.Fatal(err.Error())
	}
	if err := os.WriteFile(conf, []byte("SystemPluginDir="+filepath.Join(libDir, "vapoursynth")+"\n"), 0644); err != nil {
		hlog.Fatal(err.Error())
	}
}
