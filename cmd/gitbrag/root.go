package gitbrag

import (
	"os"
	"time"

	"github.com/radulucut/gitbrag/internal"
	"github.com/radulucut/gitbrag/internal/utils"
	"github.com/spf13/cobra"
)

var Version string

func Execute() {
	time := utils.NewTime()
	printer := internal.NewPrinter(os.Stdin, os.Stdout, os.Stderr)
	core := internal.NewCore(time, printer)
	root, err := NewRoot(Version, time, printer, core)
	if err != nil {
		printer.ErrPrintf("Error: %v\n", err)
		os.Exit(1)
	}

	err = root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type Root struct {
	Cmd *cobra.Command

	version string

	time    utils.Time
	printer *internal.Printer
	core    *internal.Core
}

func NewRoot(
	version string,
	time utils.Time,
	printer *internal.Printer,
	core *internal.Core,
) (*Root, error) {
	root := &Root{
		version: version,
		time:    time,
		printer: printer,
		core:    core,
	}

	root.Cmd = &cobra.Command{
		Use:   "gitbrag [directories...]",
		Short: "Display git statistics for local repositories",
		Long: `A terminal tool that outputs git stats for local git repositories.
It shows files changed, insertions and deletions for specified directories.

Examples:
  gitbrag ./
  gitbrag ./ projects

  # Filter by date range
  gitbrag ./ --since 2024-01-01
  gitbrag ./ --since 2024-01-01 --until 2024-12-31
  gitbrag ./ --since 7d

  # Filter by author name or email
  gitbrag ./ --author "John Doe"
  gitbrag ./ --since 7d --author john@example.com

  # Output statistics to PNG file
  gitbrag ./ -O stats.png
  gitbrag ./ --output stats.png --background "#282a36"

  # Use custom background and foreground colors
  gitbrag ./ -O stats.png -B fff
  gitbrag ./ -O stats.png --color "#50fa7b"
  gitbrag ./ -O stats.png -B "#282a36" -C "f8f8f2"
  gitbrag ./ -O stats.png -B 000 --color fff

  # Show language breakdown (top 3 + Others)
  gitbrag ./ -O stats.png --lang
`,
		Version: version,
		RunE:    root.RunRoot,
		Args:    cobra.MinimumNArgs(1),
	}

	root.Cmd.SetOut(root.printer.OutWriter)
	root.Cmd.SetErr(root.printer.ErrWriter)

	flags := root.Cmd.Flags()
	flags.String("since", "", "specific date (e.g. 2024-01-01 12:03:04) or duration (e.g. 1d)")
	flags.String("until", "", "specific date (e.g. 2024-12-31 23:59:59)")
	flags.String("author", "", "filter by author name or email")
	flags.StringP("output", "O", "", "export statistics to PNG file (e.g. stats.png)")
	flags.StringP("background", "B", "", "background color in hex format (e.g. #282a36 or 282a36), transparent by default")
	flags.StringP("color", "C", "", "text color in hex format (e.g. #f8f8f2 or f8f8f2)")
	flags.Bool("lang", false, "show language breakdown with top 3 languages and others (PNG output only)")

	root.initVersion()

	return root, nil
}

func (r *Root) RunRoot(cmd *cobra.Command, args []string) error {
	since, err := r.parseSinceFlag(cmd.Flag("since").Value.String())
	if err != nil {
		return err
	}
	until, err := r.parseUntilFlag(cmd.Flag("until").Value.String())
	if err != nil {
		return err
	}
	author := cmd.Flag("author").Value.String()
	output := cmd.Flag("output").Value.String()
	background := cmd.Flag("background").Value.String()
	color := cmd.Flag("color").Value.String()
	lang, _ := cmd.Flags().GetBool("lang")
	return r.core.Run(&internal.RunOptions{
		Dirs:       args,
		Since:      since,
		Until:      until,
		Author:     author,
		Output:     output,
		Background: background,
		Color:      color,
		Lang:       lang,
	})
}

func (r *Root) parseSinceFlag(flag string) (time.Time, error) {
	if flag == "" {
		return time.Time{}, nil
	}
	d, err := utils.ParseDuration(flag)
	if err == nil {
		return r.time.Now().Add(-d), nil
	}
	return utils.ParseDateTime(flag)
}

func (r *Root) parseUntilFlag(flag string) (time.Time, error) {
	if flag == "" {
		return time.Time{}, nil
	}
	return utils.ParseDateTime(flag)
}
