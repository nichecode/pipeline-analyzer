package shared

import (
	"regexp"
	"strings"
)

// CommandPattern represents a common command pattern
type CommandPattern struct {
	Name        string
	Regex       *regexp.Regexp
	Category    string
	Description string
	Tools       []string
}

// GetCommonPatterns returns commonly used command patterns
func GetCommonPatterns() map[string]CommandPattern {
	return map[string]CommandPattern{
		"docker": {
			Name:        "docker",
			Regex:       regexp.MustCompile(`\bdocker\s+`),
			Category:    "containerization",
			Description: "Docker containerization commands",
			Tools:       []string{"docker"},
		},
		"docker-compose": {
			Name:        "docker-compose",
			Regex:       regexp.MustCompile(`docker-compose`),
			Category:    "containerization",
			Description: "Docker Compose orchestration commands",
			Tools:       []string{"docker-compose"},
		},
		"go": {
			Name:        "go",
			Regex:       regexp.MustCompile(`\bgo\s+`),
			Category:    "build",
			Description: "Go programming language commands",
			Tools:       []string{"go"},
		},
		"npm": {
			Name:        "npm",
			Regex:       regexp.MustCompile(`\bnpm\s+`),
			Category:    "package-management",
			Description: "Node.js package manager commands",
			Tools:       []string{"npm"},
		},
		"yarn": {
			Name:        "yarn",
			Regex:       regexp.MustCompile(`\byarn\s+`),
			Category:    "package-management",
			Description: "Yarn package manager commands",
			Tools:       []string{"yarn"},
		},
		"make": {
			Name:        "make",
			Regex:       regexp.MustCompile(`\bmake\s+`),
			Category:    "build",
			Description: "GNU Make build commands",
			Tools:       []string{"make"},
		},
		"git": {
			Name:        "git",
			Regex:       regexp.MustCompile(`\bgit\s+`),
			Category:    "version-control",
			Description: "Git version control commands",
			Tools:       []string{"git"},
		},
		"curl": {
			Name:        "curl",
			Regex:       regexp.MustCompile(`\bcurl\s+`),
			Category:    "network",
			Description: "HTTP client commands",
			Tools:       []string{"curl"},
		},
		"kubectl": {
			Name:        "kubectl",
			Regex:       regexp.MustCompile(`\bkubectl\s+`),
			Category:    "deployment",
			Description: "Kubernetes cluster management commands",
			Tools:       []string{"kubectl"},
		},
		"helm": {
			Name:        "helm",
			Regex:       regexp.MustCompile(`\bhelm\s+`),
			Category:    "deployment",
			Description: "Helm Kubernetes package manager commands",
			Tools:       []string{"helm"},
		},
		"terraform": {
			Name:        "terraform",
			Regex:       regexp.MustCompile(`\bterraform\s+`),
			Category:    "infrastructure",
			Description: "Terraform infrastructure as code commands",
			Tools:       []string{"terraform"},
		},
		"python": {
			Name:        "python",
			Regex:       regexp.MustCompile(`\bpython\s+`),
			Category:    "runtime",
			Description: "Python interpreter commands",
			Tools:       []string{"python"},
		},
		"pip": {
			Name:        "pip",
			Regex:       regexp.MustCompile(`\bpip\s+`),
			Category:    "package-management",
			Description: "Python package installer commands",
			Tools:       []string{"pip"},
		},
		"mvn": {
			Name:        "mvn",
			Regex:       regexp.MustCompile(`\bmvn\s+`),
			Category:    "build",
			Description: "Maven build tool commands",
			Tools:       []string{"maven"},
		},
		"gradle": {
			Name:        "gradle",
			Regex:       regexp.MustCompile(`\bgradle\s+`),
			Category:    "build",
			Description: "Gradle build tool commands",
			Tools:       []string{"gradle"},
		},
		"cargo": {
			Name:        "cargo",
			Regex:       regexp.MustCompile(`\bcargo\s+`),
			Category:    "build",
			Description: "Rust package manager and build tool commands",
			Tools:       []string{"cargo"},
		},
		"ssh": {
			Name:        "ssh",
			Regex:       regexp.MustCompile(`\bssh\s+`),
			Category:    "network",
			Description: "Secure Shell remote access commands",
			Tools:       []string{"ssh"},
		},
		"rsync": {
			Name:        "rsync",
			Regex:       regexp.MustCompile(`\brsync\s+`),
			Category:    "file-transfer",
			Description: "Remote file synchronization commands",
			Tools:       []string{"rsync"},
		},
		"local-script": {
			Name:        "local-script",
			Regex:       regexp.MustCompile(`\./[^\s]+`),
			Category:    "script",
			Description: "Local script execution",
			Tools:       []string{"shell"},
		},
		// PHP Ecosystem
		"php": {
			Name:        "php",
			Regex:       regexp.MustCompile(`\bphp\s+`),
			Category:    "runtime",
			Description: "PHP interpreter commands",
			Tools:       []string{"php"},
		},
		"composer": {
			Name:        "composer",
			Regex:       regexp.MustCompile(`\bcomposer\s+`),
			Category:    "package-management",
			Description: "PHP dependency manager commands",
			Tools:       []string{"composer"},
		},
		"phpunit": {
			Name:        "phpunit",
			Regex:       regexp.MustCompile(`\bphpunit\s+`),
			Category:    "testing",
			Description: "PHPUnit testing framework commands",
			Tools:       []string{"phpunit"},
		},
		"pest": {
			Name:        "pest",
			Regex:       regexp.MustCompile(`\bpest\s+`),
			Category:    "testing",
			Description: "Pest PHP testing framework commands",
			Tools:       []string{"pest"},
		},
		"behat": {
			Name:        "behat",
			Regex:       regexp.MustCompile(`\bbehat\s+`),
			Category:    "testing",
			Description: "Behat BDD testing framework commands",
			Tools:       []string{"behat"},
		},
		"codeception": {
			Name:        "codeception",
			Regex:       regexp.MustCompile(`\bcodecept\s+`),
			Category:    "testing",
			Description: "Codeception testing framework commands",
			Tools:       []string{"codeception"},
		},
		"phpcs": {
			Name:        "phpcs",
			Regex:       regexp.MustCompile(`\bphpcs\s+`),
			Category:    "code-quality",
			Description: "PHP Code Sniffer - coding standards checker",
			Tools:       []string{"phpcs"},
		},
		"phpcbf": {
			Name:        "phpcbf",
			Regex:       regexp.MustCompile(`\bphpcbf\s+`),
			Category:    "code-quality",
			Description: "PHP Code Beautifier and Fixer",
			Tools:       []string{"phpcbf"},
		},
		"phpstan": {
			Name:        "phpstan",
			Regex:       regexp.MustCompile(`\bphpstan\s+`),
			Category:    "code-quality",
			Description: "PHPStan static analysis tool",
			Tools:       []string{"phpstan"},
		},
		"psalm": {
			Name:        "psalm",
			Regex:       regexp.MustCompile(`\bpsalm\s+`),
			Category:    "code-quality",
			Description: "Psalm static analysis tool",
			Tools:       []string{"psalm"},
		},
		"phan": {
			Name:        "phan",
			Regex:       regexp.MustCompile(`\bphan\s+`),
			Category:    "code-quality",
			Description: "Phan static analyzer",
			Tools:       []string{"phan"},
		},
		"php-cs-fixer": {
			Name:        "php-cs-fixer",
			Regex:       regexp.MustCompile(`\bphp-cs-fixer\s+`),
			Category:    "code-quality",
			Description: "PHP Coding Standards Fixer",
			Tools:       []string{"php-cs-fixer"},
		},
		"phpmd": {
			Name:        "phpmd",
			Regex:       regexp.MustCompile(`\bphpmd\s+`),
			Category:    "code-quality",
			Description: "PHP Mess Detector - code quality analyzer",
			Tools:       []string{"phpmd"},
		},
		"phpdoc": {
			Name:        "phpdoc",
			Regex:       regexp.MustCompile(`\bphpdoc\s+`),
			Category:    "documentation",
			Description: "phpDocumentor documentation generator",
			Tools:       []string{"phpdoc"},
		},
		"artisan": {
			Name:        "artisan",
			Regex:       regexp.MustCompile(`\bartisan\s+`),
			Category:    "framework",
			Description: "Laravel Artisan command-line interface",
			Tools:       []string{"laravel", "artisan"},
		},
		"symfony": {
			Name:        "symfony",
			Regex:       regexp.MustCompile(`\bsymfony\s+`),
			Category:    "framework",
			Description: "Symfony console commands",
			Tools:       []string{"symfony"},
		},
		"drush": {
			Name:        "drush",
			Regex:       regexp.MustCompile(`\bdrush\s+`),
			Category:    "framework",
			Description: "Drupal shell commands",
			Tools:       []string{"drupal", "drush"},
		},
		"wp": {
			Name:        "wp",
			Regex:       regexp.MustCompile(`\bwp\s+`),
			Category:    "framework",
			Description: "WordPress CLI commands",
			Tools:       []string{"wordpress", "wp-cli"},
		},
		"php-server": {
			Name:        "php-server",
			Regex:       regexp.MustCompile(`php\s+-S\s+`),
			Category:    "development",
			Description: "PHP built-in development server",
			Tools:       []string{"php"},
		},
		"phar": {
			Name:        "phar",
			Regex:       regexp.MustCompile(`\bphar\s+`),
			Category:    "packaging",
			Description: "PHP Archive (PHAR) commands",
			Tools:       []string{"phar"},
		},
		"box": {
			Name:        "box",
			Regex:       regexp.MustCompile(`\bbox\s+`),
			Category:    "packaging",
			Description: "Box PHP application builder",
			Tools:       []string{"box"},
		},
		"infection": {
			Name:        "infection",
			Regex:       regexp.MustCompile(`\binfection\s+`),
			Category:    "testing",
			Description: "Infection mutation testing framework",
			Tools:       []string{"infection"},
		},
		"robo": {
			Name:        "robo",
			Regex:       regexp.MustCompile(`\brobo\s+`),
			Category:    "task-runner",
			Description: "Robo PHP task runner",
			Tools:       []string{"robo"},
		},
		"deployer": {
			Name:        "deployer",
			Regex:       regexp.MustCompile(`\bdep\s+`),
			Category:    "deployment",
			Description: "Deployer deployment tool",
			Tools:       []string{"deployer"},
		},
		"phinx": {
			Name:        "phinx",
			Regex:       regexp.MustCompile(`\bphinx\s+`),
			Category:    "database",
			Description: "Phinx database migrations",
			Tools:       []string{"phinx"},
		},
		"doctrine": {
			Name:        "doctrine",
			Regex:       regexp.MustCompile(`\bdoctrine\s+`),
			Category:    "database",
			Description: "Doctrine ORM commands",
			Tools:       []string{"doctrine"},
		},
	}
}

