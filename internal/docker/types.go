package docker

import "time"

// DockerfileInstruction represents a single instruction in a Dockerfile
type DockerfileInstruction struct {
	Instruction string            `json:"instruction"` // FROM, RUN, COPY, etc.
	Arguments   []string          `json:"arguments"`   // Arguments for the instruction
	Flags       map[string]string `json:"flags"`       // Flags like --from, --chown, etc.
	Line        int               `json:"line"`        // Line number in the Dockerfile
	Raw         string            `json:"raw"`         // Original raw instruction
}

// DockerfileStage represents a build stage in a multi-stage Dockerfile
type DockerfileStage struct {
	Name         string                   `json:"name"`         // Stage name (AS <name>)
	BaseImage    string                   `json:"base_image"`   // FROM image
	Instructions []*DockerfileInstruction `json:"instructions"` // All instructions in this stage
	Index        int                      `json:"index"`        // Stage order (0-based)
}

// DockerfileAnalysis contains the parsed and analyzed Dockerfile content
type DockerfileAnalysis struct {
	FilePath     string             `json:"file_path"`     // Path to the Dockerfile
	Stages       []*DockerfileStage `json:"stages"`        // Multi-stage build stages
	MultiStage   bool               `json:"multi_stage"`   // Is this a multi-stage build
	BaseImages   []string           `json:"base_images"`   // All base images used
	ExposedPorts []string           `json:"exposed_ports"` // EXPOSE instructions
	Labels       map[string]string  `json:"labels"`        // LABEL instructions
	Arguments    map[string]string  `json:"arguments"`     // ARG instructions
	Environment  map[string]string  `json:"environment"`   // ENV instructions
	WorkingDir   string             `json:"working_dir"`   // WORKDIR instruction
	User         string             `json:"user"`          // USER instruction
	Entrypoint   []string           `json:"entrypoint"`    // ENTRYPOINT instruction
	Command      []string           `json:"command"`       // CMD instruction
	HealthCheck  *HealthCheck       `json:"health_check"`  // HEALTHCHECK instruction
	Volumes      []string           `json:"volumes"`       // VOLUME instructions
	SecurityScan *SecurityScan      `json:"security_scan"` // Security analysis
}

// HealthCheck represents a Docker health check configuration
type HealthCheck struct {
	Command     []string      `json:"command"`
	Interval    time.Duration `json:"interval"`
	Timeout     time.Duration `json:"timeout"`
	StartPeriod time.Duration `json:"start_period"`
	Retries     int           `json:"retries"`
}

// SecurityScan contains security analysis results
type SecurityScan struct {
	RunAsRoot              bool     `json:"run_as_root"`
	HasSecurityUpdates     bool     `json:"has_security_updates"`
	UsesLatestTags         bool     `json:"uses_latest_tags"`
	HasHealthCheck         bool     `json:"has_health_check"`
	CachingOptimized       bool     `json:"caching_optimized"`
	MultiStageOptimized    bool     `json:"multi_stage_optimized"`
	VulnerableBaseImages   []string `json:"vulnerable_base_images"`
	SecurityRecommendations []string `json:"security_recommendations"`
	BestPracticeIssues     []string `json:"best_practice_issues"`
}

// DockerComposeService represents a service in docker-compose.yml
type DockerComposeService struct {
	Name            string            `json:"name"`
	Image           string            `json:"image"`
	Build           *BuildConfig      `json:"build"`
	Ports           []string          `json:"ports"`
	Environment     map[string]string `json:"environment"`
	Volumes         []string          `json:"volumes"`
	DependsOn       []string          `json:"depends_on"`
	Networks        []string          `json:"networks"`
	Command         []string          `json:"command"`
	HealthCheck     *HealthCheck      `json:"health_check"`
	RestartPolicy   string            `json:"restart_policy"`
	Resources       *ResourceLimits   `json:"resources"`
}

