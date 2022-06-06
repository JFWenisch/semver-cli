package git

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetTags() string {
	gitCmd := exec.Command("git", "tag")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	result := strings.TrimRight(string(stdout), "\r\n")
	return (result)

}
func GetLatestTag() string {
	if getAmountOfTags() < 1 {
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
func getAmountOfTags() int {
	gitCmd := exec.Command("git", "rev-list", "--tags", "--count")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	result := strings.TrimRight(string(stdout), "\r\n")
	tagAmout, err := strconv.Atoi(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot parse amount of tags")
		os.Exit(-1)
	}

	return (tagAmout)
}
func CreateTag(version string) {
	gitTagCreateCmd := exec.Command("git", "tag", version)
	stdout, err := gitTagCreateCmd.Output()

	if err != nil {
		fmt.Println("error creating tag " + version + " " + err.Error())
		os.Exit(-1)

	}
	fmt.Println(string(stdout))

	gitTagPushCmd := exec.Command("git", "push", "orign", "--tags")
	gitTagPushCmdOut, gitTagPushCmdErr := gitTagPushCmd.CombinedOutput()

	if gitTagPushCmdErr != nil {
		fmt.Println("error pushing tag " + version + " " + gitTagPushCmdErr.Error())
		os.Exit(-1)
	}
	fmt.Println(string(gitTagPushCmdOut))

}
func GetCurrentBranch() string {
	gitCmd := exec.Command("git", "branch", "--show-current")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())

	}
	result := strings.TrimRight(string(stdout), "\r\n")
	return (result)

}
