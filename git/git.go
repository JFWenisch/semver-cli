package git

import (
	"fmt"
	"os"
	"os/exec"
)

func GetTags() string {
	gitCmd := exec.Command("git", "tag")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())

	}

	if len(stdout) < 1 {
		result := "0.0.0"
		return (result)

	} else {
		result := string(stdout)
		return (result)

	}

}
func CreateTag(version string) {
	gitTagCreateCmd := exec.Command("git", "tag  "+version)
	stdout, err := gitTagCreateCmd.Output()

	if err != nil {
		fmt.Println("error creating tag " + version + " " + err.Error())
		os.Exit(-1)

	}
	fmt.Println(string(stdout))

	gitTagPushCmd := exec.Command("git", "push orign  "+version)
	gitTagPushCmdOut, gitTagPushCmdErr := gitTagPushCmd.Output()

	if gitTagPushCmdErr != nil {
		fmt.Println("error pushing tag " + version + " " + gitTagPushCmdErr.Error())
		os.Exit(-1)
	}
	fmt.Println(string(gitTagPushCmdOut))

}
func GetCurrentBranch() string {
	gitCmd := exec.Command("git", "rev-parse --abbrev-ref HEAD")
	stdout, err := gitCmd.Output()

	if err != nil {
		fmt.Println(err.Error())

	}

	result := string(stdout)
	return (result)

}