// BuildConfig represents the build configuration for a service
type BuildConfig struct {
	Context    string            `json:"context"`
	Dockerfile string            `json:"dockerfile"`
	Args       map[string]string `json:"args"`
	Target     string            `json:"target"`
}

// ResourceLimits represents resource constraints for a service
type ResourceLimits struct {
	CPULimit    string `json:"cpu_limit"`
	MemoryLimit string `json:"memory_limit"`
	CPUReserve  string `json:"cpu_reserve"`
	MemoryReserve string `json:"memory_reserve"`
}

// DockerComposeAnalysis contains the parsed and analyzed docker-compose content
type DockerComposeAnalysis struct {
	FilePath     string                           `json:"file_path"`
	Version      string                           `json:"version"`
	Services     map[string]*DockerComposeService `json:"services"`
	Networks     map[string]interface{}           `json:"networks"`
	Volumes      map[string]interface{}           `json:"volumes"`
	Secrets      map[string]interface{}           `json:"secrets"`
	ServiceCount int                              `json:"service_count"`
	Analysis     *ComposeAnalysisResults          `json:"analysis"`
}

// ComposeAnalysisResults contains analysis results for the docker-compose setup
type ComposeAnalysisResults struct {
	DatabaseServices   []string                  `json:"database_services"`
	WebServices        []string                  `json:"web_services"`
	CacheServices      []string                  `json:"cache_services"`
	ServiceDependencies map[string][]string       `json:"service_dependencies"`
	PortConflicts      []string                  `json:"port_conflicts"`
	NetworkingIssues   []string                  `json:"networking_issues"`
	SecurityIssues     []string                  `json:"security_issues"`
	PerformanceIssues  []string                  `json:"performance_issues"`
	Recommendations    []string                  `json:"recommendations"`
	ComplexityScore    int                       `json:"complexity_score"`
}

// DockerUsageReference represents where a Docker file is used
type DockerUsageReference struct {
	Tool        string `json:"tool"`        // "circleci", "github-actions", "gotask"
	File        string `json:"file"`        // Config file path
	Location    string `json:"location"`    // Specific location (job name, task name, etc.)
	Command     string `json:"command"`     // The actual command found
	Context     string `json:"context"`     // Additional context if available
}

// DockerUsageAnalysis tracks where Docker files are referenced
type DockerUsageAnalysis struct {
	DockerfileReferences    []*DockerUsageReference `json:"dockerfile_references"`
	DockerComposeReferences []*DockerUsageReference `json:"docker_compose_references"`
	DockerCommandReferences []*DockerUsageReference `json:"docker_command_references"`
	TotalReferences         int                     `json:"total_references"`
}

// DockerAnalysis aggregates all Docker-related analysis
type DockerAnalysis struct {
	Dockerfiles     []*DockerfileAnalysis      `json:"dockerfiles"`
	DockerCompose   []*DockerComposeAnalysis   `json:"docker_compose"`
	Usage           *DockerUsageAnalysis       `json:"usage"`
	Summary         *DockerSummary             `json:"summary"`
	GeneratedAt     time.Time                  `json:"generated_at"`
}

// DockerSummary provides high-level insights across all Docker configurations
type DockerSummary struct {
	TotalDockerfiles       int      `json:"total_dockerfiles"`
	MultiStageBuilds       int      `json:"multi_stage_builds"`
	SecurityIssues         int      `json:"security_issues"`
	OptimizationIssues     int      `json:"optimization_issues"`
	HasDockerCompose       bool     `json:"has_docker_compose"`
	TotalComposeFiles      int      `json:"total_compose_files"`
	ServiceCount           int      `json:"service_count"`
	UniqueBaseImages       []string `json:"unique_base_images"`
	MostCommonInstructions []string `json:"most_common_instructions"`
	Recommendations        []string `json:"recommendations"`
	OverallScore           int      `json:"overall_score"` // 0-100
}