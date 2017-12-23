package e

import (
	"fmt"
)

type Map map[string]interface{}

type SmartError struct {
	Info  string
	Cause error
	Args  Map
}

func (e SmartError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Info, e.Cause)
}

func (e SmartError) BottommostError() error {
	var err error = e
	for {
		if se, ok := err.(SmartError); ok {
			err = se.Cause
		} else {
			break
		}
	}
	return err
}

func WrapError(info string, cause error, args ...Map) error {
	return wrap(info, cause, args...)
}

func wrap(info string, cause error, args ...Map) SmartError {
	var argsMap Map
	switch len(args) {
	case 0:
		// Leave as nil
	case 1:
		argsMap = args[0]
	default:
		argsMap = make(Map)
		for _, m := range args {
			for k, v := range m {
				argsMap[k] = v
			}
		}
	}
	return SmartError{info, cause, argsMap}
}

func Is(err1, err2 error) bool {
	if err1 == err2 {
		return true
	}

	e, ok := err1.(SmartError)
	if ok {
		return Is(e.Cause, err2)
	}

	return false
}
