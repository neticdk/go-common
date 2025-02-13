package ui

import (
	"io"

	"github.com/pterm/pterm"
)

var (
	Info        = pterm.Info
	Warning     = pterm.Warning
	Success     = pterm.Success
	Error       = pterm.Error
	Fatal       = pterm.Fatal
	Debug       = pterm.Debug
	Description = pterm.Description

	// InfoI returns a PrefixPrinter, which can be used to print text with an "info" Prefix.
	InfoI = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.InfoMessageStyle,
			Text:  "🦶",
		},
	}

	// WarningI returns a PrefixPrinter, which can be used to print text with a "warning" Prefix.
	WarningI = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.WarningMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.WarningMessageStyle,
			Text:  "⚠️",
		},
	}

	// SuccessI returns a PrefixPrinter, which can be used to print text with a "success" Prefix.
	SuccessI = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.SuccessMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.SuccessMessageStyle,
			Text:  "✓",
		},
	}

	// ErrorI returns a PrefixPrinter, which can be used to print text with an "error" Prefix.
	ErrorI = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.ErrorMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.ErrorMessageStyle,
			Text:  "✗",
		},
	}

	// Fatal returns a PrefixPrinter, which can be used to print text with an "fatal" Prefix.
	// NOTICE: Fatal terminates the application immediately!
	FatalI = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.FatalMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.FatalMessageStyle,
			Text:  "💀",
		},
		Fatal: true,
	}

	// Debug Prints debug messages. By default, it will only print if PrintDebugMessages is true.
	// You can change PrintDebugMessages with EnableDebugMessages and DisableDebugMessages, or by setting the variable itself.
	DebugI = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.DebugMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.DebugMessageStyle,
			Text:  "🪲",
		},
		Debugger: true,
	}

	// Check returns a PrefixPrinter, which can be used to print text with a "mark check" Prefix.
	CheckI = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.SuccessMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.SuccessMessageStyle,
			Text:  "✓",
		},
	}
)

func updateWriters(w io.Writer) {
	Info.Writer = w
	Warning.Writer = w
	Success.Writer = w
	Error.Writer = w
	Fatal.Writer = w
	Debug.Writer = w
	Description.Writer = w

	InfoI.Writer = w
	WarningI.Writer = w
	SuccessI.Writer = w
	ErrorI.Writer = w
	FatalI.Writer = w
	DebugI.Writer = w
	CheckI.Writer = w
}
