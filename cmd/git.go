package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"strconv"
	"strings"
)

func GetTags() string {
	if IsVerbose() {
		fmt.Println("Running 'git tag'")
	}
	gitCmd := exec.Command("git", "tag")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println("Exception running 'git tag'")
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	result := strings.TrimRight(string(stdout), "\r\n")
	return (result)

}
func GetLatestTag() string {
	if getAmountOfTags() < 1 {
		if IsVerbose() {
			fmt.Println("No tags found. Using '0.0.0'")
		}
		return "0.0.0"
	}
	gitCmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	result := strings.TrimRight(string(stdout), "\r\n")
	return (result)

}
func GetLatestTagFromBranch(branch string, releaseBranch bool) string {
	if getAmountOfTags() < 1 {
		if IsVerbose() {
			fmt.Println("Could not find any tags")
		}
		return "0.0.0"
	}
	if releaseBranch {
		if IsVerbose() {
			fmt.Println("The supplied branch '" + branch + "' is set as release branch'")
		}
		return GetLatestTag()
	}
	if IsVerbose() {
		fmt.Println("Searching latest tag on '" + branch + "' using 'git describe --match *-%BRANCH%.* --tags --abrev=0'")
	}
	args := []string{"describe", "--match", "*-develop.*", "--abbrev=0", "--tags"}
	gitCmd := exec.Command("git", args...)
	stdout, err := gitCmd.CombinedOutput()

	if err != nil {
		fmt.Println("Exception getting tags on current branch")
		fmt.Println(err)
		var latestTag string = GetLatestTag()
		//latestTag = latestTag[0:strings.Index(latestTag, "-")]
		//latestTag = latestTag + "-" + branch + ".0"

		return latestTag
	}
	result := strings.TrimRight(string(stdout), "\r\n")
	return (result)

}
func getCommitsSinceTag(tag string) []string {
	gitCmd := exec.Command("git", "log", "--oneline")
	if tag != "0.0.0" {
		gitCmd = exec.Command("git", "log", tag+"..HEAD", "--oneline")

	}
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	result := strings.Split(strings.ReplaceAll(string(stdout), "\r\n", "\n"), "\n")

	return (result)

}
func DetectBumpTypeFromTag(tag string) string {
	commits := getCommitsSinceTag(tag)

	var isMajor, isMinor, isPatch bool
	for _, commit := range commits {
		if strings.Contains(strings.ToLower(commit), "breaking change: ") {
			isMajor = true
		}
		if strings.Contains(strings.ToLower(commit), "feat: ") {
			isMinor = true
		}
		if strings.Contains(strings.ToLower(commit), "fix: ") {
			isPatch = true
		}

	}
	if isMajor {
		return "mayor"
	}
	if isMinor {
		return "minor"
	}
	if isPatch {
		return "patch"
	}

	//fmt.Println("No angular style commit detected in last " + string(len(commits)) + " commits. Aborting")
	os.Exit(-1)
	return "unknown"
}
func getAmountOfTags() int {
	if IsVerbose() {
		fmt.Println("Checking if there are already tags by running 'git rev-list --tags --count'")
	}
	gitCmd := exec.Command("git", "rev-list", "--tags", "--count")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprintln(os.Stderr, "Could not list amount of tags")
		os.Exit(-1)
	}

	result := strings.TrimRight(string(stdout), "\r\n")
	if IsVerbose() {
		fmt.Println("Found " + result + " existing Tags")
	}
	tagAmout, err := strconv.Atoi(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot parse amount of tags")
		os.Exit(-1)
	}

	return (tagAmout)
}
func CreateTag(version string) {
	if IsVerbose() {
		fmt.Println("Trying to create a new tag running 'git tag " + version + "'")
	}
	gitTagCreateCmd := exec.Command("git", "tag", version)
	stdout, err := gitTagCreateCmd.Output()

	if err != nil {
		fmt.Println("error creating tag " + version + " " + err.Error())
		os.Exit(-1)

	}
	if IsVerbose() {
		fmt.Println("Tag successfully created")
		fmt.Println(string(stdout))
		fmt.Println("Trying to push tag running 'git push --tags'")
	}

	gitTagPushCmd := exec.Command("git", "push", "--tags")
	gitTagPushCmdOut, gitTagPushCmdErr := gitTagPushCmd.CombinedOutput()
	fmt.Println(string(gitTagPushCmdOut))
	if gitTagPushCmdErr != nil {

		fmt.Println("error pushing tag " + version + " " + gitTagPushCmdErr.Error())
		os.Exit(-1)
	}
	if IsVerbose() {
		fmt.Println("Successfully created tag'")
		fmt.Println(string(gitTagPushCmdOut))
	}

}
func GetCurrentBranch() string {
	if IsVerbose() {
		fmt.Println("Trying to get current branch running 'git rev-parse --abbrev-ref HEAD'")
	}
	gitCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println("Exception running 'git branch --show-current'")
		fmt.Println(err.Error())

	}
	result := strings.TrimRight(string(stdout), "\r\n")
	if IsVerbose() {
		fmt.Println("The current branch is '" + result + "'")
	}
	return (result)

}
