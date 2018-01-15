package cmd

import (
	"github.com/spf13/cobra"
	"io"

	"github.com/jenkins-x/jx/pkg/jx/cmd/templates"
	cmdutil "github.com/jenkins-x/jx/pkg/jx/cmd/util"
	"github.com/jenkins-x/jx/pkg/kube"
	"strings"
)

const (
	optionLabelColor = "label-color"
	optionTeamColor  = "team-color"
	optionEnvColor   = "env-color"
)

// PromptOptions containers the CLI options
type PromptOptions struct {
	CommonOptions

	NoLabel  bool
	ShowIcon bool

	Prefix    string
	Label     string
	Separator string
	Divider   string
	Suffix    string

	LabelColor []string
	TeamColor  []string
	EnvColor   []string
}

var (
	get_prompt_long = templates.LongDesc(`
		Generate a command prompt for the current team and environment.
`)

	get_prompt_example = templates.Examples(`
		# Generate the current prompt
		jx prompt

		# Enable the prompt for bash
		PS1="[\u@\h \W \$(jx prompt)]\$ "

		# Enable the prompt for zsh
		PROMPT='$(jx prompt)'$PROMPT
	`)
)

// NewCmdPrompt creates the new command for: jx get prompt
func NewCmdPrompt(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &PromptOptions{
		CommonOptions: CommonOptions{
			Factory: f,
			Out:     out,
			Err:     errOut,
		},
	}
	cmd := &cobra.Command{
		Use:     "prompt",
		Short:   "Generate the command line prompt for the current team and environment",
		Long:    get_prompt_long,
		Example: get_prompt_example,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			cmdutil.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&options.Prefix, "prefix", "p", "(", "The prefix text for the prompt")
	cmd.Flags().StringVarP(&options.Label, "label", "l", "k8s", "The label for the prompt")
	cmd.Flags().StringVarP(&options.Separator, "separator", "s", "|", "The separator between the label and the rest of the prompt")
	cmd.Flags().StringVarP(&options.Divider, "divider", "d", ":", "The divider between the team and environment for the prompt")
	cmd.Flags().StringVarP(&options.Suffix, "suffix", "x", ")", "The suffix text for the prompt")

	cmd.Flags().StringArrayVarP(&options.LabelColor, optionLabelColor, "", []string{"blue"}, "The color for the label")
	cmd.Flags().StringArrayVarP(&options.TeamColor, optionTeamColor, "", []string{"cyan"}, "The color for the team")
	cmd.Flags().StringArrayVarP(&options.TeamColor, optionEnvColor, "", []string{"cyan"}, "The color for the environment")

	cmd.Flags().BoolVarP(&options.NoLabel, "no-label", "", false, "Disables the use of the label in the prompt")
	cmd.Flags().BoolVarP(&options.ShowIcon, "icon", "i", false, "Uses an icon for the label in the prompt")

	return cmd
}

// Run implements this command
func (o *PromptOptions) Run() error {
	kubeClient, currentNs, err := o.Factory.CreateClient()
	if err != nil {
		return err
	}
	team, env, err := kube.GetDevNamespace(kubeClient, currentNs)
	if err != nil {
		return err
	}
	label := o.Label
	separator := o.Separator
	divider := o.Divider
	prefix := o.Prefix
	suffix := o.Suffix

	labelColor, err := cmdutil.GetColor(optionLabelColor, o.LabelColor)
	if err != nil {
		return err
	}
	teamColor, err := cmdutil.GetColor(optionLabelColor, o.TeamColor)
	if err != nil {
		return err
	}
	envColor, err := cmdutil.GetColor(optionLabelColor, o.EnvColor)
	if err != nil {
		return err
	}
	if o.NoLabel {
		label = ""
		separator = ""
	} else {
		if o.ShowIcon {
			label = "☸️  "
			label = labelColor.Sprint(label)
		} else {
			label = labelColor.Sprint(label)
		}
	}
	team = teamColor.Sprint(team)
	if env == "" {
		divider = ""
	} else {
		env = envColor.Sprint(env)
	}
	o.Printf("%s\n", strings.Join([]string{prefix, label, separator, team, divider, env, suffix}, ""))
	return nil
}