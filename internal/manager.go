package internal

import (
	"fmt"
	"time"
	"strconv"

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

func (cm *ClipboardManager) get() error {
	text, err := cm.clipIO.Read()
	if err != nil {
		cm.log.Error(err.Error())
		return err
	}

	if text != lastText || !cm.history.Contains(text) {
		fmt.Println("Appending Text: ", text)
		text_bytes := []byte(text)
		cm.history.Append(text_bytes)

		msg := fmt.Sprintf("Added new item to history: %s\n", text)
		cm.log.Info(msg)
		lastText = text
	}

	return nil
}

func (cm *ClipboardManager) paste(idx int) error {
		var item persistence.Item
		if cm.history == nil {
        fmt.Println("cm.history is nil!")
        return fmt.Errorf("clipboard history not initialized")
    }
    if cm.clipIO == nil {
        fmt.Println("cm.clipIO is nil!")
        return fmt.Errorf("clipboard IO not initialized")
    }
		
		if idx == cm.history.Newest() {
			var found = false
			item, found = cm.history.GetNewest()

			if !found {
        cm.log.Info("No items found in clipmux history.")
        return nil
	    }
		} else {
			var found = false 
			item, found = cm.history.GetPos(idx)

			if !found {
				cm.log.Info(fmt.Sprintf("Item not found at index: %v", idx))
			}
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
	
	for pos := range 10 {
		hotkey := append([]string(nil), pasteHotkey...)
		hotkey = append(hotkey, strconv.Itoa(pos))

		hook.Register(hook.KeyDown, hotkey, func(e hook.Event) {
			fmt.Println("Callback started for hotkey index: ", pos) // NEW
			fmt.Println("Hotkey pressed")
			if err := cm.paste(pos); err != nil {
				select { 
				case errCh <- err:
				default:
				}
			}
		})
	}

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
