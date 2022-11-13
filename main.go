package main

import (
	_ "embed"
	"fmt"
	"github.com/perbu/kaboose/localgit"
	"github.com/perbu/kaboose/manifest"
	"log"
	"os"
	"path/filepath"
)

//go:embed .version
var embeddedVersion string

const localConfig = "localconfig.yaml"

func realMain() error {
	var config manifest.LocalConfig
	err := manifest.Load("localconfig.yaml", &config)
	if err != nil {
		return fmt.Errorf("loading '%s': %w", localConfig, err)
	}
	config.ExpandPaths() // fix the tilde paths:
	// subCluster is the first argument to the program:
	if len(os.Args) < 2 {
		return fmt.Errorf("missing sub-cluster name")
	}
	targetCluster := os.Args[1]
	// find the sub-cluster in the config:
	subCluster, ok := config.SubClusters[targetCluster]
	if !ok {
		return fmt.Errorf("sub-cluster '%s' not found in '%s'", targetCluster, localConfig)
	}
	err = processCluster(config, subCluster)
	if err != nil {
		return fmt.Errorf("processing sub-cluster '%s': %w", targetCluster, err)
	}
	return nil
}

func main() {
	fmt.Println("Kaboose", embeddedVersion)
	err := realMain()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type Repo struct {
	name         string
	URL          string
	branch       string
	localgitRepo *localgit.LocalRepo
}

type repoMap map[string]Repo
type SubCluster struct {
	repos map[string]Repo
}

func processCluster(config manifest.LocalConfig, subCluster manifest.SubCluster) error {
	// first we need to load the manifest file:
	var manif manifest.Manifest
	sc := SubCluster{}
	sc.repos = make(repoMap)
	manifestPath := filepath.Join(subCluster.Path, ".n2-manifest.yaml")
	err := manifest.Load(manifestPath, &manif)
	if err != nil {
		return fmt.Errorf("loading manifest file '%s': %w", manifestPath, err)
	}
	// manifest is loaded. Let's see if we can load the repositories:
	for _, repo := range manif.Repositories {
		repoPath := filepath.Join(config.RepoBaseDir, repo.Name)
		gitRepo, err := localgit.NewLocalRepo(repoPath)
		if err != nil {
			// return fmt.Errorf("loading repository '%s': %w", repo.Name, err)
			log.Printf("loading git repo '%s': %v", repoPath, err)
			continue
		}
		fmt.Println("Loaded repo", gitRepo)
		sc.repos[repo.Name] = Repo{
			name:         repo.Name,
			URL:          repo.URL,
			branch:       repo.Branch,
			localgitRepo: gitRepo,
		}
	}

	return nil
}
