package main

import "io"

type Serializable interface {
	Serialize(w io.Writer)
	Desesrialize(r io.Reader)
}
