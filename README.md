# GoMines

This is a project that I made becuase I could not find a open source local minesweeper application that has vim keybinds and customizable colors.

# Downloading and Running

Clone the repository.

`
git clone https://github.com/Robotboy26/goMines.git
`

Enter the directory.

`
cd goMines/src
`

Build the program.

`
go build
`

Run the program.

`
./goMines
`

# Customization

Customization is available through the 'settings.json' file.

Available options:
| Name | Option Type | Default | Desription |
| ---- | ----------- | ------- | ---------- |
| rows | int | 20 | The vertical height of the game field |
| cols | int | 30 | The horizontal height of the game field |
| minePersentage | int | "#00BFFF" | The persentage of squares that will have mines |
| cursorColor | string | "#FF0000" | The hex of the color for the cursor ' * ' |
| mineColor | string | "#FF0000" | The hex of the color for the mines |
| flagColor | string | "#00FFEE" | The hex of the color for the flags |
| fewAdjacentMinesColor | string | "#F6EB61" | The hex of the color for the number 1 |
| mediumAdjacentMinesColor | string | "#FF7F50" | The hex of the color for the numbers 2 and 3 |
| highAdjacentMinesColor | string | "#FF4500" | The hex of the color for the numbers 4 or greater |
| updateAdjacentOnFlag | bool | false | An assistance feature that will decrement the mine adjacency number when a flag is placed |
| autoRevealedColor | string | "#C0C0C0" | These squares can be revealed in the same way that normal sqaures can and are here to prevent an accidently placed flag from reveiling a mine and trigger a game end |

# Assistance Features

'updateAdjacentOnFlag' is an assistance features that can be used to practice how mines are positioned. This setting will decrement the mine adjacency number whenever a flag is placed.
This settings if enabled will make the game much easier.
