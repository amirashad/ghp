package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v25/github"
	"github.com/jessevdk/go-flags"
	"golang.org/x/oauth2"
)

var appVersion = "v0.0.1"

var opts struct {
	Version        bool   `long:"version" description:"Show version"`
	OrgName        string `short:"o" long:"org" description:"Organization name"`
	RepoNames      string `short:"r" long:"repos" description:"Repository names. Can be multiple"`
	BranchNames    string `short:"b" long:"branches" default:"develop master" description:"Protected branch names. Can be multiple"`
	OperationName  string `short:"p" long:"operation" default:"add" description:"Operation name [add, remove]"`
	ProtectionName string `short:"t" long:"protection" description:"Protection name. Only one name allowed"`
}

func main() {
	flags.Parse(&opts)
	if opts.Version {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present. Please, add GITHUB_TOKEN environment variable")
	}

	if len(opts.OrgName) == 0 {
		log.Fatal("org must be described!")
	}
	if len(opts.RepoNames) == 0 {
		log.Fatal("repos must be described!")
	}
	if len(opts.BranchNames) == 0 {
		log.Fatal("branches must be described!")
	}
	if len(opts.OperationName) == 0 {
		log.Fatal("operation must be described!")
	}
	if len(opts.ProtectionName) == 0 {
		log.Fatal("protection must be described!")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_repos := strings.Split(opts.RepoNames, " ")
	_branches := strings.Split(opts.BranchNames, " ")
	for _, repoName := range _repos {
		fmt.Println(repoName)

		branches, _, err := client.Repositories.ListBranches(ctx, opts.OrgName, repoName, &github.ListOptions{PerPage: 100})
		if err == nil {
			for _, branch := range branches {

				if contains(_branches, *branch.Name) {
					fmt.Println("\t", *branch.Name)

					protection, _, err := client.Repositories.GetBranchProtection(ctx, opts.OrgName, repoName, *branch.Name)
					if err == nil {
						preq := &github.ProtectionRequest{
							RequiredStatusChecks: protection.RequiredStatusChecks,
							EnforceAdmins:        protection.EnforceAdmins.Enabled,
							Restrictions: &github.BranchRestrictionsRequest{
								Teams: []string{},
								Users: []string{},
							},
						}

						if protection.RequiredPullRequestReviews != nil {
							preq.RequiredPullRequestReviews = &github.PullRequestReviewsEnforcementRequest{
								DismissStaleReviews:          protection.RequiredPullRequestReviews.DismissStaleReviews,
								RequireCodeOwnerReviews:      protection.RequiredPullRequestReviews.RequireCodeOwnerReviews,
								RequiredApprovingReviewCount: protection.RequiredPullRequestReviews.RequiredApprovingReviewCount,
							}
						}

						for _, v := range protection.Restrictions.Users {
							preq.Restrictions.Users = append(preq.Restrictions.Users, *v.Login)
						}

						for _, v := range protection.Restrictions.Teams {
							preq.Restrictions.Teams = append(preq.Restrictions.Teams, *v.Slug)
						}

						if protection.RequiredStatusChecks == nil {
							protection.RequiredStatusChecks = &github.RequiredStatusChecks{
								Strict:   true,
								Contexts: []string{},
							}
							preq.RequiredStatusChecks = protection.RequiredStatusChecks
						} else {
							for _, v := range protection.RequiredStatusChecks.Contexts {
								preq.RequiredStatusChecks.Contexts = append(protection.RequiredStatusChecks.Contexts, v)
							}
						}

						if opts.OperationName == "add" {
							preq.RequiredStatusChecks.Contexts = append(preq.RequiredStatusChecks.Contexts, opts.ProtectionName)
						} else if opts.OperationName == "remove" {
							fmt.Println("Operation not supported yet: remove")
							return
						} else {
							fmt.Println("Operation not supported: ", opts.OperationName)
							return
						}
						fmt.Println("\t\tpreq.RequiredStatusChecks.Contexts ", preq.RequiredStatusChecks.Contexts)

						_, _, err := client.Repositories.UpdateBranchProtection(ctx, opts.OrgName, repoName, *branch.Name, preq)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println("\t\tDone!")
					}
				}
			}
		}
	}
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
