# git-download
git-download is download the GitHub releases.

## Usage
### Basic usage
```bash
# download neko-neko/git-download
$ git download -repo neko-neko/git-download

# specific version
$ git download -repo neko-neko/git-download -version v1.0.0

# specific files
$ git download -repo neko-neko/git-download -include darwin

# save files inside specific directory
$ git download -repo neko-neko/git-download -dir /tmp
```

### Private repo
```bash
$ GITHUB_TOKEN=hoge git download -repo neko-neko/git-download
```

### For GitHub Enterprise
```bash
$ GITHUB_TOKEN=hoge GITHUB_API=example.ghe.com/api/v3 git download -repo neko-neko/git-download
```

## Installation
```bash
$ go get github.com/neko-neko/git-download
```

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## Credits
neko-neko

## License
MIT