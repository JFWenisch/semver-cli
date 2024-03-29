/*
Copyright © 2022 Jean-Fabian Wenisch <hello@wenisch.tech>


*/
package cmd

import (
	"fmt"
	"os"
	"regexp"

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
	tagsBumpCmd.Flags().StringVarP(&releaseBranch, "release-branches", "r", "main", "Comma seperated list of release branches. When command is executed on a non-release branch, a pre-release version is created'")
	tagsBumpCmd.Flags().StringVarP(&BumpType, "type", "t", "", "Type of commit, e.g. 'major', 'minor', 'patch' or 'detect'")
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
			fmt.Println(GetLatestTag())
		} else {
			fmt.Println(GetTags())
		}

	},
}

var tagsBumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "creates the next tag",
	Long: `Funcitionality to identify and/or create the next tag in relation to semver conventions. For example:

	semver-cli tags bump -t patch  -p v 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if BumpType != "major" && BumpType != "minor" && BumpType != "patch" && BumpType != "detect" {
			fmt.Fprintln(os.Stderr, "specified type is not 'major','minor', 'patch' or 'detect'")
			os.Exit(-1)
		}

		var branch string = GetCurrentBranch()
		var tag string = GetLatestTagFromBranch(branch, branch == releaseBranch)
		if IsVerbose() {
			fmt.Println("The latest found tag is '" + tag + "'")
		}
		var nextTag = generateNextTag(branch, tag, BumpType)

		if DryRun {
			fmt.Println(nextTag)
		} else {
			CreateTag(nextTag)
		}
	},
}

func generateNextTag(branch string, tag string, bumptype string) string {

	if bumptype == "detect" {
		bumptype = DetectBumpTypeFromTag(tag)
	}

	//First 3 numbers found are used as version
	re := regexp.MustCompile("[0-9]+")
	splittedTag := re.FindAllString(tag, -1)
	majorVerion, minorVersion, patchVersion := splittedTag[0], splittedTag[1], splittedTag[2]

	if branch == releaseBranch {

		if bumptype == "major" {
			currentVersion, err := strconv.Atoi(majorVerion)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot parse "+bumptype+" of tag "+tag)
				os.Exit(-1)
			}
			currentVersion++
			majorVerion = strconv.Itoa(currentVersion)
		}
		if bumptype == "minor" {
			currentVersion, err := strconv.Atoi(minorVersion)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot parse "+bumptype+" of tag "+tag)
				os.Exit(-1)
			}
			currentVersion++
			minorVersion = strconv.Itoa(currentVersion)
		}
		if bumptype == "patch" {
			currentVersion, err := strconv.Atoi(patchVersion)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot parse "+bumptype+" of tag "+tag)
				os.Exit(-1)
			}
			currentVersion++
			patchVersion = strconv.Itoa(currentVersion)
		}
		return tagPrefix + majorVerion + "." + minorVersion + "." + patchVersion
	} else {

		if len(splittedTag) > 2 {

			branchVersion, err := strconv.Atoi(splittedTag[3])
			strconv.Itoa(branchVersion)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot parse "+bumptype+" of tag "+tag)
				os.Exit(-1)
			}
			branchVersion++
			return tagPrefix + majorVerion + "." + minorVersion + "." + patchVersion + "-" + branch + "." + strconv.Itoa(branchVersion)
		} else {
			return tagPrefix + majorVerion + "." + minorVersion + "." + patchVersion + "-" + branch + "." + "1"

		}
	}
}
