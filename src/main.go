package main

import (
    "fmt"
    "math/rand"
    "encoding/json"
    "os"
    "log"
    "io"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell/v2"
)

type Settings struct {
    Rows int `json:"rows"`
    Cols int `json:"cols"`
    MinePercentage int `json:"minePercentage"`
    numMines int
    // Colors
    CursorColor string `json:"cursorColor"`
    MineColor string `json:"mineColor"`
    FlagColor string `json:"flagColor"`
    FewAdjacentMinesColor string `json:"fewAdjacentMinesColor"`
    MediumAdjacentMinesColor string `json:"mediumAdjacentMinesColor"`
    HighAdjacentMinesColor string `json:"highAdjacentMinesColor"`
    UpdateAdjacentOnFlag bool `json:"updateAdjacentOnFlag"`
    AutoRevealedColor string `json:"autoRevealedColor"`
}

type Cell struct {
    isMine     bool
    isRevealed bool
    isAutoRevealed bool
    isFlagged  bool
    adjacent   int
}

type Minesweeper struct {
    settings    *Settings
    grid        [][]Cell
    gameOver    bool
    win         bool
    cursorX     int
    cursorY     int
    flags       int
}

func NewMinesweeper(settings *Settings) *Minesweeper {
    ms := &Minesweeper{
        grid: make([][]Cell, settings.Rows),
    }
    for i := range ms.grid {
        ms.grid[i] = make([]Cell, settings.Cols)
    }
    ms.settings = settings
    ms.settings.numMines = int(float64(settings.Rows * settings.Cols) * (float64(settings.MinePercentage) / 100))
    ms.placeMines()
    ms.flags = settings.numMines
    ms.calculateAdjacent()
    return ms
}

func (ms *Minesweeper) placeMines() {
    for i := 0; i < ms.settings.numMines; {
        x := rand.Intn(ms.settings.Rows)
        y := rand.Intn(ms.settings.Cols)
        if !ms.grid[x][y].isMine {
            ms.grid[x][y].isMine = true
            i++
        }
    }
}

func (ms *Minesweeper) calculateAdjacent() {
    for i := range ms.grid {
        for j := range ms.grid[i] {
            if ms.grid[i][j].isMine {
                continue
            }
            for x := -1; x <= 1; x++ {
                for y := -1; y <= 1; y++ {
                    if x == 0 && y == 0 {
                        continue
                    }
                    nx, ny := i+x, j+y
                    if nx >= 0 && nx < ms.settings.Rows && ny >= 0 && ny < ms.settings.Cols && ms.grid[nx][ny].isMine {
                        ms.grid[i][j].adjacent++
                    }
                }
            }
        }
    }
}

func (ms *Minesweeper) checkWinCondition() bool {
    revealedCells := 0
    for i := range ms.grid {
        for j := range ms.grid[i] {
            // Count revealed cells that are not mines
            if ms.grid[i][j].isRevealed && !ms.grid[i][j].isMine {
                revealedCells++
            }
        }
    }
    // Check if the number of revealed cells is equal to total cells minus mines
    return revealedCells == (ms.settings.Rows * ms.settings.Cols - ms.settings.numMines)
}

func (ms *Minesweeper) revealCell(x, y int) {
    if x < 0 || x >= ms.settings.Rows || y < 0 || y >= ms.settings.Cols || ms.grid[x][y].isFlagged {
        return
    }

    if ms.grid[x][y].isRevealed && !ms.grid[x][y].isAutoRevealed {
        return
    } else if ms.grid[x][y].isRevealed && ms.grid[x][y].isAutoRevealed {
        ms.grid[x][y].isAutoRevealed = false
    }

    ms.grid[x][y].isRevealed = true
    if ms.grid[x][y].isMine {
        ms.gameOver = true
        return
    }
    if ms.grid[x][y].adjacent == 0 {
        for nx := x - 1; nx <= x+1; nx++ {
            for ny := y - 1; ny <= y+1; ny++ {
                ms.revealCell(nx, ny)
            }
        }
    }
    
    if ms.checkWinCondition() {
        ms.win = true
    }
}

