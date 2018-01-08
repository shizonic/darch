package hooks

import (
	"fmt"

	"../../hooks"
	"../../stage"
	"github.com/ryanuber/columnize"
	"github.com/urfave/cli"
)

func listCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List hooks currently installed.",
		Action: func(c *cli.Context) error {
			err := list()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func detailsCommand() cli.Command {
	return cli.Command{
		Name:      "details",
		Usage:     "Get info about a hook, including what images the hook applies to (based on configuration).",
		ArgsUsage: "HOOK_NAME",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "include-matched-images",
			},
		},
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			err := details(c.Args().First(), c.Bool("include-matched-images"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func helpCommand() cli.Command {
	return cli.Command{
		Name:      "help",
		Usage:     "Print help for the given hook.",
		ArgsUsage: "HOOK_NAME",
		Action: func(c *cli.Context) error {
			if len(c.Args()) != 1 {
				return cli.NewExitError(fmt.Errorf("Unexpected arguements"), 1)
			}
			err := help(c.Args().First())
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:  "hooks",
		Usage: "Commands manage/view hooks.",
		Subcommands: []cli.Command{
			listCommand(),
			detailsCommand(),
			helpCommand(),
		},
	}
}

func list() error {

	hooks, err := hooks.GetHooks()

	if err != nil {
		return err
	}

	result := []string{
		"Name | Path",
	}

	for _, hook := range hooks {
		result = append(result, hook.Name+" | "+hook.Path)
	}

	fmt.Println(columnize.SimpleFormat(result))

	return nil
}

func details(hookName string, includeMatchedImages bool) error {
	hook, err := hooks.GetHook(hookName)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", hook.Name)
	fmt.Printf("Path: %s\n", hook.Path)
	fmt.Printf("ExecutionOrder: %d\n", hook.ExecutionOrder)
	fmt.Printf("IncludeImages:\n")
	for _, includeImage := range hook.IncludeImages {
		fmt.Printf("\t%s\n", includeImage)
	}
	fmt.Printf("ExcludeImages:\n")
	for _, excludeImage := range hook.ExcludeImages {
		fmt.Printf("\t%s\n", excludeImage)
	}

	if includeMatchedImages {
		fmt.Printf("MatchedImages:\n")
		images, err := stage.GetAllStaged("/var/darch/staged")
		if err != nil {
			return err
		}
		for _, image := range images {
			for _, imageTag := range image.Tags {
				if hooks.AppliesToStagedTag(hook, imageTag) {
					fmt.Printf("\t%s\n", imageTag.FullName)
				}
			}
		}
	}

	return nil
}

func help(hookName string) error {
	hook, err := hooks.GetHook(hookName)
	if err != nil {
		return err
	}

	return hooks.PrintHookHelp(hook)
}
