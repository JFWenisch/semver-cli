/*
Copyright © 2022 Jean-Fabian Wenisch <hello@wenisch.tech>


*/
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"semver-cli/git"
	"strconv"

	"github.com/spf13/cobra"
)

var Latest bool
var DryRun bool
var releaseBranch string
var BumpType string
var tagPrefix string

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: must also specify a subcommand")

		}

	},
}

func init() {
	tagsListCmd.PersistentFlags().BoolVarP(&Latest, "latest", "l", false, "Only display latest tag")
	tagsCmd.AddCommand(tagsListCmd)

	tagsBumpCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", false, "Outputs the next determined version without creating it")
	tagsBumpCmd.Flags().StringVarP(&releaseBranch, "release-branches", "r", "main,master", "Comma seperated list of release branches. When command is executed on a non-release branch, a pre-release version is created'")
	tagsBumpCmd.Flags().StringVarP(&BumpType, "type", "t", "", "Type of commit, e.g. 'major', 'minor' or 'patch'")
	tagsBumpCmd.Flags().StringVarP(&tagPrefix, "prefix", "p", "", "The prefix for tagging e.g. 'v'")
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
		if Latest {
			fmt.Println(git.GetLatestTag())
		} else {
			fmt.Println(git.GetTags())
		}

	},
}

var tagsBumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "creates the next tag",
	Long:  `Funcitionality to identify and/or create the next tag in relation to semver conventions`,
	Run: func(cmd *cobra.Command, args []string) {
		if BumpType != "major" && BumpType != "minor" && BumpType != "patch" {
			fmt.Fprintln(os.Stderr, "specified type is not 'major','minor' or 'patch'")
			os.Exit(-1)
		}
		var branch string = git.GetCurrentBranch()
		var tag string = git.GetLatestTag()

		var nextTag = generateNextTag(branch, tag, BumpType)

		if DryRun {
			fmt.Println(nextTag)
		} else {
			git.CreateTag(nextTag)
		}
	},
}

func generateNextTag(branch string, tag string, bumptype string) string {
	//First 3 numbers found are used as version
	re := regexp.MustCompile("[0-9]+")
	splittedTag := re.FindAllString(tag, -1)
	majorVerion, minorVersion, patchVersion := splittedTag[0], splittedTag[1], splittedTag[2]
	var finalVersion string
	if branch == releaseBranch {

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
		finalVersion = tagPrefix + majorVerion + "." + minorVersion + "." + patchVersion
	} else {
		branchVersion, err := strconv.Atoi(splittedTag[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot parse "+BumpType+" of tag "+tag)
			os.Exit(-1)
		}
		branchVersion++
		finalVersion = tagPrefix + majorVerion + "." + minorVersion + "." + patchVersion + "-" + branch + "." + strconv.Itoa(branchVersion)
	}

	return finalVersion
}