// func (ms *Minesweeper) toggleFlag(x, y int) {
//     if x < 0 || x >= ms.settings.Rows || y < 0 || y >= ms.settings.Cols {
//         return
//     }
//
//     if ms.grid[x][y].isRevealed {
//         return
//     }
//
//     if ms.grid[x][y].isFlagged {
//         ms.grid[x][y].isFlagged = false
//         ms.flags = ms.flags + 1
//     } else {
//         if ms.flags > 0 {
//             ms.grid[x][y].isFlagged = true
//             ms.flags = ms.flags - 1
//         }
//     }
// }

func (ms *Minesweeper) toggleFlag(x, y int) {
    if x < 0 || x >= ms.settings.Rows || y < 0 || y >= ms.settings.Cols {
        return
    }

    if ms.grid[x][y].isRevealed {
        return
    }

    // Function to change the adjacent mine counts around a cell
    // Not currently working as it will show the mines because there is not a guarantee that you know where the flags are being placed
    adjustAdjacentCount := func(x, y int, delta int) {
        for dx := -1; dx <= 1; dx++ {
            for dy := -1; dy <= 1; dy++ {
                if dx == 0 && dy == 0 {
                    continue
                }
                nx, ny := x+dx, y+dy
                if nx >= 0 && nx < ms.settings.Rows && ny >= 0 && ny < ms.settings.Cols {
                    ms.grid[nx][ny].adjacent += delta
                    if ms.grid[nx][ny].adjacent != 0 && ms.grid[nx][ny].isAutoRevealed {
                        ms.grid[nx][ny].isRevealed = false
                        ms.grid[nx][ny].isAutoRevealed = false
                    } else if ms.grid[nx][ny].adjacent == 0 && !ms.grid[nx][ny].isRevealed && !ms.grid[nx][ny].isMine {
                        // Don't reveal because you might be killed a mine if you accidently put a flag in the wrong spot
                        // ms.revealCell(nx, ny) // Reveal the cells that are now empty after a flag
                        // Instead use another type that will give another color maybe a silver
                        ms.grid[nx][ny].isRevealed = true
                        ms.grid[nx][ny].isAutoRevealed = true
                    } else if ms.grid[nx][ny].isMine {
                        ms.grid[nx][ny].isAutoRevealed = true
                    }
                }
            }
        }
    }

    if ms.grid[x][y].isFlagged {
        ms.grid[x][y].isFlagged = false
        ms.flags = ms.flags + 1
        // Adjust adjacent counts if the flag adjustment is enabled
        if ms.settings.UpdateAdjacentOnFlag {
            adjustAdjacentCount(x, y, 1) // Increment the counts
        }
    } else {
        if ms.flags > 0 {
            ms.grid[x][y].isFlagged = true
            ms.flags = ms.flags - 1
            // Adjust adjacent counts if the flag adjustment is enabled
            if ms.settings.UpdateAdjacentOnFlag {
                adjustAdjacentCount(x, y, -1) // Decrement the counts
            }
        }
    }
}

