Output demo.gif

Require ./t

Set Shell zsh
Set FontFamily "JetBrainsMono Nerd Font"
Set FontSize 18
Set LineHeight 1.2
Set Width 1600
Set Height 650
Set Padding 14

Set BorderRadius 8
Set Margin 14
Set WindowBar Colorful
Set WindowBarSize 40

Env XDG_CONFIG_HOME "testdata"
Env XDG_DATA_HOME "testdata"

Hide
Type "source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh" Enter
Ctrl+L
Show

Type "t 'Schedule a delivery to Omicron Persei 8' --today" Sleep 500ms Enter
Sleep 1s
Type "t" Sleep 500ms Enter
Sleep 2s
Type "l" Sleep 1s
Type "l" Sleep 2s
Type "j" Sleep 2s
Type "h" Sleep 1s
Type "h" Sleep 1s Enter
Sleep 2s
Type "j" Sleep 500ms
Type "j" Sleep 1s
Type "e" Enter
Sleep 2s
Tab
Type "Try to avoid coffee spills this time."
Sleep 2s
Ctrl+S
Sleep 2s
Ctrl+S
Sleep 5s
