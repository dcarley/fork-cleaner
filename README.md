# fork-cleaner

Utility to remove forked repos from GitHub that you no longer need.

There are other tools that you could use to do this but I couldn't find any
which were simple enough to grok all of the code and trust for something
that is potentially dangerous.

## Usage

1. Create a [GitHub personal access token][token] and export it as `GITHUB_ACCESS_TOKEN`.
1. Install dependencies:

        dep ensure

1. Run the utility:

        go run main.go

1. Answer the prompts about whether you want to delete each repo.

[token]: https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
