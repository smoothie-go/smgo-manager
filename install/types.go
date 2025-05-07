package install

type pkg struct {
	name      string
	filename  string
	extractFn func(src, dst string) error
}
