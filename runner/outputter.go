package runner

import (
	"encoding/json"
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

type jsonResult struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  string `json:"value"`
}

func WriteJSONResult(writer io.Writer, source, value, target string) error {
	data, err := json.Marshal(jsonResult{
		Source: source,
		Target: target,
		Value:  value,
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(writer, "%s\n", data)
	return err
}
