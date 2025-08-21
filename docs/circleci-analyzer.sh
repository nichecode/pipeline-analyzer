dependency information
    echo "6. **Test task equivalents** locally"
    echo ""
    
    echo "## Key Files to Examine"
    echo ""
    echo "- [Job Usage Analysis](summaries/job-usage.md) - How jobs are reused and their dependencies"
    echo "- [All Jobs](summaries/all-jobs.md) - Complete job list with descriptions"
    echo "- [Docker & Scripts](summaries/docker-and-scripts.md) - Docker/script patterns"
    echo "- [Workflows](summaries/workflows.md) - Workflow structure and job orchestration"
    echo ""
    
    echo "## Suggested go-task Structure"
    echo ""
    echo '```yaml'
    echo "version: '3'"
    echo ""
    echo "tasks:"
    while IFS= read -r job; do
        if [[ -n "$job" ]]; then
            desc=$(yq ".jobs.\"$job\".description" "$CONFIG_FILE" 2>/dev/null)
            if [[ -n "$desc" && "$desc" != "null" ]]; then
                echo "  $job:"
                echo "    desc: \"$desc\""
            else
                echo "  $job:"
                echo "    desc: \"Migrated from CircleCI job\""
            fi
            
            # Try to detect dependencies
            deps=$(yq ".workflows[].jobs.\"$job\".requires[]?" "$CONFIG_FILE" 2>/dev/null | sort | uniq | tr '\n' ', ' | sed 's/,$//')
            if [[ -n "$deps" ]]; then
                echo "    deps: [$deps]"
            fi
            echo "    cmds:"
            echo "      - # Convert run commands from jobs/$job.md"
            echo ""
        fi
    done < "$OUTPUT_DIR/summaries/unique-jobs.txt"
    echo '```'
    echo ""
    
    echo "## Navigation"
    echo ""
    echo "- [← Back to Overview](../README.md)"
    echo "- [All Jobs](summaries/all-jobs.md)"
    echo "- [Job Usage Analysis](summaries/job-usage.md)"
    
} > "$OUTPUT_DIR/MIGRATION-CHECKLIST.md"

# 11. Create main README
echo "11. Creating main README..."
{
    echo "# CircleCI Analysis Report"
    echo ""
    echo "**Generated:** $(date)"
    echo "**Config:** $CONFIG_FILE"
    echo ""
    
    echo "## 📊 Overview"
    echo ""
    echo "- **Unique jobs:** $JOB_COUNT"
    if [[ -s "$OUTPUT_DIR/summaries/workflow-names.txt" ]]; then
        workflow_list=$(cat "$OUTPUT_DIR/summaries/workflow-names.txt" | tr '\n' ', ' | sed 's/,$//')
        echo "- **Workflows:** $workflow_list"
    else
        echo "- **Workflows:** None found"
    fi
    echo ""
    
    echo "## 🚀 Quick Start"
    echo ""
    echo "1. **[📋 Migration Checklist](MIGRATION-CHECKLIST.md)** - Your step-by-step guide"
    echo "2. **[📈 Job Usage Analysis](summaries/job-usage.md)** - Job reuse patterns and dependencies"
    echo "3. **[⚡ Commands Analysis](summaries/commands.md)** - All run commands for conversion"
    echo ""
    
    echo "## 📁 Directory Structure"
    echo ""
    echo "### Jobs"
    echo "Individual job analysis with run commands and configuration:"
    echo ""
    while IFS= read -r job; do
        if [[ -n "$job" ]]; then
            echo "- [jobs/$job.md](jobs/$job.md)"
        fi
    done < "$OUTPUT_DIR/summaries/unique-jobs.txt"
    echo ""
    
    echo "### Workflows"
    echo "Workflow structure and job dependencies:"
    echo ""
    if [[ -s "$OUTPUT_DIR/summaries/workflow-names.txt" ]]; then
        while IFS= read -r workflow; do
            if [[ -n "$workflow" ]]; then
                echo "- [workflows/$workflow.md](workflows/$workflow.md)"
            fi
        done < "$OUTPUT_DIR/summaries/workflow-names.txt"
    else
        echo "- No workflows found"
    fi
    echo ""
    
    echo "### Analysis Summaries"
    echo ""
    echo "- [📈 Job Usage & Dependencies](summaries/job-usage.md)"
    echo "- [📝 All Jobs Index](summaries/all-jobs.md)"
    echo "- [⚡ Commands Analysis](summaries/commands.md)"
    echo "- [🐳 Docker & Scripts](summaries/docker-and-scripts.md)"
    echo "- [⚙️ Executors & Images](summaries/executors-and-images.md)"
    echo "- [🔄 Workflows Index](summaries/workflows.md)"
    echo ""
    
    echo "## 🎯 Next Steps"
    echo ""
    echo "1. **Start with [Migration Checklist](MIGRATION-CHECKLIST.md)**"
    echo "2. **Review most frequently used jobs** from [job usage analysis](summaries/job-usage.md)"
    echo "3. **Examine job dependencies** to understand execution order"
    echo "4. **Begin converting** highest-impact jobs to go-task format"
    echo ""
    
    echo "## 🔍 Most Used Jobs"
    echo ""
    if yq '.workflows' "$CONFIG_FILE" >/dev/null 2>&1; then
        echo "| Job | Usage Count | Link |"
        echo "|-----|-------------|------|"
        yq '.workflows[].jobs[]' "$CONFIG_FILE" 2>/dev/null | \
        grep -E '^[a-zA-Z_-]+' | \
        sort | uniq -c | sort -nr | head -10 | \
        while read -r count job; do
            echo "| $job | $count | [View Details](jobs/$job.md) |"
        done
    else
        echo "No workflow data available"
    fi
    
} > "$OUTPUT_DIR/README.md"

echo ""
echo "✅ Analysis complete!"
echo ""
echo "📁 **Output directory:** $OUTPUT_DIR"
echo "📖 **Start here:** $OUTPUT_DIR/README.md"
echo ""
echo "🔗 **Key entry points:**"
echo "   - 📋 Migration guide: $OUTPUT_DIR/MIGRATION-CHECKLIST.md"
echo "   - 📈 Job analysis: $OUTPUT_DIR/summaries/job-usage.md"
echo "   - 📝 All jobs: $OUTPUT_DIR/summaries/all-jobs.md"
echo ""
echo "🚀 **Ready to start your go-task migration!**"
echo "   Open $OUTPUT_DIR/README.md in your browser or markdown viewer."