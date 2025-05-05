//go:build linux

package install

var links map[string]string = map[string]string{
	"vapoursynth": "https://github.com/Stefan-Olt/vs-plugin-build/releases/download/vapoursynth-fa1f579a63b4fdef5f436c086875fe9c4879b5d2/vapoursynth-build-linux-x86_64.tar.gz",
	"fmtc":        `https://github.com/Stefan-Olt/vs-plugin-build/releases/download/vsplugin%2Ffmtconv%2Fgit-18a9cecb%2Flinux-glibc-x86_64%2F2024-10-10T14.44.42%2B00.00Z/fmtconv-git-18a9cecb-linux-glibc-x86_64.zip`,
	"bs":          `https://github.com/Stefan-Olt/vs-plugin-build/releases/download/vsplugin%2Fcom.vapoursynth.bestsource%2FR8%2Flinux-glibc-x86_64%2F2024-12-12T18.37.59%2B00.00Z/BestSource-R8-linux-glibc-x86_64.zip`,
	"svp1":        "https://github.com/smoothie-go/smoothie-go/raw/refs/heads/master/resources/vapoursynth/libsvpflow1.so",
}

func Package() {

}
