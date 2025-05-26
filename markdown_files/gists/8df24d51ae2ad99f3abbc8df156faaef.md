# [Golang GitHub API Client] #go #golang #github #api

## 3. Golang GitHub API Client (added flags).go

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

var (
	bold  = color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
	title = color.New(color.FgHiRed).Add(color.Underline).Add(color.Bold).SprintFunc()
)

func main() {
	md := flag.Bool("missing", false, "Display repos missing descriptions")
	ds := flag.Bool("squatters", false, "Display repos that are domain squatting")
	flag.Parse()

	client := NewGitHubClient()
	org := "fastly"

	if !*md && !*ds {
		flag.PrintDefaults()
	}

	var wg sync.WaitGroup
	if *md {
		wg.Add(1)
		go DisplayReposMissingDescription(client, org, &wg)
	}
	if *ds {
		wg.Add(1)
		go DisplayDomainSquatters(client, org, &wg)
	}
	wg.Wait()
}

// NewGitHubClient returns a client for making API requests to GitHub.
//
// NOTE:
// Ensure you have an API token exposed via the environment variable:
// GITHUB_TOKEN_TERMINAL_TOOL
func NewGitHubClient() *github.Client {
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv("GITHUB_TOKEN_TERMINAL_TOOL"),
		},
	)
	c := oauth2.NewClient(context.Background(), sts)

	return github.NewClient(c)
}

// NewOptions configures the ListByOrg API endpoint.
func NewOptions() *github.RepositoryListByOrgOptions {
	return &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
}

// DisplayReposMissingDescription prints all the repos that have no
// description, including whether the repo is private or public.
//
// NOTE:
// This function is run concurrently with DisplayDomainSquatters.
// To prevent interweaving of printed output we write data to a buffer which is
// then flushed once the search operation is finished.
func DisplayReposMissingDescription(client *github.Client, org string, wg *sync.WaitGroup) {
	var b strings.Builder

	fmt.Fprintf(&b, "%s\n\n", title("Repos with no description..."))

	opts := NewOptions()
	count := 0

	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, opts)
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range repos {
			if r.GetDescription() == "" {
				count += 1
				fmt.Fprintf(&b, "%s: %s\n%s %t\n\n", bold("URL"), r.GetHTMLURL(), bold("Private?"), r.GetPrivate())
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	fmt.Fprintf(&b, "Total repos without a description: %d\n\n", count)
	fmt.Println(b.String())

	wg.Done()
}

// DisplayDomainSquatters prints all the repos that contain only a single
// README.md file, including whether the repo is private or public.
//
// NOTE:
// This function is run concurrently with DisplayReposMissingDescription.
// To prevent interweaving of printed output we write data to a buffer which is
// then flushed once the search operation is finished.
//
// TODO:
// Some repos have content in non-default branches. This suggests the repo
// might be WIP (Work in Progress). We should update the script to inspect the
// last modified date for the branch. If modified recently (e.g. within the
// last month), then it's probably safe to keep the repo. Otherwise we should
// reach out to the team/person responsible for the repo about removing it.
func DisplayDomainSquatters(client *github.Client, org string, wg *sync.WaitGroup) {
	var b strings.Builder

	fmt.Fprintf(&b, "%s\n\n", title("Repos that are domain squatting..."))

	opts := NewOptions()
	count := 0

	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), "fastly", opts)
		if err != nil {
			log.Fatal(err)
		}
		sort.Slice(repos, func(i, j int) bool {
			return repos[i].GetName() < repos[j].GetName()
		})
		for _, r := range repos {
			_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), org, r.GetName(), "/", nil)
			if err != nil {
				if !strings.Contains(err.Error(), "This repository is empty") {
					log.Fatal(err)
				}
			}
			if len(directoryContent) == 1 && strings.ToLower(directoryContent[0].GetName()) == "readme.md" {
				count += 1
				fmt.Fprintf(&b, "%s: %s\n%s %t\n\n", bold("URL"), r.GetHTMLURL(), bold("Private?"), r.GetPrivate())
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	fmt.Fprintf(&b, "Total repos that are domain squatting: %d", count)
	fmt.Println(b.String())

	wg.Done()
}
```

## 4. Golang GitHub API Client (added concurrency).go

```go
// NOTE: GitHub API rate limit is 5,000 requests per hour for 'user-to-server'.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/fatih/color"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

