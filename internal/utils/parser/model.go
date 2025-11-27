package parser

type Parser[T any] interface {
	Parse(s string) (T, error)
}
