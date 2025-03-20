## Overview
This command-line tool reads a tab-indented list from stdin and displays it as a collapsible tree in your terminal. You can move the cursor to different items and expand or collapse their children.

## install 

```bash
go install github.com/kungfusheep/colaps@latest
```

## Usage
1. Pipe any tab-indented text into the program:
   ```bash
   cat my-list.txt | colaps
   ```
2. Use the arrow keys (or `j/k` for down/up) to move the cursor.
3. Press `l` or `Tab` to expand or collapse the currently selected item.
4. Press `left` (or `h`) to collapse the current item or jump back to its parent.
5. Press `q` or `Ctrl+C` to quit.

