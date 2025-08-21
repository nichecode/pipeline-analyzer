package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Writer handles generating markdown documentation for Docker analysis
type Writer struct {
	outputDir string
}

// NewWriter creates a new Docker writer
func NewWriter(outputDir string) *Writer {
	return &Writer{
		outputDir: outputDir,
	}
}

// WriteAllFiles generates all Docker analysis documentation
func (w *Writer) WriteAllFiles(analysis *DockerAnalysis, configPath string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write main README
	if err := w.writeMainReadme(analysis, configPath); err != nil {
		return fmt.Errorf("failed to write main README: %w", err)
	}

	// Write individual Dockerfile analyses
	dockerfileDir := filepath.Join(w.outputDir, "dockerfiles")
	if len(analysis.Dockerfiles) > 0 {
		if err := os.MkdirAll(dockerfileDir, 0755); err != nil {
			return fmt.Errorf("failed to create dockerfiles directory: %w", err)
		}

		for i, dockerfile := range analysis.Dockerfiles {
			filename := fmt.Sprintf("dockerfile-%d.md", i+1)
			if err := w.writeDockerfileAnalysis(dockerfile, filepath.Join(dockerfileDir, filename)); err != nil {
				return fmt.Errorf("failed to write Dockerfile analysis: %w", err)
			}
		}
	}

	// Write docker-compose analysis if present
	if analysis.DockerCompose != nil {
		if err := w.writeDockerComposeAnalysis(analysis.DockerCompose, filepath.Join(w.outputDir, "docker-compose.md")); err != nil {
			return fmt.Errorf("failed to write docker-compose analysis: %w", err)
		}
	}

	// Write security analysis
	if err := w.writeSecurityAnalysis(analysis, filepath.Join(w.outputDir, "security-analysis.md")); err != nil {
		return fmt.Errorf("failed to write security analysis: %w", err)
	}

	// Write optimization guide
	if err := w.writeOptimizationGuide(analysis, filepath.Join(w.outputDir, "optimization-guide.md")); err != nil {
		return fmt.Errorf("failed to write optimization guide: %w", err)
	}

	return nil
}

// writeMainReadme writes the main Docker analysis README
func (w *Writer) writeMainReadme(analysis *DockerAnalysis, configPath string) error {
	content := fmt.Sprintf(`# Docker Analysis Report

**Generated:** %s
**Config Path:** %s

## 📊 Overview

- **Total Dockerfiles:** %d
- **Multi-stage builds:** %d
- **Docker Compose:** %t
- **Security issues:** %d
- **Optimization opportunities:** %d
- **Overall score:** %d/100

## 🐳 Summary

%s

## 📁 Files Analyzed

`,
		analysis.GeneratedAt.Format(time.RFC3339),
		configPath,
		analysis.Summary.TotalDockerfiles,
		analysis.Summary.MultiStageBuilds,
		analysis.Summary.HasDockerCompose,
		analysis.Summary.SecurityIssues,
		analysis.Summary.OptimizationIssues,
		analysis.Summary.OverallScore,
		generateSummaryText(analysis.Summary),
	)

	// List Dockerfiles
	if len(analysis.Dockerfiles) > 0 {
		content += "### 🐳 Dockerfiles\n\n"
		for i, dockerfile := range analysis.Dockerfiles {
			filename := filepath.Base(dockerfile.FilePath)
			stages := "single-stage"
			if dockerfile.MultiStage {
				stages = fmt.Sprintf("%d stages", len(dockerfile.Stages))
			}
			
			content += fmt.Sprintf("- [%s](dockerfiles/dockerfile-%d.md) - %s build\n", filename, i+1, stages)
		}
		content += "\n"
	}

	// Docker Compose
	if analysis.DockerCompose != nil {
		content += "### 🔧 Docker Compose\n\n"
		content += fmt.Sprintf("- [docker-compose configuration](docker-compose.md) - %d services\n\n", analysis.DockerCompose.ServiceCount)
	}

	// Analysis sections
	content += `## 📋 Analysis Sections

- [🔒 Security Analysis](security-analysis.md) - Security issues and recommendations
- [⚡ Optimization Guide](optimization-guide.md) - Performance and caching optimizations

## 🔍 Key Insights

`

	// Base images
	if len(analysis.Summary.UniqueBaseImages) > 0 {
		content += "### Base Images Used\n\n"
		for _, image := range analysis.Summary.UniqueBaseImages {
			content += fmt.Sprintf("- `%s`\n", image)
		}
		content += "\n"
	}

	// Most common instructions
	if len(analysis.Summary.MostCommonInstructions) > 0 {
		content += "### Most Common Instructions\n\n"
		for i, instruction := range analysis.Summary.MostCommonInstructions {
			content += fmt.Sprintf("%d. %s\n", i+1, instruction)
		}
		content += "\n"
	}

	// Recommendations
	if len(analysis.Summary.Recommendations) > 0 {
		content += "## 🎯 Recommendations\n\n"
		for _, recommendation := range analysis.Summary.Recommendations {
			content += fmt.Sprintf("- %s\n", recommendation)
		}
		content += "\n"
	}

	content += `## 🚀 Getting Started

1. **Review the overview above** to understand your Docker setup
2. **Check individual Dockerfile analyses** for specific optimizations
3. **Follow security recommendations** to improve container security
4. **Implement optimization suggestions** for better performance

## Navigation

- [← Back to Discovery Overview](../README.md)
`

	return os.WriteFile(filepath.Join(w.outputDir, "README.md"), []byte(content), 0644)
}

