package docker

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseDockerfile parses a Dockerfile and returns its analysis
func ParseDockerfile(dockerfilePath string) (*DockerfileAnalysis, error) {
	content, err := readDockerfile(dockerfilePath)
	if err != nil {
		return nil, err
	}

	instructions := parseInstructions(content)
	stages := parseStages(instructions)
	
	analysis := &DockerfileAnalysis{
		FilePath:     dockerfilePath,
		Stages:       stages,
		MultiStage:   len(stages) > 1,
		BaseImages:   extractBaseImages(stages),
		ExposedPorts: extractExposedPorts(instructions),
		Labels:       extractLabels(instructions),
		Arguments:    extractArguments(instructions),
		Environment:  extractEnvironment(instructions),
		WorkingDir:   extractWorkingDir(instructions),
		User:         extractUser(instructions),
		Entrypoint:   extractEntrypoint(instructions),
		Command:      extractCommand(instructions),
		HealthCheck:  extractHealthCheck(instructions),
		Volumes:      extractVolumes(instructions),
	}

	// Perform security analysis
	analysis.SecurityScan = performSecurityScan(analysis)

	return analysis, nil
}

// readDockerfile reads the content of a Dockerfile
func readDockerfile(dockerfilePath string) ([]string, error) {
	file, err := os.Open(dockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Dockerfile: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	return lines, nil
}

// parseInstructions parses all instructions from Dockerfile content
func parseInstructions(lines []string) []*DockerfileInstruction {
	var instructions []*DockerfileInstruction
	var currentInstruction *DockerfileInstruction

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle line continuations
		if strings.HasSuffix(line, "\\") {
			if currentInstruction == nil {
				currentInstruction = parseInstructionLine(line, lineNum+1)
			} else {
				// Continue previous instruction
				currentInstruction.Raw += " " + strings.TrimSuffix(line, "\\")
				// Re-parse arguments
				currentInstruction.Arguments, currentInstruction.Flags = parseArguments(currentInstruction.Raw)
			}
			continue
		}

		if currentInstruction != nil {
			// Complete the multi-line instruction
			currentInstruction.Raw += " " + line
			currentInstruction.Arguments, currentInstruction.Flags = parseArguments(currentInstruction.Raw)
			instructions = append(instructions, currentInstruction)
			currentInstruction = nil
		} else {
			// Single-line instruction
			instruction := parseInstructionLine(line, lineNum+1)
			if instruction != nil {
				instructions = append(instructions, instruction)
			}
		}
	}

	// Handle case where file ends with continuation
	if currentInstruction != nil {
		instructions = append(instructions, currentInstruction)
	}

	return instructions
}

// parseInstructionLine parses a single instruction line
func parseInstructionLine(line string, lineNum int) *DockerfileInstruction {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}

	instruction := &DockerfileInstruction{
		Instruction: strings.ToUpper(parts[0]),
		Line:        lineNum,
		Raw:         line,
		Flags:       make(map[string]string),
	}

	instruction.Arguments, instruction.Flags = parseArguments(line)
	return instruction
}

// parseArguments extracts arguments and flags from an instruction
func parseArguments(line string) ([]string, map[string]string) {
	parts := strings.Fields(line)
	if len(parts) <= 1 {
		return []string{}, make(map[string]string)
	}

	args := []string{}
	flags := make(map[string]string)

	for i := 1; i < len(parts); i++ {
		part := parts[i]
		if strings.HasPrefix(part, "--") {
			// This is a flag
			if strings.Contains(part, "=") {
				flagParts := strings.SplitN(part, "=", 2)
				flags[flagParts[0]] = flagParts[1]
			} else {
				// Flag without value, check next part
				if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "--") {
					flags[part] = parts[i+1]
					i++ // Skip next part
				} else {
					flags[part] = "true"
				}
			}
		} else {
			args = append(args, part)
		}
	}

	return args, flags
}

// parseStages groups instructions into build stages
func parseStages(instructions []*DockerfileInstruction) []*DockerfileStage {
	var stages []*DockerfileStage
	var currentStage *DockerfileStage

	for _, instruction := range instructions {
		if instruction.Instruction == "FROM" {
			// Start a new stage
			if currentStage != nil {
				stages = append(stages, currentStage)
			}

			currentStage = &DockerfileStage{
				Instructions: []*DockerfileInstruction{instruction},
				Index:        len(stages),
			}

			// Extract base image and stage name
			if len(instruction.Arguments) > 0 {
				currentStage.BaseImage = instruction.Arguments[0]
				
				// Look for AS clause
				for i, arg := range instruction.Arguments {
					if strings.ToUpper(arg) == "AS" && i+1 < len(instruction.Arguments) {
						currentStage.Name = instruction.Arguments[i+1]
						break
					}
				}
			}
		} else if currentStage != nil {
			currentStage.Instructions = append(currentStage.Instructions, instruction)
		}
	}

	// Add the last stage
	if currentStage != nil {
		stages = append(stages, currentStage)
	}

	return stages
}

