package syslgo

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Fatal(args ...interface{})
	Log(args ...interface{})
}
