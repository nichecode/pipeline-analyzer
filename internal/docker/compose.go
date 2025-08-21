package docker

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ParseDockerCompose parses a docker-compose.yml file and returns its analysis
func ParseDockerCompose(composePath string) (*DockerComposeAnalysis, error) {
	content, err := os.ReadFile(composePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read docker-compose file: %w", err)
	}

	var compose map[string]interface{}
	if err := yaml.Unmarshal(content, &compose); err != nil {
		return nil, fmt.Errorf("failed to parse docker-compose YAML: %w", err)
	}

	analysis := &DockerComposeAnalysis{
		FilePath: composePath,
		Services: make(map[string]*DockerComposeService),
		Networks: make(map[string]interface{}),
		Volumes:  make(map[string]interface{}),
		Secrets:  make(map[string]interface{}),
	}

	// Extract version
	if version, ok := compose["version"]; ok {
		if versionStr, ok := version.(string); ok {
			analysis.Version = versionStr
		}
	}

	// Parse services
	if services, ok := compose["services"]; ok {
		if servicesMap, ok := services.(map[string]interface{}); ok {
			for serviceName, serviceConfig := range servicesMap {
				service, err := parseService(serviceName, serviceConfig)
				if err != nil {
					continue // Skip invalid services
				}
				analysis.Services[serviceName] = service
			}
		}
	}

	// Parse networks
	if networks, ok := compose["networks"]; ok {
		if networksMap, ok := networks.(map[string]interface{}); ok {
			analysis.Networks = networksMap
		}
	}

	// Parse volumes
	if volumes, ok := compose["volumes"]; ok {
		if volumesMap, ok := volumes.(map[string]interface{}); ok {
			analysis.Volumes = volumesMap
		}
	}

	// Parse secrets
	if secrets, ok := compose["secrets"]; ok {
		if secretsMap, ok := secrets.(map[string]interface{}); ok {
			analysis.Secrets = secretsMap
		}
	}

	analysis.ServiceCount = len(analysis.Services)

	// Perform analysis
	analysis.Analysis = analyzeDockerCompose(analysis)

	return analysis, nil
}

// parseService parses a single service configuration
func parseService(name string, config interface{}) (*DockerComposeService, error) {
	serviceMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid service configuration for %s", name)
	}

	service := &DockerComposeService{
		Name:        name,
		Environment: make(map[string]string),
	}

	// Parse image
	if image, ok := serviceMap["image"]; ok {
		if imageStr, ok := image.(string); ok {
			service.Image = imageStr
		}
	}

	// Parse build configuration
	if build, ok := serviceMap["build"]; ok {
		service.Build = parseBuildConfig(build)
	}

	// Parse ports
	if ports, ok := serviceMap["ports"]; ok {
		service.Ports = parseStringArray(ports)
	}

	// Parse environment
	if env, ok := serviceMap["environment"]; ok {
		service.Environment = parseEnvironment(env)
	}

	// Parse volumes
	if volumes, ok := serviceMap["volumes"]; ok {
		service.Volumes = parseStringArray(volumes)
	}

	// Parse depends_on
	if dependsOn, ok := serviceMap["depends_on"]; ok {
		service.DependsOn = parseStringArray(dependsOn)
	}

	// Parse networks
	if networks, ok := serviceMap["networks"]; ok {
		service.Networks = parseStringArray(networks)
	}

	// Parse command
	if command, ok := serviceMap["command"]; ok {
		service.Command = parseStringArray(command)
	}

	// Parse healthcheck
	if healthcheck, ok := serviceMap["healthcheck"]; ok {
		service.HealthCheck = parseHealthCheck(healthcheck)
	}

	// Parse restart policy
	if restart, ok := serviceMap["restart"]; ok {
		if restartStr, ok := restart.(string); ok {
			service.RestartPolicy = restartStr
		}
	}

	// Parse resource limits
	if deploy, ok := serviceMap["deploy"]; ok {
		if deployMap, ok := deploy.(map[string]interface{}); ok {
			if resources, ok := deployMap["resources"]; ok {
				service.Resources = parseResourceLimits(resources)
			}
		}
	}

	return service, nil
}

