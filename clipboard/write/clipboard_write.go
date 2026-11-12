package clipboard

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.design/x/clipboard"
)

type ClipboardContent struct {
	DataType int
	Bytes    []byte
}

func ClipboardWrite(contents string) {
	var clipboardFormatted ClipboardContent
	err := json.Unmarshal([]byte(contents), &clipboardFormatted)
	if err != nil {
		fmt.Println(fmt.Errorf("Clipboard translation error: %v", err))
	} else {
		_, err := clipboard.Write(context.Background(), clipboard.Format(clipboardFormatted.DataType), clipboardFormatted.Bytes)
		if err != nil {
			fmt.Println(err)
		}
	}
}
