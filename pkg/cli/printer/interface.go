package printer

import (
	"io"
)

type Interface interface {
	Format() string
	Print(obj any, w io.Writer) error
}
