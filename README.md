# GitPacker
Command line tool for packaging Git repo(s) for long term storage.

This tool works best when you have configured SSH authentication between your device and Github/the source code hosting provider. Other authentication methods such as username/password will work, but require more input from the user (at a minimum once per repo to pack).

## Usage

### Installing from source

Download source code for desired version;
```bash
git clone git@github.com:tozny/GitPacker.git --branch 0.0.1
```

Build binary and move to a location in your shell path;
```bash
cd GitPacker
make build
mv gitpacker /SOME/PATH/IN/YOUR/SHELL/PATH/
```

### Running


### Executing using parameters from a file

Create a file named `pack.json` in the same directory with the appropriate values for the repos you want to clone and packaging parameters.

```json
{
	// Name of the directory to store all cloned repos in
	"root_clone_directory": "archive",
	// Array of cloning config for repos to clone
	"repos": [
		{
			"clone_directory": "hook-service",
			"git_url": "git@github.com:tozny/hook-service.git",
			// Leave blank if latest is find
			"commit": "be36802ee0f09919e696208ed01a18e2320d395f",
			// Defaults to false
			"shallow": true
		}
	],
	// Whether to generate zip file of all cloned repos, defaults to false
	"archive": true,
	"archive_filename": "packedByGitPacker"
}
```

### Executing using parameters via command line arguments

 #TODO

## Development

```bash
make all
```

```bash
make lint
```

```bash
make build
```

```bash
make run
```

## Publish

Follow [semantic versioning](https://semver.org) when releasing new versions of this program.

Releasing involves tagging a commit in this repository, and pushing the tag. Tagging and releasing of new versions should only be done from the master branch after an approved Pull Request has been merged, or on the branch of an approved Pull Request.

To publish a new version, run

```bash
make version version=X.Y.Z
```

To consume published updates from other repositories that depends on this module run

```bash
go get github.com/tozny/GitPacker@vX.Y.Z
```

and the go `get` tool will fetch the published artifact and update that modules `go.mod` and`go.sum` files with the updated dependency.
