package github

import (
	"os"
	"os/exec"
)

const repoDir = "data/repositories/"

func CloneRepo(name, url string) error {
	dir := repoDir
	if _, err := os.Stat(dir); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("git", "clone", url, name)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func PullRepo(name string) error {
	dir := repoDir + name
	cmd := exec.Command("git", "pull")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ReadFile(name, path string) ([]byte, error) {
	dir := repoDir + name + "/" + path

	return os.ReadFile(dir)
}

func ListFiles(name string, subdir string) ([]string, error) {
	dir := repoDir + name + "/" + subdir

	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdirnames(-1)
}
