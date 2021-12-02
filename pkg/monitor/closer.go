package monitor

import "io"

var _closerList = []io.Closer{}

func appendClose(closer io.Closer) {
	_closerList = append(_closerList, closer)
}

func Close() {
	for _, closer := range _closerList {
		closer.Close()
	}
}
