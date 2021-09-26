package reporter

// Color represents a terminal color.
// These values are specifically used to define terminal escape characters having
// the same width, so they place nice with tabwriter.
// Approach is taken from Dave (https://stackoverflow.com/users/295163/dave)
// For more details, see https://stackoverflow.com/a/46208644/151445
type Color string

//goland:noinspection GoUnusedConst
const (
	Reset                   Color = "\x1b[0000m"
	Bright                  Color = "\x1b[0001m"
	BlackText               Color = "\x1b[0030m"
	RedText                 Color = "\x1b[0031m"
	GreenText               Color = "\x1b[0032m"
	YellowText              Color = "\x1b[0033m"
	BlueText                Color = "\x1b[0034m"
	MagentaText             Color = "\x1b[0035m"
	CyanText                Color = "\x1b[0036m"
	WhiteText               Color = "\x1b[0037m"
	DefaultText             Color = "\x1b[0039m"
	BrightRedText           Color = "\x1b[1;31m"
	BrightGreenText         Color = "\x1b[1;32m"
	BrightYellowText        Color = "\x1b[1;33m"
	BrightBlueText          Color = "\x1b[1;34m"
	BrightMagentaText       Color = "\x1b[1;35m"
	BrightCyanText          Color = "\x1b[1;36m"
	BrightWhiteText         Color = "\x1b[1;37m"
	BlackBackground         Color = "\x1b[0040m"
	RedBackground           Color = "\x1b[0041m"
	GreenBackground         Color = "\x1b[0042m"
	YellowBackground        Color = "\x1b[0043m"
	BlueBackground          Color = "\x1b[0044m"
	MagentaBackground       Color = "\x1b[0045m"
	CyanBackground          Color = "\x1b[0046m"
	WhiteBackground         Color = "\x1b[0047m"
	BrightBlackBackground   Color = "\x1b[0100m"
	BrightRedBackground     Color = "\x1b[0101m"
	BrightGreenBackground   Color = "\x1b[0102m"
	BrightYellowBackground  Color = "\x1b[0103m"
	BrightBlueBackground    Color = "\x1b[0104m"
	BrightMagentaBackground Color = "\x1b[0105m"
	BrightCyanBackground    Color = "\x1b[0106m"
	BrightWhiteBackground   Color = "\x1b[0107m"
)
