package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"gopkg.in/yaml.v3"
)

// PackageOptions represent the options presented to users to be able to configure the way a package is built using these scripts
// The YAML that corresponds to these options are stored within packages/<package-name>/package.yaml for each package
type PackageOptions struct {
	// Version represents the version of the package. It will override other values if it exists
	Version *string `yaml:"version,omitempty"`
	// PackageVersion represents the current version of the package. It needs to be incremented whenever there are changes
	PackageVersion *int `yaml:"packageVersion" default:"0"`
	// MainChartOptions represent options presented to the user to configure the main chart
	MainChartOptions ChartOptions `yaml:",inline"`
	// AdditionalChartOptions represent options presented to the user to configure any additional charts
	AdditionalChartOptions []AdditionalChartOptions `yaml:"additionalCharts,omitempty"`
	// DoNotRelease represents a boolean flag that indicates a package should not be tracked in make charts
	DoNotRelease bool `yaml:"doNotRelease,omitempty"`
}

// ChartOptions represent the options presented to users to be able to configure the way a main chart is built using these scripts
type ChartOptions struct {
	// WorkingDir is the working directory for this chart within packages/<package-name>
	WorkingDir string `yaml:"workingDir" default:"charts"`
	// UpstreamOptions is any options provided on how to get this chart from upstream
	UpstreamOptions UpstreamOptions `yaml:",inline"`
	// IgnoreDependencies drops certain dependencies from the list that is parsed from upstream
	IgnoreDependencies []string `yaml:"ignoreDependencies"`
	// ReplacePaths marks paths as those that should be replaced instead of patches. Consequently, these paths will exist in both generated-changes/excludes and generated-changes/overlay
	ReplacePaths []string `yaml:"replacePaths"`
}

// AdditionalChartOptions represent the options presented to users to be able to configure the way an additional chart is built using these scripts
type AdditionalChartOptions struct {
	// WorkingDir is the working directory for this chart within packages/<package-name>
	WorkingDir string `yaml:"workingDir"`
	// UpstreamOptions is any options provided on how to get this chart from upstream. It is mutually exclusive with CRDChartOptions
	UpstreamOptions *UpstreamOptions `yaml:"upstreamOptions,omitempty"`
	// CRDChartOptions is any options provided on how to generate a CRD chart. It is mutually exclusive with UpstreamOptions
	CRDChartOptions *CRDChartOptions `yaml:"crdOptions,omitempty"`
	// IgnoreDependencies drops certain dependencies from the list that is parsed from upstream
	IgnoreDependencies []string `yaml:"ignoreDependencies"`
	// ReplacePaths marks paths as those that should be replaced instead of patches. Consequently, these paths will exist in both generated-changes/excludes and generated-changes/overlay
	ReplacePaths []string `yaml:"replacePaths"`
}

// UpstreamOptions represents the options presented to users to define where the upstream Helm chart is located
type UpstreamOptions struct {
	// URL represents a source for your upstream (e.g. a Github repository URL or a download link for an archive)
	URL string `yaml:"url,omitempty"`
	// Subdirectory represents a specific directory within the upstream pointed to by the URL to treat as the root
	Subdirectory *string `yaml:"subdirectory,omitempty"`
	// Commit represents a specific commit hash to treat as the head, if the URL points to a Github repository
	Commit *string `yaml:"commit,omitempty"`
}

// CRDChartOptions represent any options that are configurable for CRD charts
type CRDChartOptions struct {
	// The directory within packages/<package-name>/templates/ that will contain the template for your CRD chart
	TemplateDirectory string `yaml:"templateDirectory"`
	// The directory within your templateDirectory in which CRD files should be placed
	CRDDirectory string `yaml:"crdDirectory" default:"templates"`
	// Whether to add a validation file to your main chart to check that CRDs exist
	AddCRDValidationToMainChart bool `yaml:"addCRDValidationToMainChart"`
	// UseTarArchive indicates whether to bundle and compress CRD files into a tgz file
	UseTarArchive bool `yaml:"useTarArchive"`
}

type chartsVersion struct {
	Name          string
	DevBranch     string
	ReleaseBranch string
}

type gatheredVersion struct {
	PackageName    string
	VersionName    string
	DevVersion     string
	ReleaseVersion string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("must pass package name")
	}
	if err := printVersionsForPackage(os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printVersionsForPackage(packageName string) error {
	rawVersions := []string{"v2.6", "v2.7", "v2.8", "v2.9"}
	inputVersions := make([]chartsVersion, 0, len(rawVersions))
	for _, rawVersion := range rawVersions {
		version := chartsVersion{
			Name:          rawVersion,
			DevBranch:     "dev-" + rawVersion,
			ReleaseBranch: "release-" + rawVersion,
		}
		inputVersions = append(inputVersions, version)
	}

	fs := memfs.New()
	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: "https://github.com/rancher/charts",
	})
	if err != nil {
		return fmt.Errorf("failed to get repo %w", err)
	}
	err = repo.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		return fmt.Errorf("failed to fetch refs: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	outputVersions := make([]gatheredVersion, 0, len(inputVersions))
	for _, inputVersion := range inputVersions {
		devVersion, err := getVersion(worktree, packageName, inputVersion.DevBranch)
		if err != nil {
			return fmt.Errorf("failed to get dev version of %q: %w", inputVersion.Name, err)
		}
		releaseVersion, err := getVersion(worktree, packageName, inputVersion.ReleaseBranch)
		if err != nil {
			return fmt.Errorf("failed to get release version of %q: %w", inputVersion.Name, err)
		}
		outputVersion := gatheredVersion{
			PackageName:    packageName,
			VersionName:    inputVersion.Name,
			DevVersion:     devVersion,
			ReleaseVersion: releaseVersion,
		}
		outputVersions = append(outputVersions, outputVersion)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	fmt.Fprintf(writer, "\tdev\trelease\n")
	for _, outputVersion := range outputVersions {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", outputVersion.VersionName, outputVersion.DevVersion, outputVersion.ReleaseVersion)
	}
	writer.Flush()

	return nil
}

func getVersion(worktree *git.Worktree, packageName, branchName string) (string, error) {
	branch := "refs/heads/" + branchName
	err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branch),
	})
	if err != nil {
		return "", fmt.Errorf("failed to check out branch %q: %w", branch, err)
	}

	fd, err := worktree.Filesystem.Open("packages/" + packageName + "/package.yaml")
	if err != nil {
		return "", fmt.Errorf("failed to open package.yaml: %w", err)
	}
	defer fd.Close()

	pkg := new(PackageOptions)
	decoder := yaml.NewDecoder(fd)
	if err := decoder.Decode(&pkg); err != nil {
		return "", fmt.Errorf("failed to decode package.yaml on branch %q: %w", branchName, err)
	}
	return *pkg.Version, nil
}