// ClassifyCommand attempts to classify a command string
func ClassifyCommand(command string) CommandClassification {
	patterns := GetCommonPatterns()
	lowerCommand := strings.ToLower(command)
	
	// Define pattern priority order - more specific patterns first
	priorityOrder := []string{
		// Framework-specific patterns (most specific)
		"artisan", "symfony", "drush", "wp",
		// PHP-specific tools
		"phpunit", "pest", "behat", "codeception", "infection",
		"phpcs", "phpcbf", "phpstan", "psalm", "phan", "php-cs-fixer", "phpmd",
		"phpdoc", "php-server", "phar", "box", "robo", "deployer",
		"phinx", "doctrine", "composer",
		// Other specific tools
		"docker-compose", "kubectl", "helm", "terraform", "cargo",
		"gradle", "local-script",
		// General tools (least specific)
		"docker", "go", "npm", "yarn", "make", "git", "curl", "ssh", "rsync",
		"python", "pip", "mvn", "php",
	}
	
	// Check patterns in priority order
	for _, patternName := range priorityOrder {
		if pattern, exists := patterns[patternName]; exists {
			if pattern.Regex.MatchString(lowerCommand) {
				return CommandClassification{
					Category:    pattern.Category,
					Tools:       pattern.Tools,
					Complexity:  calculateCommandComplexity(command),
					Risk:        assessCommandRisk(command),
					Suggestions: generateCommandSuggestions(command, pattern),
				}
			}
		}
	}
	
	// Fallback to any remaining patterns
	for _, pattern := range patterns {
		if pattern.Regex.MatchString(lowerCommand) {
			return CommandClassification{
				Category:    pattern.Category,
				Tools:       pattern.Tools,
				Complexity:  calculateCommandComplexity(command),
				Risk:        assessCommandRisk(command),
				Suggestions: generateCommandSuggestions(command, pattern),
			}
		}
	}
	
	// Default classification
	return CommandClassification{
		Category:   "utility",
		Tools:      []string{"shell"},
		Complexity: calculateCommandComplexity(command),
		Risk:       assessCommandRisk(command),
	}
}