// perPage is the number of items per page to be loaded from the API.
const perPage = 100

var (
	bold  = color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
	title = color.New(color.FgHiRed).Add(color.Underline).Add(color.Bold).SprintFunc()
)

func main() {
	md := flag.Bool("missing", false, "Display repos missing descriptions")
	ds := flag.Bool("squatters", false, "Display repos that are domain squatting")
	flag.Parse()

	client := NewGitHubClient()
	org := "fastly"

	if !*md && !*ds {
		flag.PrintDefaults()
	}

	pages := numOfPages(client, org)
	repos := fetchRepos(client, org, pages)

	var wg sync.WaitGroup
	if *md {
		wg.Add(1)
		go DisplayReposMissingDescription(repos, &wg)
	}
	if *ds {
		wg.Add(1)
		go DisplayDomainSquatters(repos, client, org, &wg)
	}
	wg.Wait()
}

// NewGitHubClient returns a client for making API requests to GitHub.
//
// NOTE:
// Ensure you have an API token exposed via the environment variable:
// GITHUB_TOKEN_TERMINAL_TOOL
func NewGitHubClient() *github.Client {
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv("GITHUB_TOKEN_TERMINAL_TOOL"),
		},
	)
	c := oauth2.NewClient(context.Background(), sts)

	return github.NewClient(c)
}

// numOfPages calculates the number of pages that need to be paginated based on
// the total public and private repos across the `perPage` constant.
func numOfPages(client *github.Client, org string) int {
	o, _, err := client.Organizations.Get(context.Background(), org)
	if err != nil {
		log.Fatal(err)
	}
	return (o.GetPublicRepos() + o.GetOwnedPrivateRepos()) / perPage
}

// fetchRepos makes multiple concurrent API requests for repositories.
func fetchRepos(client *github.Client, org string, numOfPages int) []*github.Repository {
	var m sync.Mutex
	repos := []*github.Repository{}

	process := func(page int, wg *sync.WaitGroup) {
		opts := NewOptions()
		opts.Page = page

		repositories, _, err := client.Repositories.ListByOrg(context.Background(), org, opts)
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range repositories {
			m.Lock()
			repos = append(repos, r)
			m.Unlock()
		}
		wg.Done()
	}

	var wg sync.WaitGroup
	wg.Add(numOfPages)
	for i := 1; i <= numOfPages; i++ {
		go process(i, &wg)
	}
	wg.Wait()

	return repos
}

// NewOptions configures the ListByOrg API endpoint.
func NewOptions() *github.RepositoryListByOrgOptions {
	return &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: perPage,
		},
	}
}

// DisplayReposMissingDescription prints all the repos that have no
// description, including whether the repo is private or public.
//
// NOTE:
// This function is run concurrently with DisplayDomainSquatters.
// To prevent interweaving of printed output we write data to a buffer which is
// then flushed once the search operation is finished.
//
// TODO:
// The two Display* functions both need to paginate over the same amount of
// pages, so only do it once and then store off the data to pass into both
// functions instead of repeating the operations.
func DisplayReposMissingDescription(repos []*github.Repository, wg *sync.WaitGroup) {
	var (
		b     strings.Builder
		count int
	)

	fmt.Fprintf(&b, "%s\n\n", title("Repos with no description..."))

	for _, r := range repos {
		if r.GetDescription() == "" {
			count += 1
			fmt.Fprintf(&b, "%s: %s\n%s %t\n\n", bold("URL"), r.GetHTMLURL(), bold("Private?"), r.GetPrivate())
		}
	}

	fmt.Fprintf(&b, "Total repos without a description: %d\n\n", count)
	fmt.Println(b.String())

	wg.Done()
}

