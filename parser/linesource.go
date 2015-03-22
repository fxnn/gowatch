package parser

type LineSource interface {
	Lines() <-chan string
}