// Extract functions for various Dockerfile elements
func extractBaseImages(stages []*DockerfileStage) []string {
	var images []string
	seen := make(map[string]bool)

	for _, stage := range stages {
		if stage.BaseImage != "" && !seen[stage.BaseImage] {
			images = append(images, stage.BaseImage)
			seen[stage.BaseImage] = true
		}
	}

	return images
}

func extractExposedPorts(instructions []*DockerfileInstruction) []string {
	var ports []string
	for _, instruction := range instructions {
		if instruction.Instruction == "EXPOSE" {
			ports = append(ports, instruction.Arguments...)
		}
	}
	return ports
}

func extractLabels(instructions []*DockerfileInstruction) map[string]string {
	labels := make(map[string]string)
	for _, instruction := range instructions {
		if instruction.Instruction == "LABEL" {
			for _, arg := range instruction.Arguments {
				if strings.Contains(arg, "=") {
					parts := strings.SplitN(arg, "=", 2)
					if len(parts) == 2 {
						labels[parts[0]] = strings.Trim(parts[1], "\"'")
					}
				}
			}
		}
	}
	return labels
}

func extractArguments(instructions []*DockerfileInstruction) map[string]string {
	args := make(map[string]string)
	for _, instruction := range instructions {
		if instruction.Instruction == "ARG" {
			for _, arg := range instruction.Arguments {
				if strings.Contains(arg, "=") {
					parts := strings.SplitN(arg, "=", 2)
					if len(parts) == 2 {
						args[parts[0]] = parts[1]
					}
				} else {
					args[arg] = ""
				}
			}
		}
	}
	return args
}

func extractEnvironment(instructions []*DockerfileInstruction) map[string]string {
	env := make(map[string]string)
	for _, instruction := range instructions {
		if instruction.Instruction == "ENV" {
			for _, arg := range instruction.Arguments {
				if strings.Contains(arg, "=") {
					parts := strings.SplitN(arg, "=", 2)
					if len(parts) == 2 {
						env[parts[0]] = parts[1]
					}
				}
			}
		}
	}
	return env
}

func extractWorkingDir(instructions []*DockerfileInstruction) string {
	var workdir string
	for _, instruction := range instructions {
		if instruction.Instruction == "WORKDIR" && len(instruction.Arguments) > 0 {
			workdir = instruction.Arguments[0]
		}
	}
	return workdir
}

func extractUser(instructions []*DockerfileInstruction) string {
	var user string
	for _, instruction := range instructions {
		if instruction.Instruction == "USER" && len(instruction.Arguments) > 0 {
			user = instruction.Arguments[0]
		}
	}
	return user
}

func extractEntrypoint(instructions []*DockerfileInstruction) []string {
	for _, instruction := range instructions {
		if instruction.Instruction == "ENTRYPOINT" {
			return instruction.Arguments
		}
	}
	return nil
}

func extractCommand(instructions []*DockerfileInstruction) []string {
	for _, instruction := range instructions {
		if instruction.Instruction == "CMD" {
			return instruction.Arguments
		}
	}
	return nil
}

func extractHealthCheck(instructions []*DockerfileInstruction) *HealthCheck {
	for _, instruction := range instructions {
		if instruction.Instruction == "HEALTHCHECK" {
			hc := &HealthCheck{
				Command: []string{},
			}

			// Parse flags for healthcheck options
			if interval, ok := instruction.Flags["--interval"]; ok {
				if d, err := time.ParseDuration(interval); err == nil {
					hc.Interval = d
				}
			}
			if timeout, ok := instruction.Flags["--timeout"]; ok {
				if d, err := time.ParseDuration(timeout); err == nil {
					hc.Timeout = d
				}
			}
			if startPeriod, ok := instruction.Flags["--start-period"]; ok {
				if d, err := time.ParseDuration(startPeriod); err == nil {
					hc.StartPeriod = d
				}
			}
			if retries, ok := instruction.Flags["--retries"]; ok {
				if r, err := strconv.Atoi(retries); err == nil {
					hc.Retries = r
				}
			}

			hc.Command = instruction.Arguments
			return hc
		}
	}
	return nil
}

func extractVolumes(instructions []*DockerfileInstruction) []string {
	var volumes []string
	for _, instruction := range instructions {
		if instruction.Instruction == "VOLUME" {
			volumes = append(volumes, instruction.Arguments...)
		}
	}
	return volumes
}