// CommandClassification represents command analysis results
type CommandClassification struct {
	Category    string
	Tools       []string
	Complexity  int
	Risk        string
	Suggestions []string
}

// calculateCommandComplexity estimates command complexity
func calculateCommandComplexity(command string) int {
	complexity := 1
	
	// Pipe operations increase complexity
	if strings.Contains(command, "|") {
		complexity += strings.Count(command, "|")
	}
	
	// Redirection increases complexity
	if strings.Contains(command, ">") || strings.Contains(command, ">>") {
		complexity++
	}
	
	// Background processes increase complexity
	if strings.Contains(command, "&") {
		complexity++
	}
	
	// Conditional execution increases complexity
	if strings.Contains(command, "&&") || strings.Contains(command, "||") {
		complexity += strings.Count(command, "&&") + strings.Count(command, "||")
	}
	
	// Command substitution increases complexity
	if strings.Contains(command, "$(") || strings.Contains(command, "`") {
		complexity++
	}
	
	// Multiple statements increase complexity
	if strings.Contains(command, ";") {
		complexity += strings.Count(command, ";")
	}
	
	return complexity
}

// assessCommandRisk evaluates potential security risks
func assessCommandRisk(command string) string {
	lowerCommand := strings.ToLower(command)
	
	// High risk patterns
	highRisk := []string{
		"rm -rf",
		"dd if=",
		"mkfs",
		"fdisk",
		"format",
		"del /f",
		"rmdir /s",
		"sudo rm",
		"chmod 777",
		"curl | sh",
		"wget | sh",
		"eval",
		"> /dev/",
	}
	
	for _, risk := range highRisk {
		if strings.Contains(lowerCommand, risk) {
			return "high"
		}
	}
	
	// Medium risk patterns
	mediumRisk := []string{
		"sudo",
		"su ",
		"chmod",
		"chown",
		"rm ",
		"mv /",
		"cp /",
		"curl",
		"wget",
		"ssh",
		"scp",
	}
	
	for _, risk := range mediumRisk {
		if strings.Contains(lowerCommand, risk) {
			return "medium"
		}
	}
	
	return "low"
}

