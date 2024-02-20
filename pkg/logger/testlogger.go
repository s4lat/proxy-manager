package logger

import "testing"

type TestLogger struct {
	t *testing.T
}

func NewTestLogger(t *testing.T) TestLogger {
	return TestLogger{t}
}

func (t TestLogger) Debug(message interface{}, args ...interface{}) {
	args = append([]interface{}{message}, args)
	t.t.Log(args...)
}

func (t TestLogger) Info(message string, args ...interface{}) {
	t.t.Logf(message, args...)
}

func (t TestLogger) Warn(message string, args ...interface{}) {
	t.t.Logf(message, args...)
}

func (t TestLogger) Error(message interface{}, args ...interface{}) {
	args = append([]interface{}{message}, args)
	t.t.Log(args...)
}

func (t TestLogger) Fatal(message interface{}, args ...interface{}) {
	args = append([]interface{}{message}, args)
	t.t.Fatal(args...)
}
