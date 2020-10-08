package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	// This file should be located in the same directory the
	// the `gitpacker` binary is located in
	// TODO: allow loading of file from environment or command line
	// filepath argument.
	DefaultGitPackerConfigFilename = "pack.json"
	ZipFileSuffix                  = ".zip"
)

// CloneConfig wraps parameters for cloning a git repository.
type CloneConfig struct {
	// The name of the directory to clone the repository to,
	// the directory will be created if it does not exist.
	CloneDirectory string `json:"clone_directory"`
	GitURL         string `json:"git_url"`
	Commit         string `json:"commit,omitempty"`
	Shallow        bool   `json:"shallow"`
}

// PackConfig wraps parameters for an execution of the
// GitPacker program.
type PackConfig struct {
	// The name of the top level directory to clone all the repos into
	RootCloneDirectory string        `json:"root_clone_directory"`
	Repos              []CloneConfig `json:"repos"`
	// Whether an single file archive should be generated containing
	// all repositories specified in Repos.
	Archive         bool   `json:"archive"`
	ArchiveFilename string `json:"archive_filename"`
}

var (
	packConfig PackConfig
)

func (c CloneConfig) Clone(rootDirectory string) error {
	// Clone the given repository to the given directory
	cloneDirectoryPath := c.CloneDirectory
	if rootDirectory != "" {
		cloneDirectoryPath = fmt.Sprintf("%s/%s", rootDirectory, c.CloneDirectory)
	}
	r, err := git.PlainClone(cloneDirectoryPath, false, &git.CloneOptions{
		URL: c.GitURL,
	})
	if err != nil {
		return err
	}
	// Only checkout to specific commit if specified
	// otherwise latest commit is fine
	if c.Commit == "" {
		// Delete git history if shallow is specified
		if c.Shallow {
			err = os.RemoveAll(fmt.Sprintf("%s/.git/", cloneDirectoryPath))
			if err != nil {
				return err
			}
		}
		return nil
	}
	// ... retrieving the commit being pointed by HEAD
	fmt.Println("git show-ref --head HEAD")
	ref, err := r.Head()
	if err != nil {
		return err
	}
	fmt.Println(ref.Hash())
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	// ... checking out to commit
	fmt.Printf("git checkout %s\n", c.Commit)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(c.Commit),
	})
	if err != nil {
		return err
	}
	// Delete git history if shallow is specified
	// Must be done after checking out to specific commit as otherwise
	// there would be no info for git to checkout to that commit locally
	if c.Shallow {
		err = os.RemoveAll(fmt.Sprintf("%s/.git/", cloneDirectoryPath))
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Load gitpacker config
	err := loadJSON(DefaultGitPackerConfigFilename, &packConfig)
	if err != nil {
		fmt.Printf("Error %s loading GitPacker config from %s", err, DefaultGitPackerConfigFilename)
		return
	}
	fmt.Printf("Loaded GitPack config %+v\n", packConfig)
	if strings.HasPrefix(packConfig.RootCloneDirectory, "/") {
		fmt.Printf("RootCloneDirectory %s can not be an absolute path / begin with /\n", packConfig.RootCloneDirectory)
		return
	}
	// Clone all specified git repositories
	var errs []error
	for _, repo := range packConfig.Repos {
		// Can't reliably be done with naive goroutine for each clone
		// as some remote providers (i.e. Github) limit number of open connections per IP
		err = repo.Clone(packConfig.RootCloneDirectory)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Printf("Cloning error %+v\n", err)
		}
		// Exit early if any errors cloning
		return
	}
	// If specified, compress all cloned repositories
	// into a single zip file
	if !packConfig.Archive {
		return
	}
	if packConfig.ArchiveFilename == "" {
		fmt.Println("Must specify archive_filename if archive is true.")
		return
	}
	// Ensure filename ends with .zip
	archiveFilename := packConfig.ArchiveFilename
	if !(strings.HasSuffix(archiveFilename, ZipFileSuffix)) {
		archiveFilename += ZipFileSuffix
	}
	err = zipDirectory(packConfig.RootCloneDirectory, archiveFilename)
	if err != nil {
		fmt.Printf("Error %+s zipping directory %+s to %+s\n", err, packConfig.RootCloneDirectory, archiveFilename)
	}
}

func loadJSON(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		return err
	}

	return nil
}

// Stackoverflow provideth
// https://stackoverflow.com/questions/37869793/how-do-i-zip-a-directory-containing-sub-directories-or-files-in-golang
func zipDirectory(directory string, zipFilename string) error {
	file, err := os.Create(zipFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	return filepath.Walk(directory, walker)
}
