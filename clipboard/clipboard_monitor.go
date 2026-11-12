package clipboard

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/benn0x/clipboard-sync/networking/ws"
	"github.com/benn0x/clipboard-sync/security"
	"golang.design/x/clipboard"
)

type ClipboardContent struct {
	DataType int
	Bytes    []byte
}

func StartClipboardMonitor() {
	ctx := context.Background()

	for data := range clipboard.Watch(ctx) {
		clipboardFormatted := ClipboardContent{
			DataType: int(data.Format),
			Bytes:    data.Bytes,
		}

		bytes, err := json.Marshal(clipboardFormatted)
		fmt.Println(string(bytes), clipboardFormatted)
		if err != nil {
			fmt.Println(fmt.Errorf("Couldn't process clipboard data: %v", err))
		} else {
			security.Encrypt <- string(bytes)
			encryptedClipboard := <-security.Encrypt

			ws.SendOverWS <- encryptedClipboard
		}
	}
}
