package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetTags() string {
	gitCmd := exec.Command("git", "tag")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if len(stdout) < 1 {
		result := "0.0.0"
		return (result)

	} else {
		result := strings.TrimRight(string(stdout), "\r\n")
		return (result)

	}

}
func GetLatestTag() string {
	gitCmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if len(stdout) < 1 {
		result := "0.0.0"
		return (result)

	} else {
		result := strings.TrimRight(string(stdout), "\r\n")
		return (result)

	}

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
	result := string(stdout)
	return (result)

}
