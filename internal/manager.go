package internal

import (
	"fmt"
	"time"

	"github.com/brandonmakai/clipmux/persistence"
	"github.com/brandonmakai/clipmux/internal/logger"
	hook "github.com/robotn/gohook"
)

var lastText string
var pasteHotkey = []string{"q", "ctrl", "i"}

type ClipboardManager struct { 
	clipIO ReadPaster
	history *persistence.ClipboardHistory
	log *logger.Logger
}

func NewClipboardManager(io ReadPaster, history *persistence.ClipboardHistory, log *logger.Logger) *ClipboardManager{
	return &ClipboardManager{
		clipIO: io, 
		history: history,
		log: log,
	}
}

// TODO: Add 'Get' by index
func (cm *ClipboardManager) get() error {
	text, err := cm.clipIO.Read()
	if err != nil {
		cm.log.Error(err.Error())
		return err
	}

	if text != lastText {
		text_bytes := []byte(text)
		cm.history.Append(text_bytes)	

		msg := fmt.Sprintf("Added new item to history: %s\n", text)
		cm.log.Info(msg)
	}
	
	return nil
}

func (cm *ClipboardManager) paste() error {
	item, found := cm.history.GetNewest()
	if !found {
		cm.log.Info("No items found in clipmux history.")
		return nil
	} 
	
	text := string(item.Data)
	err := cm.clipIO.Paste(text)

	cm.log.Error(err.Error())
	return err 
}

func (cm *ClipboardManager) Run() error {
    errCh := make(chan error)

    hook.Register(hook.KeyDown, pasteHotkey, func(e hook.Event) {
        if err := cm.paste(); err != nil {
            errCh <- err
        }
    })

    hook.Start()
    defer hook.End()

    for {
        select {
        case err := <-errCh:
						// Propogate dependency errors to main so the application can terminate 
            return err
        default:
            // Retrieve clipboard periodically
            if err := cm.get(); err != nil {
                return err
            }
            time.Sleep(100 * time.Millisecond)
        }
    }
}
