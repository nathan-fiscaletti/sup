package main

import (
	"log"
	"os"

	"github.com/nathan-fiscaletti/sup/internal"
	"github.com/nathan-fiscaletti/sup/internal/color"
)

func main() {
	isGitRepo, err := internal.IsGitRepository()
	if err != nil {
		log.Fatalf("Error checking if git repository: %v", err)
	}

	if !isGitRepo {
		log.Fatal("Not a git repository")
	}

	repo, err := internal.GetRepositoryDetails()
	if err != nil {
		log.Fatalf("Error getting repository ID: %v", err)
	}
	log.Printf("Repository: %s%s%s, MD5: %s%s%s\n", color.Green, repo.URL, color.Reset, color.Magenta, repo.UrlMd5Sum, color.Reset)

	cache, err := repo.GetChangeCache()
	if err != nil {
		log.Fatalf("Error getting cache location: %v", err)
	}

	isFirstLoad := false
	cachedChanges, err := cache.Get()
	if err != nil {
		isFirstLoad = os.IsNotExist(err)
		if !isFirstLoad {
			panic(err)
		}
	}

	_, err = internal.Git("fetch", "--all")
	if err != nil {
		panic(err)
	}

	changes, err := internal.GetRemoteUpdates()
	if err != nil {
		panic(err)
	}

	err = cache.Set(changes)
	if err != nil {
		panic(err)
	}

	if isFirstLoad {
		log.Println("Cache updated with most recent changes.")
		log.Println("After changes have been made to a branch, run again to see changes.")
		return
	}

	difference := changes.Compare(cachedChanges)
	if len(difference) == 0 {
		log.Println("No changes found since last run.")
		return
	}

	for _, change := range difference {
		log.Printf("%s%s%s\n", color.Cyan, change.BranchName, color.Reset)
		changes, err := change.Branch().CommitsSince(change.Date)
		if err != nil {
			log.Fatalf("Error getting commits for branch: %s", err)
		}

		for _, commit := range changes {
			subject, err := commit.Subject()
			if err != nil {
				log.Fatalf("Error getting commit subject: %v", err)
			}

			shortHash, err := commit.ShortHash()
			if err != nil {
				log.Fatalf("Error getting commit short hash: %v", err)
			}

			log.Printf("%s%s%s %s%s%s -- %s\n", color.Blue, commit.Author, color.Reset, color.Yellow, shortHash, color.Reset, subject)
		}
	}
}
