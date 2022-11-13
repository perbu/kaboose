package manifest

// LocalConfig contains the per-user configuration for kaboose.
type LocalConfig struct {
	RepoBaseDir string `yaml:"repoBaseDir"` // the base directory for all repos of the applications - e.g. ~/git
	// SubClusters are the various parts of the clusters containing sets of config files
	SubClusters map[string]SubCluster `yaml:"subClusters"`
}

func (c *LocalConfig) ExpandPaths() {
	c.RepoBaseDir = expandTilde(c.RepoBaseDir)
	for k, sc := range c.SubClusters {
		sc.Path = expandTilde(sc.Path)
		c.SubClusters[k] = sc
	}
}

type SubCluster struct {
	Path string `yaml:"path"`
}

// Manifest contains the cluster configuration
// This is the main configuration file for kaboose and resides alongside the k8s yaml
type Manifest struct {
	Config struct {
		WorkingDir       string `yaml:"workingDir"`
		ValuesFormat     string `yaml:"valuesFormat"`
		ValuesPath       string `yaml:"valuesPath"`
		DockerRepoPrefix string `yaml:"dockerRepoPrefix"`
		BuildSecrets     []struct {
			Repo       string `yaml:"repo"`
			ArgName    string `yaml:"argName"`
			SecretName string `yaml:"secretName"`
		} `yaml:"buildSecrets"`
	} `yaml:"config"`
	Repositories []struct {
		Name   string `yaml:"name"`
		URL    string `yaml:"url"`
		Branch string `yaml:"branch"`
	} `yaml:"repositories"`
}

// Kustomization contains the kustomization configuration
type Kustomization struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Images     []struct {
		Name   string `yaml:"name"`
		NewTag string `yaml:"newTag"`
	} `yaml:"images"`
	Resources []string `yaml:"resources"`
}
