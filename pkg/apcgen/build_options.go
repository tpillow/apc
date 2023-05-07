package apcgen

type BuildOptions struct {
	SkipWhitespace bool
}

var DefaultBuildOptions = BuildOptions{
	SkipWhitespace: true,
}