// DisplayDomainSquatters prints all the repos that contain only a single
// README.md file, including whether the repo is private or public.
//
// NOTE:
// This function is run concurrently with DisplayReposMissingDescription.
// To prevent interweaving of printed output we write data to a buffer which is
// then flushed once the search operation is finished.
//
// Additionally we execute the API requests for repo file contents
// concurrently, and we use a semaphore pattern to batch those requests
// otherwise the GitHub API will issue us a secondary rate limit message.
//
// TODO:
// Some repos have content in non-default branches. This suggests the repo
// might be WIP (Work in Progress). We should update the script to inspect the
// last modified date for the branch. If modified recently (e.g. within the
// last month), then it's probably safe to keep the repo. Otherwise we should
// reach out to the team/person responsible for the repo about removing it.
func DisplayDomainSquatters(repos []*github.Repository, client *github.Client, org string, wg *sync.WaitGroup) {
	var (
		b     strings.Builder
		count uint64
		m     sync.Mutex
	)

	fmt.Fprintf(&b, "%s\n\n", title("Repos that are domain squatting..."))

	semaphore := make(chan struct{}, 50) // 50 is the maximum number of concurrent processes that may run at any time

	process := func(r *github.Repository, wg2 *sync.WaitGroup) {
		// Will block once more than 100 instances of `process` are called.
		semaphore <- struct{}{}

		// Empty buffered channel by one so the next call to `process` can begin.
		defer func() { <-semaphore }()

		// Done must be called before pulling from the semaphore channel.
		defer wg2.Done()

		_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), org, r.GetName(), "/", nil)
		if err != nil {
			if !strings.Contains(err.Error(), "This repository is empty") {
				log.Fatal(err)
			}
		}
		if len(directoryContent) == 1 && strings.ToLower(directoryContent[0].GetName()) == "readme.md" {
			atomic.AddUint64(&count, 1)

			// The strings.Builder type is not thread-safe.
			m.Lock()
			fmt.Fprintf(&b, "%s: %s\n%s %t\n\n", bold("URL"), r.GetHTMLURL(), bold("Private?"), r.GetPrivate())
			m.Unlock()
		}
	}

	var wg2 sync.WaitGroup
	wg2.Add(len(repos))
	for _, r := range repos {
		go process(r, &wg2)
	}
	wg2.Wait()
	close(semaphore)

	fmt.Fprintf(&b, "Total repos that are domain squatting: %d", count)
	fmt.Println(b.String())

	wg.Done()
}
```

## 1. Golang GitHub API Client (simple inlined proof-of-concept).go

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN_TERMINAL_TOOL")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// orgs, _, err := client.Organizations.List(context.Background(), "integralist", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("orgs: %+v\n", orgs)

	// opt := &github.RepositoryListByOrgOptions{Type: "private"}
	// repos, _, err := client.Repositories.ListByOrg(context.Background(), "fastly", opt)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("private repos: %d\n", len(repos))

	bold := color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
	title := color.New(color.FgHiRed).Add(color.Underline).Add(color.Bold)
	title.Print("Repos with no description...\n\n")

	// repos, _, err := client.Repositories.ListByOrg(context.Background(), "fastly", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, repo := range repos {
	// 	if repo.GetDescription() == "" {
	// 		fmt.Printf("%s: %s\n%s %t\n\n", bold("URL"), repo.GetHTMLURL(), bold("Private?"), repo.GetPrivate())
	// 	}
	// }

	org := "fastly"
	opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	count := 0

	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, opts)
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range repos {
			if r.GetDescription() == "" {
				count += 1
				fmt.Printf("%s: %s\n%s %t\n\n", bold("URL"), r.GetHTMLURL(), bold("Private?"), r.GetPrivate())
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	fmt.Printf("Total repos without a description: %d\n\n", count)

	title.Print("Repos that are domain squatting...\n\n")

	opts = &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	count = 0

	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, opts)
		if err != nil {
			log.Fatal(err)
		}
		sort.Slice(repos, func(i, j int) bool {
			return repos[i].GetName() < repos[j].GetName()
		})
		for _, r := range repos {
			_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), org, r.GetName(), "/", nil)
			if err != nil {
				if !strings.Contains(err.Error(), "This repository is empty") {
					log.Fatal(err)
				}
			}
			if len(directoryContent) == 1 && strings.ToLower(directoryContent[0].GetName()) == "readme.md" {
				count += 1
				fmt.Printf("%s: %s\n%s %t\n\n", bold("URL"), r.GetHTMLURL(), bold("Private?"), r.GetPrivate())
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	fmt.Printf("Total repos that are domain squatting: %d", count)
}
```