// parseBuildConfig parses build configuration
func parseBuildConfig(build interface{}) *BuildConfig {
	buildConfig := &BuildConfig{
		Args: make(map[string]string),
	}

	switch b := build.(type) {
	case string:
		buildConfig.Context = b
	case map[string]interface{}:
		if context, ok := b["context"]; ok {
			if contextStr, ok := context.(string); ok {
				buildConfig.Context = contextStr
			}
		}
		if dockerfile, ok := b["dockerfile"]; ok {
			if dockerfileStr, ok := dockerfile.(string); ok {
				buildConfig.Dockerfile = dockerfileStr
			}
		}
		if args, ok := b["args"]; ok {
			buildConfig.Args = parseEnvironment(args)
		}
		if target, ok := b["target"]; ok {
			if targetStr, ok := target.(string); ok {
				buildConfig.Target = targetStr
			}
		}
	}

	return buildConfig
}

// parseStringArray converts various YAML formats to string array
func parseStringArray(value interface{}) []string {
	switch v := value.(type) {
	case []interface{}:
		var result []string
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	case string:
		return []string{v}
	}
	return []string{}
}

// parseEnvironment parses environment variables from various formats
func parseEnvironment(env interface{}) map[string]string {
	envMap := make(map[string]string)

	switch e := env.(type) {
	case map[string]interface{}:
		for key, value := range e {
			if str, ok := value.(string); ok {
				envMap[key] = str
			} else {
				envMap[key] = fmt.Sprintf("%v", value)
			}
		}
	case []interface{}:
		for _, item := range e {
			if str, ok := item.(string); ok {
				if strings.Contains(str, "=") {
					parts := strings.SplitN(str, "=", 2)
					if len(parts) == 2 {
						envMap[parts[0]] = parts[1]
					}
				} else {
					envMap[str] = ""
				}
			}
		}
	}

	return envMap
}

// parseHealthCheck parses health check configuration
func parseHealthCheck(healthcheck interface{}) *HealthCheck {
	hcMap, ok := healthcheck.(map[string]interface{})
	if !ok {
		return nil
	}

	hc := &HealthCheck{}

	if test, ok := hcMap["test"]; ok {
		hc.Command = parseStringArray(test)
	}

	if interval, ok := hcMap["interval"]; ok {
		if intervalStr, ok := interval.(string); ok {
			if d, err := parseDuration(intervalStr); err == nil {
				hc.Interval = d
			}
		}
	}

	if timeout, ok := hcMap["timeout"]; ok {
		if timeoutStr, ok := timeout.(string); ok {
			if d, err := parseDuration(timeoutStr); err == nil {
				hc.Timeout = d
			}
		}
	}

	if startPeriod, ok := hcMap["start_period"]; ok {
		if startPeriodStr, ok := startPeriod.(string); ok {
			if d, err := parseDuration(startPeriodStr); err == nil {
				hc.StartPeriod = d
			}
		}
	}

	if retries, ok := hcMap["retries"]; ok {
		if retriesNum, ok := retries.(int); ok {
			hc.Retries = retriesNum
		} else if retriesStr, ok := retries.(string); ok {
			if retriesNum, err := strconv.Atoi(retriesStr); err == nil {
				hc.Retries = retriesNum
			}
		}
	}

	return hc
}

