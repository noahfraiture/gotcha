package image

import (
	"io"
)

type Printer interface {
	Fprint(w io.Writer) error
	Fprintln(w io.Writer) error
}
