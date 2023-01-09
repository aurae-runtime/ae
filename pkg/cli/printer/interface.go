package printer

import (
	"io"
)

type Interface interface {
	Format() string
	Print(w io.Writer, obj any) error
}
