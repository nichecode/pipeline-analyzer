package githubactions

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parser handles GitHub Actions workflow parsing
type Parser struct{}

// NewParser creates a new GitHub Actions parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses a GitHub Actions workflow file
func (p *Parser) ParseFile(filePath string) (*Workflow, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return p.Parse(data)
}

// Parse parses GitHub Actions workflow YAML content
func (p *Parser) Parse(data []byte) (*Workflow, error) {
	var workflow Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}

	return &workflow, nil
}

// ParseWorkflowsDirectory parses all workflow files in a directory
func (p *Parser) ParseWorkflowsDirectory(dirPath string) (map[string]*Workflow, error) {
	workflows := make(map[string]*Workflow)
	
	files, err := filepath.Glob(filepath.Join(dirPath, "*.yml"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob yml files: %v", err)
	}
	
	yamlFiles, err := filepath.Glob(filepath.Join(dirPath, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob yaml files: %v", err)
	}
	
	files = append(files, yamlFiles...)
	
	for _, file := range files {
		workflow, err := p.ParseFile(file)
		if err != nil {
			// Log warning but continue with other files
			continue
		}
		
		fileName := filepath.Base(file)
		workflows[fileName] = workflow
	}
	
	return workflows, nil
}

// ExtractRunCommands extracts run commands from a step
func (p *Parser) ExtractRunCommands(step Step) []string {
	if step.Run == "" {
		return []string{}
	}
	
	// Split multi-line run commands
	commands := strings.Split(step.Run, "\n")
	var cleanCommands []string
	
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" && !strings.HasPrefix(cmd, "#") {
			cleanCommands = append(cleanCommands, cmd)
		}
	}
	
	return cleanCommands
}

// GetJobDependencies extracts job dependencies from needs field
func (p *Parser) GetJobDependencies(job Job) []string {
	if job.Needs == nil {
		return []string{}
	}
	
	switch needs := job.Needs.(type) {
	case string:
		return []string{needs}
	case []interface{}:
		var deps []string
		for _, dep := range needs {
			if depStr, ok := dep.(string); ok {
				deps = append(deps, depStr)
			}
		}
		return deps
	case []string:
		return needs
	}
	
	return []string{}
}

// GetRunnerType extracts runner type from runs-on field
func (p *Parser) GetRunnerType(job Job) string {
	if job.RunsOn == nil {
		return "unknown"
	}
	
	switch runner := job.RunsOn.(type) {
	case string:
		return runner
	case []interface{}:
		if len(runner) > 0 {
			if runnerStr, ok := runner[0].(string); ok {
				return runnerStr
			}
		}
	case []string:
		if len(runner) > 0 {
			return runner[0]
		}
	}
	
	return "unknown"
}

// ExtractDockerImages finds Docker images used in jobs
func (p *Parser) ExtractDockerImages(workflow *Workflow) map[string][]string {
	images := make(map[string][]string)
	
	for jobName, job := range workflow.Jobs {
		var jobImages []string
		
		// Check container field
		if job.Container != nil {
			switch container := job.Container.(type) {
			case string:
				jobImages = append(jobImages, container)
			case map[interface{}]interface{}:
				if image, ok := container["image"].(string); ok {
					jobImages = append(jobImages, image)
				}
			}
		}
		
		// Check service containers
		for _, service := range job.Services {
			jobImages = append(jobImages, service.Image)
		}
		
		// Check for Docker commands in run steps
		for _, step := range job.Steps {
			commands := p.ExtractRunCommands(step)
			for _, cmd := range commands {
				dockerImages := p.extractDockerImagesFromCommand(cmd)
				jobImages = append(jobImages, dockerImages...)
			}
		}
		
		if len(jobImages) > 0 {
			images[jobName] = jobImages
		}
	}
	
	return images
}

// extractDockerImagesFromCommand finds Docker images in run commands
func (p *Parser) extractDockerImagesFromCommand(command string) []string {
	var images []string
	
	// Match docker run, docker pull, etc.
	dockerRegex := regexp.MustCompile(`docker\s+(?:run|pull|build|push)\s+(?:-[^\s]+\s+)*([^\s]+)`)
	matches := dockerRegex.FindAllStringSubmatch(command, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			images = append(images, match[1])
		}
	}
	
	return images
}