func (ms *Minesweeper) display() string {
    display := ""
    for i := 0; i < ms.settings.Rows; i++ {
        for j := 0; j < ms.settings.Cols; j++ {
            // Check if the cursor is at the current cell's position
            if ms.cursorX == i && ms.cursorY == j {
                display += fmt.Sprintf(" [%s]*[white] ", ms.settings.CursorColor) // Cursor indication
                continue
            }

            cell := ms.grid[i][j]
            if cell.isFlagged {
                display += fmt.Sprintf(" [%s]F[white] ", ms.settings.FlagColor)
            } else if !cell.isRevealed {
                display += fmt.Sprintf("░░░") // Just keep white
            } else if cell.isMine {
                display += fmt.Sprintf("[%s]💣💣💣[white]", ms.settings.MineColor)
            } else if cell.isAutoRevealed {
                display += fmt.Sprintf("[%s]░░░[white]", ms.settings.AutoRevealedColor)
            } else {
                // Minimal color coding based on the number of adjacent mines
                switch {
                case cell.adjacent == 0:
                    display += "   " // Three spaces without any number
                case cell.adjacent == 1:
                    display += fmt.Sprintf(" [%s]%d[white] ", ms.settings.FewAdjacentMinesColor, cell.adjacent)
                case cell.adjacent == 2:
                    display += fmt.Sprintf(" [%s]%d[white] ", ms.settings.MediumAdjacentMinesColor, cell.adjacent)
                case cell.adjacent == 3:
                    display += fmt.Sprintf(" [%s]%d[white] ", ms.settings.MediumAdjacentMinesColor, cell.adjacent)
                case cell.adjacent >= 4:
                    display += fmt.Sprintf(" [%s]%d[white] ", ms.settings.HighAdjacentMinesColor, cell.adjacent)
                }
            }
        }
        display += "\n" // Get to the next row
    }
    return display
}

func main() {
    // Load json settings
    settingsFile, err := os.Open("settings.json")
    if err != nil {
        log.Fatal(err)
    }
    defer settingsFile.Close()
    settingsFileData, err := io.ReadAll(settingsFile)
    if err != nil {
        log.Fatal(err)
    }
    var settings Settings
    err = json.Unmarshal(settingsFileData, &settings)
    if err != nil {
        log.Fatal(err)
    }

    var display string
    ms := NewMinesweeper(&settings)
    app := tview.NewApplication()

    // Initial screen setup
    display = ms.display()
    display += "Press any key to start!"

    textView := tview.NewTextView().
        SetText(display).
        SetTextAlign(tview.AlignLeft).
        SetDynamicColors(true)

        textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if ms.gameOver {
            app.Stop()
            return nil
        }

        switch event.Key() {
        case tcell.KeyEscape:
            app.Stop()
            return nil
        case tcell.KeyUp:
            if ms.cursorX > 0 {
                ms.cursorX--
            }
        case tcell.KeyDown:
            if ms.cursorX < settings.Rows - 1 {
                ms.cursorX++
            }
        case tcell.KeyLeft:
            if ms.cursorY > 0 {
                ms.cursorY--
            }
        case tcell.KeyRight:
            if ms.cursorY < settings.Cols - 1 {
                ms.cursorY++
            }
        case tcell.KeyEnter:
            ms.revealCell(ms.cursorX, ms.cursorY)
        case tcell.KeyRune: // Handle single character input
            switch event.Rune() {
            case 'h':
                if ms.cursorY > 0 {
                    ms.cursorY--
                }
            case 'j':
                if ms.cursorX < settings.Rows - 1 {
                    ms.cursorX++
                }
            case 'k':
                if ms.cursorX > 0 {
                    ms.cursorX--
                }
            case 'l':
                if ms.cursorY < settings.Cols - 1 {
                    ms.cursorY++
                }
            case 'f':
                ms.toggleFlag(ms.cursorX, ms.cursorY)
            case 'r':
                if !ms.grid[ms.cursorX][ms.cursorY].isFlagged {
                    ms.revealCell(ms.cursorX, ms.cursorY)
                }
            case ' ':
                if !ms.grid[ms.cursorX][ms.cursorY].isFlagged {
                    ms.revealCell(ms.cursorX, ms.cursorY)
                }
            }
        }

        display = ms.display()
        display += "Use h j k l to move (Vim keybinds)\n"
        display += fmt.Sprintf("Flags remaining: %d\n", ms.flags)
        if ms.gameOver {
            display += "[red] Game Over"
        }
        if ms.win {
            display += "[green] Game Complete"
        }

        textView.SetText(display)
        return nil
    })

    if err := app.SetRoot(textView, true).Run(); err != nil {
        panic(err)
    }
}
