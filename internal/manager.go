package internal

import (
	"fmt"
	"time"

	"github.com/brandonmakai/clipmux/internal/logger"
	"github.com/brandonmakai/clipmux/persistence"
	hook "github.com/robotn/gohook"
)

var lastText string
var pasteHotkey = []string{"q", "ctrl"}

type ClipboardManager struct {
	clipIO  ReadPaster
	history *persistence.ClipboardHistory
	log     *logger.Logger
}

func NewClipboardManager(io ReadPaster, history *persistence.ClipboardHistory, log *logger.Logger) *ClipboardManager {
	return &ClipboardManager{
		clipIO:  io,
		history: history,
		log:     log,
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
		lastText = text
	}

	return nil
}

func (cm *ClipboardManager) paste() error {
    if cm.history == nil {
        fmt.Println("cm.history is nil!")
        return fmt.Errorf("clipboard history not initialized")
    }
    if cm.clipIO == nil {
        fmt.Println("cm.clipIO is nil!")
        return fmt.Errorf("clipboard IO not initialized")
    }

    item, found := cm.history.GetNewest()
    if !found {
        cm.log.Info("No items found in clipmux history.")
        return nil
    }

    if item.Data == nil {
        fmt.Println("item.Data is nil!")
        return fmt.Errorf("clipboard item is nil")
    }

    text := string(item.Data)
    err := cm.clipIO.Paste(text)

    fmt.Println("Pasted: ", text)
    if err != nil {
        cm.log.Error(err.Error())
    }
    return err
}

func (cm *ClipboardManager) Run() error {
	errCh := make(chan error)
	
	fmt.Printf("ClipIO: %v, History %v\n", cm.clipIO, cm.history)
	hook.Register(hook.KeyDown, pasteHotkey, func(e hook.Event) {
		fmt.Println("Callback started for hotkey") // NEW
    defer func() {
    if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
		fmt.Println("Hotkey pressed")
		if err := cm.paste(); err != nil {
			select { 
			case errCh <- err:
			default:
			}
		}
	})

	go func() {
		s := hook.Start()
		<- hook.Process(s)
	}()
	
	for {
		select {
		case err := <- errCh:
			return err
		default:
			if err := cm.get(); err != nil {
				return err
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

}
