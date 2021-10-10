# Overview

Githubsearch uses the github API to search exact text matches, including special characters [that are ignored](https://docs.github.com/en/search-github/searching-on-github/searching-code#considerations-for-code-search) when searching github code.  This is done by replacing the special characters with a space, then checking matching file contents for an exact match.

Support for searching for usage of exported go resources is also available.  The is done by looking for the imported package name, then searching for usage of the given variable/function.

# Usage

A github token can be used for API calls by either providing the `GH_TOKEN` or `GITHUB_TOKEN` environment variable, or by using the `--token` command line option.

Positional parameters are used at the search text.

## Organization search

Searching within an organization is done with the `--organization` option.  Omitting this option searches public repositories

## Go package search

The `--searchtype=gopackage` option will check the given import path, followed by a `.`, then the resource.  `github.com/someusername/reponame/package.New` would search for files containing the import path `github.com/someusername/reponame/package`, and also containing `package.New()`.  Import aliases are used if provided.

# Example

To search for the text `/a/file/path`:

`githubsearch search /a/file/path`

To search for the text `one:two:three` in the `my-company` organization:

`githubsearch search one:two:three`

To search for usage of the `New` function in the go package `github.com/someusername/reponame/package`:


`githubsearch search --searchtype=gopackage github.com/someusername/reponame/package.New`