// generateCommandSuggestions provides improvement suggestions
func generateCommandSuggestions(command string, pattern CommandPattern) []string {
	var suggestions []string
	
	switch pattern.Category {
	case "containerization":
		if !strings.Contains(command, "--rm") && strings.Contains(command, "docker run") {
			suggestions = append(suggestions, "Consider adding --rm flag to automatically remove containers")
		}
		if strings.Contains(command, "latest") {
			suggestions = append(suggestions, "Avoid using 'latest' tag, specify explicit version")
		}
	
	case "build":
		if strings.Contains(command, "go build") && !strings.Contains(command, "-ldflags") {
			suggestions = append(suggestions, "Consider using -ldflags to inject version information")
		}
		if strings.Contains(command, "npm install") && !strings.Contains(command, "ci") {
			suggestions = append(suggestions, "Use 'npm ci' for faster, reliable builds in CI environments")
		}
	
	case "network":
		if strings.Contains(command, "curl") && !strings.Contains(command, "--fail") {
			suggestions = append(suggestions, "Consider using --fail flag to exit on HTTP errors")
		}
	
	case "package-management":
		if strings.Contains(command, "composer install") && !strings.Contains(command, "--no-dev") {
			suggestions = append(suggestions, "Consider using --no-dev flag for production installs")
		}
		if strings.Contains(command, "composer install") && !strings.Contains(command, "--optimize-autoloader") {
			suggestions = append(suggestions, "Add --optimize-autoloader for better performance")
		}
	
	case "testing":
		if strings.Contains(command, "phpunit") && !strings.Contains(command, "--coverage") {
			suggestions = append(suggestions, "Consider adding code coverage reporting")
		}
		if strings.Contains(command, "pest") && !strings.Contains(command, "--parallel") {
			suggestions = append(suggestions, "Use --parallel flag to run tests in parallel for faster execution")
		}
	
	case "code-quality":
		if strings.Contains(command, "phpcs") && !strings.Contains(command, "--standard") {
			suggestions = append(suggestions, "Specify coding standard with --standard flag")
		}
		if strings.Contains(command, "phpstan") && !strings.Contains(command, "--level") {
			suggestions = append(suggestions, "Specify analysis level with --level flag (0-9)")
		}
	
	case "framework":
		if strings.Contains(command, "artisan") && strings.Contains(command, "migrate") && !strings.Contains(command, "--seed") {
			suggestions = append(suggestions, "Consider using --seed flag to run database seeders")
		}
	}
	
	// General suggestions based on complexity
	complexity := calculateCommandComplexity(command)
	if complexity > 3 {
		suggestions = append(suggestions, "Complex command - consider breaking into multiple steps")
	}
	
	// Security suggestions
	risk := assessCommandRisk(command)
	if risk == "high" {
		suggestions = append(suggestions, "High-risk command detected - review security implications")
	}
	
	return suggestions
}