// writeDockerfileAnalysis writes analysis for a single Dockerfile
func (w *Writer) writeDockerfileAnalysis(dockerfile *DockerfileAnalysis, outputPath string) error {
	filename := filepath.Base(dockerfile.FilePath)
	
	content := fmt.Sprintf(`# Dockerfile Analysis: %s

**File:** %s
**Multi-stage:** %t
**Stages:** %d

## 📊 Overview

`,
		filename,
		dockerfile.FilePath,
		dockerfile.MultiStage,
		len(dockerfile.Stages),
	)

	// Base images
	if len(dockerfile.BaseImages) > 0 {
		content += "### Base Images\n\n"
		for _, image := range dockerfile.BaseImages {
			content += fmt.Sprintf("- `%s`\n", image)
		}
		content += "\n"
	}

	// Stages
	content += "## 🏗️ Build Stages\n\n"
	for i, stage := range dockerfile.Stages {
		stageTitle := fmt.Sprintf("Stage %d", i+1)
		if stage.Name != "" {
			stageTitle += fmt.Sprintf(" (%s)", stage.Name)
		}
		
		content += fmt.Sprintf("### %s\n\n", stageTitle)
		content += fmt.Sprintf("- **Base image:** `%s`\n", stage.BaseImage)
		content += fmt.Sprintf("- **Instructions:** %d\n\n", len(stage.Instructions))

		// List key instructions
		content += "**Key instructions:**\n\n"
		for _, instruction := range stage.Instructions {
			if instruction.Instruction != "FROM" {
				args := strings.Join(instruction.Arguments, " ")
				if len(args) > 80 {
					args = args[:80] + "..."
				}
				content += fmt.Sprintf("- `%s %s`\n", instruction.Instruction, args)
			}
		}
		content += "\n"
	}

	// Configuration
	content += "## ⚙️ Configuration\n\n"

	if len(dockerfile.ExposedPorts) > 0 {
		content += "### Exposed Ports\n\n"
		for _, port := range dockerfile.ExposedPorts {
			content += fmt.Sprintf("- `%s`\n", port)
		}
		content += "\n"
	}

	if len(dockerfile.Environment) > 0 {
		content += "### Environment Variables\n\n"
		for key, value := range dockerfile.Environment {
			content += fmt.Sprintf("- `%s=%s`\n", key, value)
		}
		content += "\n"
	}

	if dockerfile.WorkingDir != "" {
		content += fmt.Sprintf("### Working Directory\n\n`%s`\n\n", dockerfile.WorkingDir)
	}

	if dockerfile.User != "" {
		content += fmt.Sprintf("### User\n\n`%s`\n\n", dockerfile.User)
	}

	if len(dockerfile.Volumes) > 0 {
		content += "### Volumes\n\n"
		for _, volume := range dockerfile.Volumes {
			content += fmt.Sprintf("- `%s`\n", volume)
		}
		content += "\n"
	}

	// Health check
	if dockerfile.HealthCheck != nil {
		content += "### Health Check\n\n"
		content += fmt.Sprintf("- **Command:** `%s`\n", strings.Join(dockerfile.HealthCheck.Command, " "))
		if dockerfile.HealthCheck.Interval > 0 {
			content += fmt.Sprintf("- **Interval:** %s\n", dockerfile.HealthCheck.Interval)
		}
		if dockerfile.HealthCheck.Timeout > 0 {
			content += fmt.Sprintf("- **Timeout:** %s\n", dockerfile.HealthCheck.Timeout)
		}
		if dockerfile.HealthCheck.Retries > 0 {
			content += fmt.Sprintf("- **Retries:** %d\n", dockerfile.HealthCheck.Retries)
		}
		content += "\n"
	}

	// Security analysis
	if dockerfile.SecurityScan != nil {
		content += "## 🔒 Security Analysis\n\n"
		
		if dockerfile.SecurityScan.RunAsRoot {
			content += "⚠️ **Warning:** Container runs as root user\n\n"
		} else {
			content += "✅ **Good:** Container uses non-root user\n\n"
		}

		if dockerfile.SecurityScan.UsesLatestTags {
			content += "⚠️ **Warning:** Uses 'latest' tag for base images\n\n"
		} else {
			content += "✅ **Good:** Base images use specific versions\n\n"
		}

		if dockerfile.SecurityScan.HasHealthCheck {
			content += "✅ **Good:** Health check is configured\n\n"
		} else {
			content += "⚠️ **Warning:** No health check configured\n\n"
		}

		// Security recommendations
		if len(dockerfile.SecurityScan.SecurityRecommendations) > 0 {
			content += "### Security Recommendations\n\n"
			for _, rec := range dockerfile.SecurityScan.SecurityRecommendations {
				content += fmt.Sprintf("- %s\n", rec)
			}
			content += "\n"
		}

		// Best practice issues
		if len(dockerfile.SecurityScan.BestPracticeIssues) > 0 {
			content += "### Best Practice Issues\n\n"
			for _, issue := range dockerfile.SecurityScan.BestPracticeIssues {
				content += fmt.Sprintf("- %s\n", issue)
			}
			content += "\n"
		}
	}

	content += `## Navigation

- [← Back to Docker Overview](../README.md)
- [Security Analysis](../security-analysis.md)
- [Optimization Guide](../optimization-guide.md)
`

	return os.WriteFile(outputPath, []byte(content), 0644)
}

