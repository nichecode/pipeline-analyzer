package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BuildTool represents a discovered build tool configuration
type BuildTool struct {
	Type        string `json:"type"`        // circleci, gotask, npm, etc.
	Name        string `json:"name"`        // Display name
	ConfigPath  string `json:"config_path"` // Path to config file
	Description string `json:"description"` // Tool description
}

// Repository represents a discovered repository structure
type Repository struct {
	RootPath   string      `json:"root_path"`   // Repository root path
	GitRepo    bool        `json:"git_repo"`    // Is this a git repository
	BuildTools []BuildTool `json:"build_tools"` // Discovered build tools
}

// Scanner handles repository scanning for build tools
type Scanner struct {
	rootPath string
}

// NewScanner creates a new repository scanner
func NewScanner(rootPath string) *Scanner {
	return &Scanner{
		rootPath: rootPath,
	}
}

// ScanRepository scans the repository for all build tools
func (s *Scanner) ScanRepository() (*Repository, error) {
	// Ensure we're scanning from a valid directory
	if _, err := os.Stat(s.rootPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", s.rootPath)
	}

	repo := &Repository{
		RootPath:   s.rootPath,
		GitRepo:    s.isGitRepository(),
		BuildTools: []BuildTool{},
	}

	// Scan for different build tools
	buildTools, err := s.discoverBuildTools()
	if err != nil {
		return nil, err
	}

	repo.BuildTools = buildTools
	return repo, nil
}

// isGitRepository checks if the root path is a git repository
func (s *Scanner) isGitRepository() bool {
	gitPath := filepath.Join(s.rootPath, ".git")
	_, err := os.Stat(gitPath)
	return err == nil
}

// discoverBuildTools discovers all build tools in the repository
func (s *Scanner) discoverBuildTools() ([]BuildTool, error) {
	var tools []BuildTool

	// Define build tool patterns to search for
	patterns := []struct {
		toolType    string
		name        string
		patterns    []string
		description string
	}{
		{
			toolType:    "circleci",
			name:        "CircleCI",
			patterns:    []string{".circleci/config.yml", ".circleci/config.yaml"},
			description: "CircleCI continuous integration",
		},
		{
			toolType:    "gotask",
			name:        "Go Task",
			patterns:    []string{"Taskfile.yml", "Taskfile.yaml"},
			description: "Go Task runner",
		},
		{
			toolType:    "github-actions",
			name:        "GitHub Actions",
			patterns:    []string{".github/workflows/*.yml", ".github/workflows/*.yaml"},
			description: "GitHub Actions workflows",
		},
		{
			toolType:    "npm",
			name:        "npm",
			patterns:    []string{"package.json"},
			description: "Node.js package manager",
		},
		{
			toolType:    "composer",
			name:        "Composer",
			patterns:    []string{"composer.json"},
			description: "PHP dependency manager",
		},
		{
			toolType:    "cargo",
			name:        "Cargo",
			patterns:    []string{"Cargo.toml"},
			description: "Rust package manager",
		},
		{
			toolType:    "maven",
			name:        "Maven",
			patterns:    []string{"pom.xml"},
			description: "Java build tool",
		},
		{
			toolType:    "gradle",
			name:        "Gradle",
			patterns:    []string{"build.gradle", "build.gradle.kts", "gradle.build"},
			description: "Java/Kotlin build tool",
		},
		{
			toolType:    "makefile",
			name:        "Makefile",
			patterns:    []string{"Makefile", "makefile", "GNUmakefile"},
			description: "GNU Make build system",
		},
		{
			toolType:    "docker",
			name:        "Docker",
			patterns:    []string{"Dockerfile", "docker-compose.yml", "docker-compose.yaml"},
			description: "Docker containerization",
		},
		{
			toolType:    "python",
			name:        "Python",
			patterns:    []string{"requirements.txt", "pyproject.toml", "setup.py", "Pipfile"},
			description: "Python package management",
		},
		{
			toolType:    "terraform",
			name:        "Terraform",
			patterns:    []string{"*.tf", "main.tf"},
			description: "Infrastructure as code",
		},
		{
			toolType:    "kubernetes",
			name:        "Kubernetes",
			patterns:    []string{"*.yaml", "*.yml"},
			description: "Kubernetes manifests",
		},
	}

	// Track found tools to avoid duplicates
	foundTools := make(map[string]bool)

	// Search for each pattern
	for _, pattern := range patterns {
		foundFiles, err := s.findFiles(pattern.patterns)
		if err != nil {
			return nil, err
		}

		// Add found tools, avoiding duplicates
		for _, file := range foundFiles {
			// Skip Kubernetes check for now as it's too broad
			if pattern.toolType == "kubernetes" {
				continue
			}

			// Create unique key for deduplication
			toolKey := pattern.toolType + ":" + file
			if foundTools[toolKey] {
				continue
			}
			foundTools[toolKey] = true

			tools = append(tools, BuildTool{
				Type:        pattern.toolType,
				Name:        pattern.name,
				ConfigPath:  file,
				Description: pattern.description,
			})
		}
	}

	return tools, nil
}

// findFiles searches for files matching the given patterns
func (s *Scanner) findFiles(patterns []string) ([]string, error) {
	var foundFiles []string

	for _, pattern := range patterns {
		// Handle glob patterns for workflows
		if strings.Contains(pattern, "*") {
			matches, err := filepath.Glob(filepath.Join(s.rootPath, pattern))
			if err != nil {
				return nil, err
			}
			for _, match := range matches {
				// Convert to relative path
				relPath, err := filepath.Rel(s.rootPath, match)
				if err == nil {
					foundFiles = append(foundFiles, relPath)
				}
			}
		} else {
			// Direct file check
			fullPath := filepath.Join(s.rootPath, pattern)
			if _, err := os.Stat(fullPath); err == nil {
				foundFiles = append(foundFiles, pattern)
			}
		}
	}

	return foundFiles, nil
}

// CreateDiscoveryDir creates the .discovery directory structure
func (s *Scanner) CreateDiscoveryDir() (string, error) {
	discoveryPath := filepath.Join(s.rootPath, ".discovery", "pipeline-analyzer")
	
	// Create the directory structure
	if err := os.MkdirAll(discoveryPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create discovery directory: %w", err)
	}

	return discoveryPath, nil
}

// GetOutputDir returns the appropriate output directory for a tool type
func (s *Scanner) GetOutputDir(baseDir, toolType string) string {
	return filepath.Join(baseDir, toolType)
}

// ScanAndCreateStructure scans the repository and creates the discovery structure
func (s *Scanner) ScanAndCreateStructure() (*Repository, string, error) {
	// Scan the repository
	repo, err := s.ScanRepository()
	if err != nil {
		return nil, "", err
	}

	// Create discovery directory
	discoveryDir, err := s.CreateDiscoveryDir()
	if err != nil {
		return nil, "", err
	}

	return repo, discoveryDir, nil
}