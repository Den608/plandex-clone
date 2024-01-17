package cmd

import (
	"fmt"
	"os"
	"plandex/api"
	"plandex/auth"
	"plandex/format"
	"plandex/lib"
	"plandex/term"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"ls"},
	Short:   "List everything in context",
	Run:     context,
}

func context(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	contexts, err := api.Client.ListContext(lib.CurrentPlanId)

	if err != nil {
		color.New(color.FgRed).Fprintln(os.Stderr, "Error listing context:", err)
		os.Exit(1)
	}

	totalTokens := 0
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Name", "Type", "🪙", "Added", "Updated"})
	table.SetAutoWrapText(false)

	if len(contexts) == 0 {
		fmt.Println("🤷‍♂️ No context")
		fmt.Println()
		term.PrintCmds("", "load")
		return
	}

	for i, context := range contexts {
		totalTokens += context.NumTokens

		t, icon := lib.GetContextTypeAndIcon(context)

		row := []string{
			strconv.Itoa(i + 1),
			" " + icon + " " + context.Name,
			t,
			strconv.Itoa(context.NumTokens), //+ " 🪙",
			format.Time(context.CreatedAt),
			format.Time(context.UpdatedAt),
		}
		table.Rich(row, []tablewriter.Colors{
			{tablewriter.FgHiWhiteColor, tablewriter.Bold},
			{tablewriter.FgHiGreenColor, tablewriter.Bold},
			{tablewriter.FgHiWhiteColor},
			{tablewriter.FgHiWhiteColor},
			{tablewriter.FgHiWhiteColor},
			{tablewriter.FgHiWhiteColor},
		})
	}

	table.Render()

	tokensTbl := tablewriter.NewWriter(os.Stdout)
	tokensTbl.SetAutoWrapText(false)
	tokensTbl.Append([]string{color.New(color.FgHiCyan, color.Bold).Sprintf("Total tokens →") + color.New(color.FgHiWhite, color.Bold).Sprintf(" %d 🪙", totalTokens)})

	tokensTbl.Render()

	fmt.Println()
	term.PrintCmds("", "load", "rm", "clear")

}

func init() {
	RootCmd.AddCommand(contextCmd)

}
