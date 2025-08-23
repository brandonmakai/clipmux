package internal

import (
	"fmt"
	"time"
	"strconv"

	"github.com/brandonmakai/clipmux/internal/logger"
	"github.com/brandonmakai/clipmux/persistence"
	hook "github.com/robotn/gohook"
)

// TODO: Make these a config option
const pasteHotkeyBase = "q"
const maxHotkeys = 10
const linearHistory = true

type ClipboardManager struct {
	clipIO  ReadPaster
	history *persistence.ClipboardHistory
	log     *logger.Logger
	lastText string
}

func NewClipboardManager(io ReadPaster, history *persistence.ClipboardHistory, log *logger.Logger) *ClipboardManager {
	return &ClipboardManager{
		clipIO:  io,
		history: history,
		log:     log,
		lastText: "",
	}
}

func (cm *ClipboardManager) get() error {
	text, err := cm.clipIO.Read()
	if err != nil {
		cm.log.Error(err.Error())
		return err
	}

	if text != cm.lastText || !cm.history.Contains(text) {
		fmt.Println("Appending Text: ", text)
		text_bytes := []byte(text)
		cm.history.Append(text_bytes)

		msg := fmt.Sprintf("Added new item to history: %s\n", text)
		cm.log.Info(msg)
		cm.lastText = text
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
    cm.clipIO.Paste(text)

    fmt.Println("Pasted: ", text)
    return nil 
}

func (cm *ClipboardManager) Run() error {
	errCh := make(chan error)
	
	for i := 0; i < maxHotkeys; i++ {
	  pos := i // shadow the loop variable to prevent callback from only getting final value
		hotkey := append([]string{pasteHotkeyBase, "ctrl"}, strconv.Itoa(pos))

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
