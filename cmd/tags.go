/*
Copyright © 2022 Jean-Fabian Wenisch <hello@wenisch.tech>


*/
package cmd

import (
	"fmt"
	"os"
	"semver-cli/git"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var Latest bool
var DryRun bool
var BumpType string

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Error: must also specify a subcommand")
		}

	},
}

func init() {
	tagsListCmd.PersistentFlags().BoolVarP(&Latest, "latest", "l", false, "Only display latest tag")
	tagsCmd.AddCommand(tagsListCmd)

	tagsBumpCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", false, "Outputs the next determined version without creating it")
	tagsBumpCmd.Flags().StringVarP(&BumpType, "type", "t", "", "Type of commit, e.g. 'major', 'minor' or 'patch'")
	tagsBumpCmd.MarkFlagRequired("type")
	tagsCmd.AddCommand(tagsBumpCmd)
	rootCmd.AddCommand(tagsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists tags",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(git.GetTags())

	},
}

var tagsBumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "creates the next tag",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if BumpType != "major" && BumpType != "minor" && BumpType != "patch" {
			fmt.Fprintln(os.Stderr, "specified type is not 'major','minor' or 'patch'")
			os.Exit(-1)
		}
		tag := git.GetTags()
		splittedTag := strings.Split(tag, ".")
		majorVerion, minorVersion, patchVersion := splittedTag[0], splittedTag[1], splittedTag[2]

		if BumpType == "major" {
			currentVersion, err := strconv.Atoi(majorVerion)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot parse "+BumpType+" of tag "+tag)
				os.Exit(-1)
			}
			currentVersion++
			majorVerion = strconv.Itoa(currentVersion)
		}
		if BumpType == "minor" {
			currentVersion, err := strconv.Atoi(minorVersion)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot parse "+BumpType+" of tag "+tag)
				os.Exit(-1)
			}
			currentVersion++
			minorVersion = strconv.Itoa(currentVersion)
		}
		if BumpType == "patch" {
			currentVersion, err := strconv.Atoi(patchVersion)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot parse "+BumpType+" of tag "+tag)
				os.Exit(-1)
			}
			currentVersion++
			patchVersion = strconv.Itoa(currentVersion)
		}

		var finalVersion string = majorVerion + "." + minorVersion + "." + patchVersion
		if DryRun {
			fmt.Println(finalVersion)
		} else {
			git.CreateTag(finalVersion)
		}
	},
}