// GetPatternsByCategory groups patterns by category
func GetPatternsByCategory() map[string][]string {
	patterns := GetCommonPatterns()
	categories := make(map[string][]string)
	
	for name, pattern := range patterns {
		categories[pattern.Category] = append(categories[pattern.Category], name)
	}
	
	return categories
}

// DetectToolEcosystem attempts to identify the primary tool ecosystem
func DetectToolEcosystem(commands []string) string {
	ecosystems := map[string][]string{
		"go":         {"go", "gofmt", "golint", "go-build"},
		"nodejs":     {"npm", "yarn", "node", "pnpm"},
		"python":     {"python", "pip", "poetry", "conda"},
		"java":       {"mvn", "gradle", "java", "javac"},
		"rust":       {"cargo", "rustc", "rustfmt"},
		"php":        {"php", "composer", "phpunit", "artisan", "symfony", "pest", "behat", "phpcs", "phpstan", "psalm"},
		"docker":     {"docker", "docker-compose", "dockerfile"},
		"kubernetes": {"kubectl", "helm", "k8s", "kustomize"},
	}
	
	scores := make(map[string]int)
	
	for _, command := range commands {
		lowerCommand := strings.ToLower(command)
		for ecosystem, tools := range ecosystems {
			for _, tool := range tools {
				if strings.Contains(lowerCommand, tool) {
					scores[ecosystem]++
				}
			}
		}
	}
	
	// Find the highest scoring ecosystem
	maxScore := 0
	primaryEcosystem := "mixed"
	
	for ecosystem, score := range scores {
		if score > maxScore {
			maxScore = score
			primaryEcosystem = ecosystem
		}
	}
	
	if maxScore == 0 {
		return "shell"
	}
	
	return primaryEcosystem
}