# ClipMux
A lightweight clipboard multiplexer for power users.

Access your clipboard history via global hotkeys, just like tmux for your clipboard.

## Planned Features
* Tracks the last N clipboard entries (default: 10).
* Access items directly via numeric hotkeys (Ctrl+Alt+1 → first item, Ctrl+Alt+2 → second item).
* Deduplicates clipboard entries automatically.
* Persists clipboard history across sessions.
* Minimal design – no menus or distractions, just direct clipboard access.
* macOS supported

## Planned Usage
Run ClipMux
`./clipmux`
The program will run in the background, listening for clipboard changes and global hotkeys.

### Hotkeys
* Ctrl+Alt+1 → Paste first clipboard item
* Ctrl+Alt+2 → Paste second clipboard item
* Ctrl+Alt+N → Paste N-th clipboard item

### Configuration
History size, hotkeys, and persistence path will be configurable in future releases.

Currently uses JSON file in $HOME/.clipmux/history.json to store clipboard history.

### Contributing
ClipMux is open source and welcomes contributions!
1. Fork the repo
2. Create a feature branch (git checkout -b feature-name)
3. Commit your changes (git commit -m 'Add feature')
4. Push to branch (git push origin feature-name)
5. Open a Pull Request

### Roadmap
1. MVP: Clipboard history tracking + hotkeys (WOP)
2. Customer Bindings 
3. Terminal CLI Commands
4. Cross-platform support and testing
5. Optional encryption for sensitive clipboard items

### License
MIT License © Brandon Williams

