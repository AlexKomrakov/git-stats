package main

import (
    "github.com/stretchr/graceful"
    "github.com/gorilla/mux"
    "github.com/codegangsta/negroni"
    "github.com/unrolled/render"

    "gopkg.in/src-d/go-git.v3"
    "gopkg.in/src-d/go-git.v3/core"
    githttp "gopkg.in/src-d/go-git.v3/clients/http"

    "flag"
    "time"
    "log"
    "net/http"
    "io"
)

var serverAddress string
var r *render.Render

type RepoInfo struct {
    Commits int
    Authors []AuthorInfo
}

type AuthorInfo struct {
    Commits int
    Email string
}

func init() {
    r = render.New(render.Options {
        // Specify extensions to load for templates.
        Extensions: []string{".tmpl", ".html"},
        // Render will now recompile the templates on every HTML response.
        IsDevelopment: true,
    })
}

func Index(res http.ResponseWriter, req *http.Request) {
    r.HTML(res, http.StatusOK, "index", map[string]interface{}{
        "year": time.Now().Year(),
    })
}

func Stats(res http.ResponseWriter, req *http.Request) {
    repo := req.FormValue("repo")
    user := req.FormValue("user")
    pass := req.FormValue("pass")

    r.JSON(res, http.StatusOK, map[string]interface{}{
        "refs": getRefs(repo, user, pass),
        "commits": getCommits(repo, user, pass),
    })
}

func getRefs(repo, user, pass string) map[string]core.Hash {
    auth := githttp.NewBasicAuth(user, pass)
    r, err := git.NewAuthenticatedRemote(repo, auth)
    if err != nil {
        panic(err)
    }

    err = r.Connect()
    if err != nil {
        panic(err)
    }

    return r.Refs()
}

func getCommits(repo, user, pass string) RepoInfo {
    auth := githttp.NewBasicAuth(user, pass)
    repository, err := git.NewRepository(repo, auth)
    if err != nil {
        panic(err)
    }
    if err := repository.PullDefault(); err != nil {
        panic(err)
    }
    commits, _ := repository.Commits()
    defer commits.Close()
    authors_commits := make(map[string]int)
    commits_count := 0
    for {
        // The commits are not shorted in any special order
        commit, err := commits.Next()
        if err != nil {
            if err == io.EOF {
                break
            }
            panic(err)
        }
        authors_commits[commit.Author.Email]++
        commits_count++
    }

    repoInfo := RepoInfo{Commits: commits_count}
    for email, commits := range authors_commits {
        repoInfo.Authors = append(repoInfo.Authors, AuthorInfo{commits, email})
    }

    return repoInfo
}

func main() {
    flag.StringVar(&serverAddress, "a", "0.0.0.0:8000", "Server address: host:port")
    flag.Parse()

    router := mux.NewRouter()
    router.NotFoundHandler = http.HandlerFunc(Index)
    router.HandleFunc("/", Index)
    router.HandleFunc("/stats", Stats).Methods("GET")

    n := negroni.Classic()
    n.UseHandler(router)
    log.Println("Server started on address: " + serverAddress)
    graceful.Run(serverAddress, 30*time.Second, n)
}