// parseResourceLimits parses resource limits configuration
func parseResourceLimits(resources interface{}) *ResourceLimits {
	resourcesMap, ok := resources.(map[string]interface{})
	if !ok {
		return nil
	}

	limits := &ResourceLimits{}

	if limitsSection, ok := resourcesMap["limits"]; ok {
		if limitsMap, ok := limitsSection.(map[string]interface{}); ok {
			if cpus, ok := limitsMap["cpus"]; ok {
				if cpuStr, ok := cpus.(string); ok {
					limits.CPULimit = cpuStr
				}
			}
			if memory, ok := limitsMap["memory"]; ok {
				if memStr, ok := memory.(string); ok {
					limits.MemoryLimit = memStr
				}
			}
		}
	}

	if reservationsSection, ok := resourcesMap["reservations"]; ok {
		if reservationsMap, ok := reservationsSection.(map[string]interface{}); ok {
			if cpus, ok := reservationsMap["cpus"]; ok {
				if cpuStr, ok := cpus.(string); ok {
					limits.CPUReserve = cpuStr
				}
			}
			if memory, ok := reservationsMap["memory"]; ok {
				if memStr, ok := memory.(string); ok {
					limits.MemoryReserve = memStr
				}
			}
		}
	}

	return limits
}

// parseDuration parses Docker compose duration format (like 30s, 1m30s)
func parseDuration(durationStr string) (time.Duration, error) {
	// Docker compose duration format can be like "30s", "1m", "1h30m"
	return time.ParseDuration(durationStr)
}

// analyzeDockerCompose performs analysis on the docker-compose configuration
func analyzeDockerCompose(compose *DockerComposeAnalysis) *ComposeAnalysisResults {
	analysis := &ComposeAnalysisResults{
		ServiceDependencies: make(map[string][]string),
		PortConflicts:      []string{},
		NetworkingIssues:   []string{},
		SecurityIssues:     []string{},
		PerformanceIssues:  []string{},
		Recommendations:    []string{},
	}

	// Analyze services
	for serviceName, service := range compose.Services {
		// Categorize services
		categorizeService(service, analysis)
		
		// Check dependencies
		if len(service.DependsOn) > 0 {
			analysis.ServiceDependencies[serviceName] = service.DependsOn
		}
		
		// Check for security issues
		checkServiceSecurity(service, analysis)
		
		// Check for performance issues
		checkServicePerformance(service, analysis)
	}

	// Check for port conflicts
	checkPortConflicts(compose.Services, analysis)

	// Generate recommendations
	generateComposeRecommendations(compose, analysis)

	// Calculate complexity score
	analysis.ComplexityScore = calculateComplexityScore(compose, analysis)

	return analysis
}

// categorizeService categorizes services by type
func categorizeService(service *DockerComposeService, analysis *ComposeAnalysisResults) {
	serviceName := strings.ToLower(service.Name)
	image := strings.ToLower(service.Image)

	// Database services
	if strings.Contains(serviceName, "db") || 
	   strings.Contains(serviceName, "database") ||
	   strings.Contains(image, "postgres") ||
	   strings.Contains(image, "mysql") ||
	   strings.Contains(image, "mongo") ||
	   strings.Contains(image, "redis") {
		analysis.DatabaseServices = append(analysis.DatabaseServices, service.Name)
	}

	// Cache services
	if strings.Contains(serviceName, "cache") ||
	   strings.Contains(serviceName, "redis") ||
	   strings.Contains(image, "redis") ||
	   strings.Contains(image, "memcached") {
		analysis.CacheServices = append(analysis.CacheServices, service.Name)
	}

	// Web services
	if strings.Contains(serviceName, "web") ||
	   strings.Contains(serviceName, "app") ||
	   strings.Contains(serviceName, "api") ||
	   strings.Contains(serviceName, "frontend") ||
	   strings.Contains(serviceName, "backend") ||
	   len(service.Ports) > 0 {
		analysis.WebServices = append(analysis.WebServices, service.Name)
	}
}

