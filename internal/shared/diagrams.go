package shared

import (
	"fmt"
	"strings"
)

// MermaidDiagram generates a simple Mermaid flowchart
type MermaidDiagram struct {
	Title string
	Nodes []MermaidNode
	Edges []MermaidEdge
}

type MermaidNode struct {
	ID          string
	Label       string
	Description string
	Commands    []string
	NodeType    string // workflow, setup, test, build, deploy, utility
}

type MermaidEdge struct {
	From string
	To   string
}

// Generate creates the Mermaid diagram syntax
func (d *MermaidDiagram) Generate() string {
	var sb strings.Builder
	
	sb.WriteString("```mermaid\n")
	sb.WriteString("flowchart TD\n")
	
	// Add nodes
	for _, node := range d.Nodes {
		sb.WriteString(fmt.Sprintf("    %s[\"`**%s**\n", node.ID, node.Label))
		
		if node.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n", node.Description))
		}
		
		// Add key commands (max 3)
		commandCount := len(node.Commands)
		if commandCount > 3 {
			commandCount = 3
		}
		
		for i := 0; i < commandCount; i++ {
			cmd := node.Commands[i]
			if len(cmd) > 40 {
				cmd = cmd[:37] + "..."
			}
			sb.WriteString(fmt.Sprintf("• %s\n", cmd))
		}
		
		if len(node.Commands) > 3 {
			sb.WriteString(fmt.Sprintf("• ... (%d more)\n", len(node.Commands)-3))
		}
		
		sb.WriteString("`\"]\n")
	}
	
	// Add edges
	for _, edge := range d.Edges {
		sb.WriteString(fmt.Sprintf("    %s --> %s\n", edge.From, edge.To))
	}
	
	// Add styling
	sb.WriteString("\n")
	sb.WriteString("    classDef workflow fill:#e1f5fe,stroke:#01579b,stroke-width:3px\n")
	sb.WriteString("    classDef setup fill:#f3e5f5,stroke:#4a148c,stroke-width:2px\n")
	sb.WriteString("    classDef test fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px\n")
	sb.WriteString("    classDef build fill:#fff3e0,stroke:#e65100,stroke-width:2px\n")
	sb.WriteString("    classDef deploy fill:#e0f2f1,stroke:#004d40,stroke-width:2px\n")
	sb.WriteString("    classDef utility fill:#f1f8e9,stroke:#33691e,stroke-width:2px\n")
	
	// Apply classes
	for _, node := range d.Nodes {
		if node.NodeType != "" {
			sb.WriteString(fmt.Sprintf("    class %s %s\n", node.ID, node.NodeType))
		}
	}
	
	sb.WriteString("```\n")
	
	return sb.String()
}

// CleanNodeID creates a valid Mermaid node ID
func CleanNodeID(name string) string {
	// Replace invalid characters with underscores
	cleaned := strings.ReplaceAll(name, "-", "_")
	cleaned = strings.ReplaceAll(cleaned, ":", "_")
	cleaned = strings.ReplaceAll(cleaned, " ", "_")
	cleaned = strings.ReplaceAll(cleaned, "/", "_")
	cleaned = strings.ToUpper(cleaned)
	return cleaned
}

// ClassifyNodeType determines the node type based on name/commands
func ClassifyNodeType(name string, commands []string) string {
	name = strings.ToLower(name)
	allCommands := strings.ToLower(strings.Join(commands, " "))
	
	// Check name patterns
	if strings.Contains(name, "test") {
		return "test"
	}
	if strings.Contains(name, "build") || strings.Contains(name, "compile") {
		return "build"
	}
	if strings.Contains(name, "deploy") || strings.Contains(name, "publish") {
		return "deploy"
	}
	if strings.Contains(name, "install") || strings.Contains(name, "setup") || strings.Contains(name, "clean") {
		return "setup"
	}
	
	// Check command patterns
	if strings.Contains(allCommands, "test") || strings.Contains(allCommands, "lint") {
		return "test"
	}
	if strings.Contains(allCommands, "build") || strings.Contains(allCommands, "docker build") {
		return "build"
	}
	if strings.Contains(allCommands, "deploy") || strings.Contains(allCommands, "kubectl") {
		return "deploy"
	}
	if strings.Contains(allCommands, "npm install") || strings.Contains(allCommands, "pip install") {
		return "setup"
	}
	
	return "utility"
}