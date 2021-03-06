# GInit : Unique way of starting projects

## Overview [![GoDoc](https://godoc.org/github.com/hjoshi123/GInit?status.svg)](https://godoc.org/github.com/hjoshi123/ginit) [![Go Report Card](https://goreportcard.com/badge/github.com/hjoshi123/GInit)](https://goreportcard.com/report/github.com/hjoshi123/GInit)

GInit is a Command line tool built using Golang to start your project. Just enter your repo name and choose if u want private or public repo and Voilà; you have a repo and a folder (local directory) with git initialized. What's more is you can choose your `.gitignore` template for projects and it will be pushed along to your remote repo. It's a one stop tool to get your project up and going.

## Install

**Note**: To make this project work, the following environmental variable is necessary as GitHub has **deprecated** the Basic Auth mechanism. To create your own Personal Access Token refer [here](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token). \
**`export GITHUB_PAT=<Personal Access Token>`**

* To include it as a part of your project:

```go
go get github.com/hjoshi123/GInit
```

* To build the project from source:

```bash
git clone https://github.com/hjoshi123/GInit
go install
```

* Homebrew install coming soon

## Features

* Directory creation and `git init`.
* README.md and .gitignore part of Init Commit.
* .gitignore templates to be choosen from (currently 6, more coming soon).
* Command line arguments (WIP)

## Author

[Hemant Joshi](https://github.com/hjoshi123/)

### If this library helps you in anyway, show your love ❤️ by putting a ⭐ on this project ✌️

## License

This project is licensed under MIT - see the [LICENSE](https://github.com/hjoshi123/GInit/blob/master/LICENSE) file for more details.
