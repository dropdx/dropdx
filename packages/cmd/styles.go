package cmd

import (
	"github.com/pterm/pterm"
)

var (
	successStyle   = pterm.NewStyle(pterm.FgGreen, pterm.Bold)
	infoStyle      = pterm.NewStyle(pterm.FgCyan)
	warnStyle      = pterm.NewStyle(pterm.FgYellow)
	errorStyle     = pterm.NewStyle(pterm.FgRed, pterm.Bold)
	boldStyle      = pterm.NewStyle(pterm.Bold)
	mutedStyle     = pterm.NewStyle(pterm.FgGray)
	accentStyle    = pterm.NewStyle(pterm.FgYellow, pterm.Bold)
	headerStyle    = pterm.NewStyle(pterm.FgWhite, pterm.Bold, pterm.Underscore)
	tokenNameStyle = pterm.NewStyle(pterm.FgMagenta, pterm.Bold)
)

func success(s string) string    { return successStyle.Sprint(s) }
func info(s string) string       { return infoStyle.Sprint(s) }
func warn(s string) string       { return warnStyle.Sprint(s) }
func errCrit(s string) string    { return errorStyle.Sprint(s) }
func bold(s string) string       { return boldStyle.Sprint(s) }
func muted(s string) string      { return mutedStyle.Sprint(s) }
func accent(s string) string     { return accentStyle.Sprint(s) }
func header(s string) string     { return headerStyle.Sprint(s) }
func tokenStyle(s string) string { return tokenNameStyle.Sprint(s) }
