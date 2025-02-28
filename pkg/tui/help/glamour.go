// From https://github.com/analog-substance/util/blob/main/cli/glamour_help/main.go
// MIT Lincesed

package help

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var renderer *glamour.TermRenderer

func init() {
	var err error
	width, _, _ := term.GetSize(0)
	renderer, err = glamour.NewTermRenderer(
		// Detect the background color and pick either the default dark or light theme
		glamour.WithStandardStyle("auto"),
		glamour.WithEmoji(),
		glamour.WithPreservedNewLines(),
		glamour.WithStylesFromJSONBytes([]byte(`{ "document": { "margin": 0 } }`)),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		panic(err)
	}
}

func RenderMarkdown(markdown string) string {
	out, err := renderer.Render(markdown)
	if err != nil {
		panic(err)
	}

	return out
}

func GlamourUsage(c *cobra.Command) error {
	var b bytes.Buffer
	err := tmpl(&b, c.UsageTemplate(), c)
	if err != nil {
		c.PrintErrln(err)
		return err
	}
	pretty(c.ErrOrStderr(), b.String())
	return nil
}

func GlamourHelp(c *cobra.Command, _ []string) {
	var b bytes.Buffer
	err := tmpl(&b, c.HelpTemplate(), c)
	if err != nil {
		c.PrintErrln(err)
	}

	pretty(c.ErrOrStderr(), RenderMarkdown(b.String()))
}

func pretty(w io.Writer, s string) {
	fmt.Fprintf(w, "%s", s)
}

var templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"appendIfNotPresent":      appendIfNotPresent,
	"rpad":                    rpad,
	"fixIndentForPFlags":      fixIndentForPFlags,
	"gt":                      Gt,
	"eq":                      Eq,
}

func Gt(a, b any) bool {
	var left, right int64
	av := reflect.ValueOf(a)

	const intBits = 64

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = int64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = av.Int()
	case reflect.String:
		left, _ = strconv.ParseInt(av.String(), 10, intBits)
	}

	bv := reflect.ValueOf(b)

	switch bv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = int64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = bv.Int()
	case reflect.String:
		right, _ = strconv.ParseInt(bv.String(), 10, intBits)
	}

	return left > right
}

// FIXME Eq is unused by cobra and should be removed in a version 2. It exists
// only for compatibility with users of cobra.

// Eq takes two types and checks whether they are equal. Supported types are
// int and string. Unsupported types will panic.
func Eq(a, b any) bool {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		panic("Eq called on unsupported type")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return av.Int() == bv.Int()
	case reflect.String:
		return av.String() == bv.String()
	}
	return false
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// appendIfNotPresent will append stringToAppend to the end of s, but only if it's not yet present in s.
func appendIfNotPresent(s, stringToAppend string) string {
	if strings.Contains(s, stringToAppend) {
		return s
	}
	return s + " " + stringToAppend
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	formattedString := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(formattedString, s)
}

// rpad adds padding to the right of a string.
func fixIndentForPFlags(subject string) string {
	indentNeeded := "  "
	return fmt.Sprintf("%s%s", indentNeeded, strings.ReplaceAll(subject, "\n", "\n"+indentNeeded))
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data any) error {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

func AddToRootCmd(rootCmd *cobra.Command) {
	rootCmd.SetUsageFunc(GlamourUsage)
	rootCmd.SetHelpFunc(GlamourHelp)
	rootCmd.SetHelpTemplate(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)

	rootCmd.SetUsageTemplate(`## Usage{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

## Aliases
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

## Examples
` + "```bash\n{{.Example}}\n```" + `{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

## Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

## Additional Commands{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

## Flags
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces | fixIndentForPFlags}}{{end}}{{if .HasAvailableInheritedFlags}}

## Global Flags
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces | fixIndentForPFlags}}{{end}}{{if .HasHelpSubCommands}}

## Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}