// writeDockerComposeAnalysis writes analysis for docker-compose.yml
func (w *Writer) writeDockerComposeAnalysis(compose *DockerComposeAnalysis, outputPath string) error {
	content := fmt.Sprintf(`# Docker Compose Analysis

**File:** %s
**Version:** %s
**Services:** %d

## 📊 Overview

`,
		compose.FilePath,
		compose.Version,
		compose.ServiceCount,
	)

	// Services
	content += "## 🚀 Services\n\n"
	content += "| Service | Image | Ports | Dependencies |\n"
	content += "|---------|-------|-------|-------------|\n"

	for serviceName, service := range compose.Services {
		ports := strings.Join(service.Ports, ", ")
		if ports == "" {
			ports = "-"
		}
		deps := strings.Join(service.DependsOn, ", ")
		if deps == "" {
			deps = "-"
		}
		
		image := service.Image
		if image == "" && service.Build != nil {
			image = fmt.Sprintf("Built from %s", service.Build.Context)
		}
		if image == "" {
			image = "-"
		}

		content += fmt.Sprintf("| %s | %s | %s | %s |\n", serviceName, image, ports, deps)
	}
	content += "\n"

	// Service categories
	if compose.Analysis != nil {
		if len(compose.Analysis.DatabaseServices) > 0 {
			content += "### Database Services\n\n"
			for _, service := range compose.Analysis.DatabaseServices {
				content += fmt.Sprintf("- %s\n", service)
			}
			content += "\n"
		}

		if len(compose.Analysis.WebServices) > 0 {
			content += "### Web Services\n\n"
			for _, service := range compose.Analysis.WebServices {
				content += fmt.Sprintf("- %s\n", service)
			}
			content += "\n"
		}

		if len(compose.Analysis.CacheServices) > 0 {
			content += "### Cache Services\n\n"
			for _, service := range compose.Analysis.CacheServices {
				content += fmt.Sprintf("- %s\n", service)
			}
			content += "\n"
		}
	}

	// Networks and Volumes
	if len(compose.Networks) > 0 {
		content += "### Networks\n\n"
		for networkName := range compose.Networks {
			content += fmt.Sprintf("- %s\n", networkName)
		}
		content += "\n"
	}

	if len(compose.Volumes) > 0 {
		content += "### Volumes\n\n"
		for volumeName := range compose.Volumes {
			content += fmt.Sprintf("- %s\n", volumeName)
		}
		content += "\n"
	}

	// Analysis results
	if compose.Analysis != nil {
		content += "## 🔍 Analysis Results\n\n"
		
		content += fmt.Sprintf("- **Complexity Score:** %d/100\n", compose.Analysis.ComplexityScore)
		content += fmt.Sprintf("- **Security Issues:** %d\n", len(compose.Analysis.SecurityIssues))
		content += fmt.Sprintf("- **Performance Issues:** %d\n", len(compose.Analysis.PerformanceIssues))
		content += fmt.Sprintf("- **Port Conflicts:** %d\n\n", len(compose.Analysis.PortConflicts))

		// Issues
		if len(compose.Analysis.SecurityIssues) > 0 {
			content += "### Security Issues\n\n"
			for _, issue := range compose.Analysis.SecurityIssues {
				content += fmt.Sprintf("- ⚠️ %s\n", issue)
			}
			content += "\n"
		}

		if len(compose.Analysis.PerformanceIssues) > 0 {
			content += "### Performance Issues\n\n"
			for _, issue := range compose.Analysis.PerformanceIssues {
				content += fmt.Sprintf("- ⚠️ %s\n", issue)
			}
			content += "\n"
		}

		if len(compose.Analysis.PortConflicts) > 0 {
			content += "### Port Conflicts\n\n"
			for _, conflict := range compose.Analysis.PortConflicts {
				content += fmt.Sprintf("- ⚠️ %s\n", conflict)
			}
			content += "\n"
		}

		// Recommendations
		if len(compose.Analysis.Recommendations) > 0 {
			content += "### Recommendations\n\n"
			for _, rec := range compose.Analysis.Recommendations {
				content += fmt.Sprintf("- %s\n", rec)
			}
			content += "\n"
		}
	}

	content += `## Navigation

- [← Back to Docker Overview](../README.md)
- [Security Analysis](security-analysis.md)
- [Optimization Guide](optimization-guide.md)
`

	return os.WriteFile(outputPath, []byte(content), 0644)
}

