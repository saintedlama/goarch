package subpkg

import "fmt"

func subErr() error {
	return fmt.Errorf("sub error")
}
