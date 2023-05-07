package apcgen

type RuneBuildOptions struct {
	SkipWhitespace bool
}

var DefaultRuneBuildOptions = RuneBuildOptions{
	SkipWhitespace: true,
}
