package cmd

import (
	"github.com/spf13/cobra"

	"github.com/redbubble/yak/aws"
	"github.com/redbubble/yak/cli"
)

func shimCmd(cmd *cobra.Command, args []string) error {
	roleName, err := cli.ResolveRole(args[0])

	if err != nil {
		return err
	}

	command := args[1:]

	creds := cli.AssumeRoleFromCache(roleName)

	if creds == nil {
		loginData, err := cli.GetLoginDataWithTimeout()

		if err != nil {
			return err
		}

		cli.CacheLoginRoles(loginData.Roles)
		creds, err = cli.AssumeRole(loginData, roleName)

		if err != nil {
			return err
		}

		cli.CacheCredentials(roleName, creds)
		cli.WriteCache()
	}

	return cli.Exec(
		command,
		cli.EnrichedEnvironment(
			aws.EnvironmentVariables(creds),
		),
	)
}
