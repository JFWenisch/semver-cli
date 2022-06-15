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
func GetLatestTagFromBranch(branch string, releaseBranch bool) string {
	if getAmountOfTags() < 1 {
		return "0.0.0"
	}
	if releaseBranch {
		return GetLatestTag()
	}
	gitCmd := exec.Command("git", "tag", "--sort='creatordate'", "--merged "+branch)
	stdout, err := gitCmd.Output()

	if err != nil {
		var latestTag string = GetLatestTag()
		//latestTag = latestTag[0:strings.Index(latestTag, "-")]
		//latestTag = latestTag + "-" + branch + ".0"

		return latestTag
	}
	result := strings.TrimRight(string(stdout), "\r\n")
	return (result)

}
func getCommitsSinceTag(tag string) []string {

	gitCmd := exec.Command("git", "log", tag+"..HEAD", "--oneline")
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
