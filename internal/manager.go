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
var pasteHotkeyBase = []string{"ctrl", "shift", "h"}
const maxHotkeys = 10
const linearHistory = true

type ClipboardManager struct {
	clipIO  ReadPaster
	history persistence.ClipboardHistory
	log     *logger.Logger
	lastText string
}

func NewClipboardManager(io ReadPaster, history persistence.ClipboardHistory, log *logger.Logger) *ClipboardManager {
	return &ClipboardManager{
		clipIO:  io,
		history: history,
		log:     log,
		lastText: "",
	}
}

func (cm *ClipboardManager) get(linearHistory bool) error {
	text, err := cm.clipIO.Read()
	if err != nil {
		cm.log.Error(err.Error())
		return err
	}

	if linearHistory {
		if text != cm.lastText && !cm.history.Contains(text) {
			cm.appendToHistory(text)
		}
	} else {
		if text != cm.lastText {
			cm.appendToHistory(text)
		}
	}

	return nil
}

func (cm *ClipboardManager) appendToHistory(text string) {
	fmt.Println("Appending Text: ", text)
	cm.history.Append([]byte(text))
	cm.log.Info(fmt.Sprintf("Added new item to history: %s\n", text))
	cm.lastText = text
}

func (cm *ClipboardManager) paste(idx int) error {
	if cm.history == nil {
		return fmt.Errorf("clipboard history not initialized")
	}
	if cm.clipIO == nil {
		return fmt.Errorf("clipboard IO not initialized")
	}

	newest := false
	if idx < 0 {
		idx = cm.history.NewestIndex()
		newest = true
	}

	item, err := cm.retrieveItem(idx, newest)
	if err != nil {
		return err
	}	

	text := string(item.Data)
	cm.clipIO.Paste(text)

	fmt.Println("Pasted:", text)
	return nil
}

func (cm *ClipboardManager) retrieveItem(idx int, newest bool) (persistence.Item, error) {
	var item persistence.Item
	var found bool

	if newest {
		item, found = cm.history.Newest()
	} else {
		item, found = cm.history.GetPos(idx)
	}

	if !found || item.Data == nil {
		cm.log.Info(fmt.Sprintf("Item not found at index: %v", idx))
		return persistence.Item{}, fmt.Errorf("no item found at index %v", idx)
	}

	return item, nil
}

func (cm *ClipboardManager) Run() error {
	errCh := make(chan error)
	
	for i := 0; i < maxHotkeys; i++ {
	  pos := i // shadow the loop variable to prevent callback from only getting final value
		hotkey := append(pasteHotkeyBase, strconv.Itoa(pos))

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
		evChan := hook.Start()
		<- hook.Process(evChan)
	}()
	
	for {
		select {
		case err := <- errCh:
			return err
		default:
			if err := cm.get(linearHistory); err != nil {
				return err
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

}
