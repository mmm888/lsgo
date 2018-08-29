package format

import "os"

type Format interface {
	Execute([]os.FileInfo) error
}