## 2. Golang GitHub API Client (improved with refactor + simple concurrency).go

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

var (
	bold  = color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
	title = color.New(color.FgHiRed).Add(color.Underline).Add(color.Bold).SprintFunc()
)

func main() {
	client := NewGitHubClient()
	org := "fastly"

	var wg sync.WaitGroup
	wg.Add(2)

	go DisplayReposMissingDescription(client, org, &wg)
	go DisplayDomainSquatters(client, org, &wg)

	wg.Wait()
}

// NewGitHubClient returns a client for making API requests to GitHub.
//
// NOTE:
// Ensure you have an API token exposed via the environment variable:
// GITHUB_TOKEN_TERMINAL_TOOL
func NewGitHubClient() *github.Client {
	sts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv("GITHUB_TOKEN_TERMINAL_TOOL"),
		},
	)
	c := oauth2.NewClient(context.Background(), sts)

	return github.NewClient(c)
}

// NewOptions configures the ListByOrg API endpoint.
func NewOptions() *github.RepositoryListByOrgOptions {
	return &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
}

// DisplayReposMissingDescription prints all the repos that have no
// description, including whether the repo is private or public.
//
// NOTE:
// This function is run concurrently with DisplayDomainSquatters.
// To prevent interweaving of printed output we write data to a buffer which is
// then flushed once the search operation is finished.
func DisplayReposMissingDescription(client *github.Client, org string, wg *sync.WaitGroup) {
	var b strings.Builder

	fmt.Fprintf(&b, "%s\n\n", title("Repos with no description..."))

	opts := NewOptions()
	count := 0

	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, opts)
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range repos {
			if r.GetDescription() == "" {
				count += 1
				fmt.Fprintf(&b, "%s: %s\n%s %t\n\n", bold("URL"), r.GetHTMLURL(), bold("Private?"), r.GetPrivate())
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	fmt.Fprintf(&b, "Total repos without a description: %d\n\n", count)
	fmt.Println(b.String())

	wg.Done()
}

// DisplayDomainSquatters prints all the repos that contain only a single
// README.md file, including whether the repo is private or public.
//
// NOTE:
// This function is run concurrently with DisplayReposMissingDescription.
// To prevent interweaving of printed output we write data to a buffer which is
// then flushed once the search operation is finished.
//
// TODO:
// Some repos have content in non-default branches. This suggests the repo
// might be WIP (Work in Progress). We should update the script to inspect the
// last modified date for the branch. If modified recently (e.g. within the
// last month), then it's probably safe to keep the repo. Otherwise we should
// reach out to the team/person responsible for the repo about removing it.
func DisplayDomainSquatters(client *github.Client, org string, wg *sync.WaitGroup) {
	var b strings.Builder

	fmt.Fprintf(&b, "%s\n\n", title("Repos that are domain squatting..."))

	opts := NewOptions()
	count := 0

	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), "fastly", opts)
		if err != nil {
			log.Fatal(err)
		}
		sort.Slice(repos, func(i, j int) bool {
			return repos[i].GetName() < repos[j].GetName()
		})
		for _, r := range repos {
			_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), org, r.GetName(), "/", nil)
			if err != nil {
				if !strings.Contains(err.Error(), "This repository is empty") {
					log.Fatal(err)
				}
			}
			if len(directoryContent) == 1 && strings.ToLower(directoryContent[0].GetName()) == "readme.md" {
				count += 1
				fmt.Fprintf(&b, "%s: %s\n%s %t\n\n", bold("URL"), r.GetHTMLURL(), bold("Private?"), r.GetPrivate())
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	fmt.Fprintf(&b, "Total repos that are domain squatting: %d", count)
	fmt.Println(b.String())

	wg.Done()
}
```

