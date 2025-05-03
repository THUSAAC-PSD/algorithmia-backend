package environment

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"

	"emperror.dev/errors"
	"github.com/spf13/viper"
)

func FixProjectRootWorkingDirectoryPath() {
	currentWD, _ := os.Getwd()
	log.Printf("Current working directory is: `%s`", currentWD)

	rootDir := GetProjectRootWorkingDirectory()
	_ = os.Chdir(rootDir)
	newWD, _ := os.Getwd()

	log.Printf("New fixed working directory is: `%s`", newWD)
}

func GetProjectRootWorkingDirectory() string {
	var rootWorkingDirectory string
	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
	// viper will get it from `os env` or a .env file
	pn := viper.GetString(constant.ProjectNameEnv)
	if pn != "" {
		rootWorkingDirectory = getProjectRootDirectoryFromProjectName(pn)
	} else {
		wd, _ := os.Getwd()
		dir, err := searchRootDirectory(wd)
		if err != nil {
			log.Fatal(err)
		}
		rootWorkingDirectory = dir
	}

	absoluteRootWorkingDirectory, _ := filepath.Abs(rootWorkingDirectory)

	return absoluteRootWorkingDirectory
}

func getProjectRootDirectoryFromProjectName(pn string) string {
	// set root working directory of our app in the viper
	// https://stackoverflow.com/a/47785436/581476
	wd, _ := os.Getwd()

	for !strings.HasSuffix(wd, pn) {
		wd = filepath.Dir(wd)
	}

	return wd
}

func searchRootDirectory(
	dir string,
) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", errors.WrapIf(err, "failed to read directory")
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if strings.EqualFold(
				fileName,
				"go.mod",
			) {
				return dir, nil
			}
		}
	}

	parentDir := filepath.Dir(dir)
	if parentDir == dir {
		return "", errors.WrapIf(err, "could not find go.mod file")
	}

	return searchRootDirectory(parentDir)
}