// writeSecurityAnalysis writes the security analysis report
func (w *Writer) writeSecurityAnalysis(analysis *DockerAnalysis, outputPath string) error {
	content := `# Docker Security Analysis

## 🔒 Security Overview

This report analyzes security issues and provides recommendations for your Docker configuration.

`

	totalIssues := 0
	recommendations := []string{}

	// Analyze each Dockerfile
	if len(analysis.Dockerfiles) > 0 {
		content += "## 📋 Dockerfile Security Issues\n\n"
		
		for i, dockerfile := range analysis.Dockerfiles {
			filename := filepath.Base(dockerfile.FilePath)
			content += fmt.Sprintf("### %s\n\n", filename)
			
			if dockerfile.SecurityScan != nil {
				issues := []string{}
				
				if dockerfile.SecurityScan.RunAsRoot {
					issues = append(issues, "Runs as root user")
					totalIssues++
				}
				
				if dockerfile.SecurityScan.UsesLatestTags {
					issues = append(issues, "Uses 'latest' tag for base images")
					totalIssues++
				}
				
				if !dockerfile.SecurityScan.HasHealthCheck {
					issues = append(issues, "Missing health check")
					totalIssues++
				}

				if len(issues) > 0 {
					content += "**Issues found:**\n\n"
					for _, issue := range issues {
						content += fmt.Sprintf("- ⚠️ %s\n", issue)
					}
				} else {
					content += "✅ No major security issues found\n"
				}

				// Add security recommendations
				recommendations = append(recommendations, dockerfile.SecurityScan.SecurityRecommendations...)
				
				content += "\n"
				content += fmt.Sprintf("[View detailed analysis](dockerfiles/dockerfile-%d.md)\n\n", i+1)
			}
		}
	}

	// Docker Compose security
	if analysis.DockerCompose != nil && analysis.DockerCompose.Analysis != nil {
		if len(analysis.DockerCompose.Analysis.SecurityIssues) > 0 {
			content += "## 🔧 Docker Compose Security Issues\n\n"
			for _, issue := range analysis.DockerCompose.Analysis.SecurityIssues {
				content += fmt.Sprintf("- ⚠️ %s\n", issue)
			}
			content += "\n"
			totalIssues += len(analysis.DockerCompose.Analysis.SecurityIssues)
		}
	}

	// Summary
	content += fmt.Sprintf("## 📊 Security Summary\n\n")
	content += fmt.Sprintf("- **Total security issues:** %d\n", totalIssues)
	content += fmt.Sprintf("- **Dockerfiles analyzed:** %d\n", len(analysis.Dockerfiles))
	if analysis.DockerCompose != nil {
		content += fmt.Sprintf("- **Services analyzed:** %d\n", analysis.DockerCompose.ServiceCount)
	}
	content += "\n"

	// Unique recommendations
	uniqueRecs := make(map[string]bool)
	var finalRecs []string
	for _, rec := range recommendations {
		if !uniqueRecs[rec] {
			uniqueRecs[rec] = true
			finalRecs = append(finalRecs, rec)
		}
	}

	if len(finalRecs) > 0 {
		content += "## 🎯 Security Recommendations\n\n"
		for _, rec := range finalRecs {
			content += fmt.Sprintf("- %s\n", rec)
		}
		content += "\n"
	}

	content += `## 🛡️ Best Practices

### General Security
- Use official base images from trusted sources
- Pin base images to specific versions
- Run containers as non-root users
- Implement health checks for monitoring
- Regularly update base images for security patches

### Dockerfile Security
- Use multi-stage builds to reduce attack surface
- Don't store secrets in environment variables
- Use COPY instead of ADD when possible
- Minimize the number of RUN layers
- Use .dockerignore to exclude sensitive files

### Docker Compose Security
- Don't expose database ports directly
- Use secrets management for sensitive data
- Implement proper network segmentation
- Set resource limits for containers
- Use restart policies appropriately

## Navigation

- [← Back to Docker Overview](README.md)
- [Optimization Guide](optimization-guide.md)
`

	return os.WriteFile(outputPath, []byte(content), 0644)
}