// performSecurityScan analyzes the Dockerfile for security issues
func performSecurityScan(analysis *DockerfileAnalysis) *SecurityScan {
	scan := &SecurityScan{
		VulnerableBaseImages:    []string{},
		SecurityRecommendations: []string{},
		BestPracticeIssues:     []string{},
	}

	// Check if running as root
	scan.RunAsRoot = analysis.User == "" || analysis.User == "root"

	// Check for health check
	scan.HasHealthCheck = analysis.HealthCheck != nil

	// Check for latest tags
	scan.UsesLatestTags = checkForLatestTags(analysis.BaseImages)

	// Check for security updates
	scan.HasSecurityUpdates = checkForSecurityUpdates(analysis)

	// Check caching optimization
	scan.CachingOptimized = checkCachingOptimization(analysis)

	// Check multi-stage optimization
	scan.MultiStageOptimized = analysis.MultiStage

	// Generate recommendations
	generateSecurityRecommendations(scan, analysis)

	return scan
}

func checkForLatestTags(baseImages []string) bool {
	for _, image := range baseImages {
		if strings.Contains(image, ":latest") || !strings.Contains(image, ":") {
			return true
		}
	}
	return false
}

func checkForSecurityUpdates(analysis *DockerfileAnalysis) bool {
	// Look for security update patterns in RUN instructions
	securityUpdatePatterns := []string{
		"apt-get upgrade",
		"apk upgrade",
		"yum upgrade",
		"dnf upgrade",
	}

	for _, stage := range analysis.Stages {
		for _, instruction := range stage.Instructions {
			if instruction.Instruction == "RUN" {
				instructionText := strings.ToLower(instruction.Raw)
				for _, pattern := range securityUpdatePatterns {
					if strings.Contains(instructionText, pattern) {
						return true
					}
				}
			}
		}
	}
	return false
}

func checkCachingOptimization(analysis *DockerfileAnalysis) bool {
	// Check if package installation happens before copying source code
	for _, stage := range analysis.Stages {
		foundPackageInstall := false
		foundSourceCopy := false

		for _, instruction := range stage.Instructions {
			if instruction.Instruction == "RUN" {
				instructionText := strings.ToLower(instruction.Raw)
				if strings.Contains(instructionText, "install") ||
				   strings.Contains(instructionText, "pip ") ||
				   strings.Contains(instructionText, "npm ") {
					foundPackageInstall = true
				}
			}
			if instruction.Instruction == "COPY" && !foundSourceCopy {
				// Check if this looks like source code copy
				for _, arg := range instruction.Arguments {
					if strings.Contains(arg, "src") || strings.Contains(arg, ".") {
						foundSourceCopy = true
						break
					}
				}
			}
		}

		// If source is copied before package installation, caching is not optimized
		if foundSourceCopy && !foundPackageInstall {
			return false
		}
	}
	return true
}

func generateSecurityRecommendations(scan *SecurityScan, analysis *DockerfileAnalysis) {
	if scan.RunAsRoot {
		scan.SecurityRecommendations = append(scan.SecurityRecommendations,
			"Create and use a non-root user with USER instruction")
		scan.BestPracticeIssues = append(scan.BestPracticeIssues,
			"Container runs as root user")
	}

	if scan.UsesLatestTags {
		scan.SecurityRecommendations = append(scan.SecurityRecommendations,
			"Pin base images to specific versions instead of using 'latest' tag")
		scan.BestPracticeIssues = append(scan.BestPracticeIssues,
			"Uses 'latest' tag for base images")
	}

	if !scan.HasHealthCheck {
		scan.SecurityRecommendations = append(scan.SecurityRecommendations,
			"Add HEALTHCHECK instruction for container health monitoring")
		scan.BestPracticeIssues = append(scan.BestPracticeIssues,
			"Missing health check configuration")
	}

	if !scan.MultiStageOptimized && len(analysis.Stages) == 1 {
		scan.SecurityRecommendations = append(scan.SecurityRecommendations,
			"Consider using multi-stage builds to reduce final image size")
	}

	if !scan.CachingOptimized {
		scan.SecurityRecommendations = append(scan.SecurityRecommendations,
			"Optimize layer caching by copying package files before source code")
		scan.BestPracticeIssues = append(scan.BestPracticeIssues,
			"Layer caching not optimized")
	}
}

// DiscoverDockerFiles finds all Docker-related files in a directory
func DiscoverDockerFiles(rootPath string) ([]string, []string, error) {
	var dockerfiles []string
	var composeFiles []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := strings.ToLower(info.Name())
		
		// Look for Dockerfiles
		if fileName == "dockerfile" || strings.HasSuffix(fileName, ".dockerfile") {
			dockerfiles = append(dockerfiles, path)
		}
		
		// Look for docker-compose files
		if matched, _ := regexp.MatchString(`^docker-compose.*\.ya?ml$`, fileName); matched {
			composeFiles = append(composeFiles, path)
		}

		return nil
	})

	return dockerfiles, composeFiles, err
}