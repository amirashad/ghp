# ghp
GitHub repository branch protection manager

You can use this tool to add or remove some protection to some repository's branch.

## Installation
Download ghp cli tool for OSX and set as executable: 
<pre>curl -L https://github.com/amirashad/ghp/releases/download/v0.0.1/ghp_darwin_amd64 -o /usr/local/bin/ghp
chmod +x /usr/local/bin/ghp</pre>

## Setup
To create Github personal token please follow: https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line

Add to your environment path: GITHUB_TOKEN={token-created-with-github-ui}

## Example
To add protection to repo's branch: 
<pre>ghp --org Some-Org --operation add --repos "repo1 repo2" --branches "develop master" --protection "security/snyk - Dockerfile (...), security/snyk - build.gradle (...)"</pre>