// checkServiceSecurity checks for security issues in services
func checkServiceSecurity(service *DockerComposeService, analysis *ComposeAnalysisResults) {
	// Check for privileged mode
	// Note: This would need to check for privileged: true in the service config

	// Check for exposed sensitive ports
	for _, port := range service.Ports {
		if strings.Contains(port, "22:22") { // SSH
			analysis.SecurityIssues = append(analysis.SecurityIssues,
				fmt.Sprintf("Service %s exposes SSH port (22)", service.Name))
		}
		if strings.Contains(port, "3306:3306") { // MySQL
			analysis.SecurityIssues = append(analysis.SecurityIssues,
				fmt.Sprintf("Service %s exposes MySQL port directly", service.Name))
		}
		if strings.Contains(port, "5432:5432") { // PostgreSQL
			analysis.SecurityIssues = append(analysis.SecurityIssues,
				fmt.Sprintf("Service %s exposes PostgreSQL port directly", service.Name))
		}
	}

	// Check for missing health checks
	if service.HealthCheck == nil && len(service.Ports) > 0 {
		analysis.PerformanceIssues = append(analysis.PerformanceIssues,
			fmt.Sprintf("Service %s is missing health check", service.Name))
	}
}

// checkServicePerformance checks for performance issues in services
func checkServicePerformance(service *DockerComposeService, analysis *ComposeAnalysisResults) {
	// Check for missing resource limits
	if service.Resources == nil {
		analysis.PerformanceIssues = append(analysis.PerformanceIssues,
			fmt.Sprintf("Service %s has no resource limits", service.Name))
	}

	// Check restart policy
	if service.RestartPolicy == "" {
		analysis.PerformanceIssues = append(analysis.PerformanceIssues,
			fmt.Sprintf("Service %s has no restart policy", service.Name))
	}
}

// checkPortConflicts checks for port conflicts between services
func checkPortConflicts(services map[string]*DockerComposeService, analysis *ComposeAnalysisResults) {
	portMap := make(map[string][]string)

	for serviceName, service := range services {
		for _, port := range service.Ports {
			// Extract host port
			parts := strings.Split(port, ":")
			if len(parts) >= 2 {
				hostPort := parts[0]
				portMap[hostPort] = append(portMap[hostPort], serviceName)
			}
		}
	}

	// Check for conflicts
	for port, serviceList := range portMap {
		if len(serviceList) > 1 {
			analysis.PortConflicts = append(analysis.PortConflicts,
				fmt.Sprintf("Port %s is used by services: %s", port, strings.Join(serviceList, ", ")))
		}
	}
}

// generateComposeRecommendations generates recommendations for the docker-compose setup
func generateComposeRecommendations(compose *DockerComposeAnalysis, analysis *ComposeAnalysisResults) {
	if len(analysis.SecurityIssues) > 0 {
		analysis.Recommendations = append(analysis.Recommendations,
			"Review and fix security issues identified in the analysis")
	}

	if len(analysis.PerformanceIssues) > 0 {
		analysis.Recommendations = append(analysis.Recommendations,
			"Add resource limits and health checks to improve performance monitoring")
	}

	if len(analysis.PortConflicts) > 0 {
		analysis.Recommendations = append(analysis.Recommendations,
			"Resolve port conflicts between services")
	}

	if len(compose.Networks) == 0 && len(compose.Services) > 1 {
		analysis.Recommendations = append(analysis.Recommendations,
			"Consider defining custom networks for service isolation")
	}

	if len(compose.Volumes) == 0 {
		// Check if any services use volumes
		hasVolumes := false
		for _, service := range compose.Services {
			if len(service.Volumes) > 0 {
				hasVolumes = true
				break
			}
		}
		if hasVolumes {
			analysis.Recommendations = append(analysis.Recommendations,
				"Consider using named volumes for better data management")
		}
	}
}

// calculateComplexityScore calculates a complexity score for the docker-compose setup
func calculateComplexityScore(compose *DockerComposeAnalysis, analysis *ComposeAnalysisResults) int {
	score := 0

	// Base score for number of services
	score += len(compose.Services) * 10

	// Additional score for networks
	score += len(compose.Networks) * 5

	// Additional score for volumes
	score += len(compose.Volumes) * 5

	// Additional score for service dependencies
	score += len(analysis.ServiceDependencies) * 3

	// Penalty for issues
	score += len(analysis.SecurityIssues) * 5
	score += len(analysis.PerformanceIssues) * 2
	score += len(analysis.PortConflicts) * 3

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}