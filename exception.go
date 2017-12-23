package e

/*
	Package `e` provides a simple high-level way of error handling that may
	be used in higher-order packages, where classic go-way error handling turns into
	`if err != nil { return nil, err }` mantra. It's not intended to be used in a library code.
 */

type Exception struct {
	SmartError
}

func (e Exception) Error() string {
	return e.SmartError.Error()
}

func Catch(handler func(e *Exception)) {
	if r := recover(); r != nil {
		handle(r, handler)
	}
}

func OnError(handler func(e *Exception)) {
	if r := recover(); r != nil {
		handle(r, handler)
		panic(r)
	}
}

func Throw(info string, cause error, args ...Map) {
	exception := &Exception{wrap(info, cause, args...)}
	ExecHooks(exception)
	panic(exception)
}

func Must(err error, info string, args ...Map) {
	if err != nil {
		Throw(info, err, args...)
	}
}

func handle(r interface{}, handler func(e *Exception)) {
	var exception *Exception
	switch e := r.(type) {
	case *Exception:
		exception = e
	case error:
		exception = &Exception{wrap("third party exception", e)}
	default:
		exception = &Exception{wrap("third party exception", nil, Map{
			"recovered": e,
		})}
	}

	handler(exception)
}
