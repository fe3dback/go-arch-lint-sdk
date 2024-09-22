package tpl

//go:generate ../../bin/mockgen -source=$GOFILE -destination=mocks/mocks.go -package=mocks

type asciiColorizer interface {
	Colorize(color string, text string) string
}