// writeOptimizationGuide writes the optimization guide
func (w *Writer) writeOptimizationGuide(analysis *DockerAnalysis, outputPath string) error {
	content := `# Docker Optimization Guide

## ⚡ Performance Optimization

This guide provides recommendations to optimize your Docker setup for better performance, smaller images, and faster builds.

`

	// Dockerfile optimizations
	if len(analysis.Dockerfiles) > 0 {
		content += "## 🐳 Dockerfile Optimizations\n\n"
		
		for i, dockerfile := range analysis.Dockerfiles {
			filename := filepath.Base(dockerfile.FilePath)
			content += fmt.Sprintf("### %s\n\n", filename)
			
			optimizations := []string{}
			
			if !dockerfile.MultiStage {
				optimizations = append(optimizations, "Consider using multi-stage builds to reduce final image size")
			}
			
			if dockerfile.SecurityScan != nil && !dockerfile.SecurityScan.CachingOptimized {
				optimizations = append(optimizations, "Optimize layer caching by copying package files before source code")
			}
			
			if len(optimizations) > 0 {
				content += "**Optimization opportunities:**\n\n"
				for _, opt := range optimizations {
					content += fmt.Sprintf("- %s\n", opt)
				}
			} else {
				content += "✅ No major optimization issues found\n"
			}
			
			content += "\n"
			content += fmt.Sprintf("[View detailed analysis](dockerfiles/dockerfile-%d.md)\n\n", i+1)
		}
	}

	// Docker Compose optimizations
	if analysis.DockerCompose != nil && analysis.DockerCompose.Analysis != nil {
		if len(analysis.DockerCompose.Analysis.PerformanceIssues) > 0 {
			content += "## 🔧 Docker Compose Optimizations\n\n"
			for _, issue := range analysis.DockerCompose.Analysis.PerformanceIssues {
				content += fmt.Sprintf("- %s\n", issue)
			}
			content += "\n"
		}
	}

	content += `## 🚀 Best Practices for Optimization

### Image Size Optimization
- Use slim or alpine base images when possible
- Use multi-stage builds to exclude build dependencies
- Remove package caches after installation
- Use .dockerignore to exclude unnecessary files
- Minimize the number of layers

### Build Performance
- Order instructions from least to most frequently changing
- Use build cache effectively
- Copy package files before source code
- Use specific COPY instructions instead of copying everything

### Runtime Performance
- Set appropriate resource limits
- Use health checks for proper load balancing
- Implement proper logging strategies
- Use read-only filesystems when possible
- Configure proper restart policies

### Multi-stage Build Example

` + "```dockerfile" + `
# Build stage
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force
COPY src/ ./src/
RUN npm run build

# Production stage  
FROM nginx:alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
` + "```" + `

### Layer Caching Optimization

` + "```dockerfile" + `
# Good: Copy package files first for better caching
COPY package*.json ./
RUN npm install

# Then copy source code
COPY src/ ./src/
RUN npm run build
` + "```" + `

## Navigation

- [← Back to Docker Overview](README.md)
- [Security Analysis](security-analysis.md)
`

	return os.WriteFile(outputPath, []byte(content), 0644)
}

// generateSummaryText creates a summary description
func generateSummaryText(summary *DockerSummary) string {
	if summary.TotalDockerfiles == 0 {
		return "No Dockerfiles found in the project."
	}

	text := fmt.Sprintf("Found %d Dockerfile(s) in your project. ", summary.TotalDockerfiles)

	if summary.MultiStageBuilds > 0 {
		text += fmt.Sprintf("%d use multi-stage builds for optimized image sizes. ", summary.MultiStageBuilds)
	}

	if summary.HasDockerCompose {
		text += fmt.Sprintf("Docker Compose configuration manages %d services. ", summary.ServiceCount)
	}

	if summary.SecurityIssues > 0 {
		text += fmt.Sprintf("Found %d security issues that should be addressed. ", summary.SecurityIssues)
	}

	if summary.OptimizationIssues > 0 {
		text += fmt.Sprintf("Identified %d optimization opportunities. ", summary.OptimizationIssues)
	}

	if summary.OverallScore >= 80 {
		text += "Overall configuration follows Docker best practices well."
	} else if summary.OverallScore >= 60 {
		text += "Configuration is good but has room for improvement."
	} else {
		text += "Configuration needs significant improvements for security and optimization."
	}

	return text
}