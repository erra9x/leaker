package runner

import (
	"fmt"
	"io"
)

func WritePlainResult(writer io.Writer, verbose bool, source, value string) error {
	if verbose {
		_, err := fmt.Fprintf(writer, "[%s] %s\n", source, value)
		return err
	}
	_, err := fmt.Fprintf(writer, "%s\n", value)
	return err
}
