package util

type ErrorReporter interface {
	Error(interface{}, string, ...interface{})
}
