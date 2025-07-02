#!/bin/bash

# Enhanced Script to generate release notes based on git tags and commit messages
# Supports commit categorization, release dates, statistics, and filtering

set -euo pipefail

# Color codes for better output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default configuration
SHOW_STATS=false
SHOW_CONTRIBUTORS=false
INCLUDE_MERGE_COMMITS=false
MAX_RELEASES=50
OUTPUT_FORMAT="markdown"
VERBOSE=false
DETAILED_OUTPUT=true

# Help function
show_help() {
    cat << EOF
Enhanced Release Notes Generator for Kdeps Schema

Usage: $0 [OPTIONS]

Options:
    -s, --stats              Include release statistics (commit count, contributors)
    -c, --contributors       Show contributors for each release
    -m, --include-merges     Include merge commits in output
    -n, --max-releases NUM   Maximum number of releases to show (default: 50)
    -f, --format FORMAT      Output format: markdown, json, plain (default: markdown)
    -d, --detailed           Show detailed commit messages and descriptions
    --summary                Show only summary (opposite of --detailed)
    -v, --verbose            Verbose output
    -h, --help              Show this help message

Examples:
    $0                              # Generate standard release notes
    $0 --stats --contributors       # Include statistics and contributors
    $0 --max-releases 5             # Show only last 5 releases
    $0 --format json > releases.json # Output as JSON

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -s|--stats)
            SHOW_STATS=true
            shift
            ;;
        -c|--contributors)
            SHOW_CONTRIBUTORS=true
            shift
            ;;
        -m|--include-merges)
            INCLUDE_MERGE_COMMITS=true
            shift
            ;;
        -n|--max-releases)
            MAX_RELEASES="$2"
            shift 2
            ;;
        -f|--format)
            OUTPUT_FORMAT="$2"
            shift 2
            ;;
        -d|--detailed)
            DETAILED_OUTPUT=true
            shift
            ;;
        --summary)
            DETAILED_OUTPUT=false
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}Error: Unknown option $1${NC}" >&2
            show_help
            exit 1
            ;;
    esac
done

# Logging function
log() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${BLUE}[INFO]${NC} $1" >&2
    fi
}

# Check for git in the environment
if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: Git is not installed. Please install Git to use this script.${NC}" >&2
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository.${NC}" >&2
    exit 1
fi

# Function to get the latest tag
get_latest_tag() {
    git describe --tags --abbrev=0 2>/dev/null || echo ""
}

# Function to list tags in descending order
get_all_tags() {
    git fetch --all --quiet 2>/dev/null || true
    git tag --sort=-v:refname | head -n "$MAX_RELEASES"
}

# Function to get tag date
get_tag_date() {
    local tag=$1
    git log -1 --format=%ai "$tag" 2>/dev/null | cut -d' ' -f1
}

# Function to get commit count between tags
get_commit_count() {
    local previous_tag=$1
    local current_tag=$2
    if [[ "$INCLUDE_MERGE_COMMITS" == "true" ]]; then
        git rev-list --count ${previous_tag}..${current_tag} 2>/dev/null || echo "0"
    else
        git rev-list --count --no-merges ${previous_tag}..${current_tag} 2>/dev/null || echo "0"
    fi
}

# Function to get contributors between tags
get_contributors() {
    local previous_tag=$1
    local current_tag=$2
    git shortlog -sn ${previous_tag}..${current_tag} 2>/dev/null | head -5 | sed 's/^[[:space:]]*[0-9]*[[:space:]]*/  - /'
}

# Function to categorize commits
categorize_commit() {
    local commit_msg="$1"
    local commit_hash="$2"
    
    # Convert to lowercase for pattern matching
    local lower_msg=$(echo "$commit_msg" | tr '[:upper:]' '[:lower:]')
    
    # Categorization patterns
    if [[ "$lower_msg" =~ ^(feat|feature)(\(.+\))?!?: ]]; then
        echo "🚀 **Features**"
    elif [[ "$lower_msg" =~ ^(fix|bugfix)(\(.+\))?!?: ]]; then
        echo "🐛 **Bug Fixes**"
    elif [[ "$lower_msg" =~ ^(break|breaking)(\(.+\))?!?: ]] || [[ "$lower_msg" =~ ! ]]; then
        echo "💥 **Breaking Changes**"
    elif [[ "$lower_msg" =~ ^(docs?)(\(.+\))?!?: ]]; then
        echo "📚 **Documentation**"
    elif [[ "$lower_msg" =~ ^(test|tests)(\(.+\))?!?: ]]; then
        echo "🧪 **Tests**"
    elif [[ "$lower_msg" =~ ^(refactor|ref)(\(.+\))?!?: ]]; then
        echo "♻️ **Refactoring**"
    elif [[ "$lower_msg" =~ ^(perf|performance)(\(.+\))?!?: ]]; then
        echo "⚡ **Performance**"
    elif [[ "$lower_msg" =~ ^(chore|build|ci)(\(.+\))?!?: ]]; then
        echo "🔧 **Maintenance**"
    elif [[ "$lower_msg" =~ (add|new|implement) ]]; then
        echo "✨ **Enhancements**"
    elif [[ "$lower_msg" =~ (update|upgrade|bump) ]]; then
        echo "📦 **Updates**"
    else
        echo "📝 **Other Changes**"
    fi
}

# Function to get commit messages between two tags with categorization
get_commits_between_tags() {
    local previous_tag=$1
    local current_tag=$2
    local merge_flag=""
    
    if [[ "$INCLUDE_MERGE_COMMITS" != "true" ]]; then
        merge_flag="--no-merges"
    fi
    
    # Get commits with hash, message, and body if detailed
    if [[ "$DETAILED_OUTPUT" == "true" ]]; then
        git log ${merge_flag} --pretty=format:"%h|%s|%b" ${previous_tag}..${current_tag} 2>/dev/null
    else
        git log ${merge_flag} --pretty=format:"%h|%s|" ${previous_tag}..${current_tag} 2>/dev/null
    fi
}

# Function to get all commits from beginning if no previous tag
get_all_commits_from_beginning() {
    local current_tag=$1
    local merge_flag=""
    
    if [[ "$INCLUDE_MERGE_COMMITS" != "true" ]]; then
        merge_flag="--no-merges"
    fi
    
    # Get all commits up to the tag
    if [[ "$DETAILED_OUTPUT" == "true" ]]; then
        git log ${merge_flag} --pretty=format:"%h|%s|%b" ${current_tag} 2>/dev/null
    else
        git log ${merge_flag} --pretty=format:"%h|%s|" ${current_tag} 2>/dev/null
    fi
}

# Function to format commits with categorization (bash 3.x compatible)
format_commits_categorized() {
    local commits="$1"
    
    if [[ -z "$commits" ]]; then
        echo "  - No commits found"
        return
    fi
    
    # Create temporary files for each category (bash 3.x compatible)
    local temp_dir=$(mktemp -d)
    local breaking_file="$temp_dir/breaking"
    local features_file="$temp_dir/features"
    local enhancements_file="$temp_dir/enhancements"
    local bugfixes_file="$temp_dir/bugfixes"
    local performance_file="$temp_dir/performance"
    local refactoring_file="$temp_dir/refactoring"
    local updates_file="$temp_dir/updates"
    local tests_file="$temp_dir/tests"
    local docs_file="$temp_dir/docs"
    local maintenance_file="$temp_dir/maintenance"
    local other_file="$temp_dir/other"
    
    # Process each commit
    while IFS='|' read -r hash message body; do
        [[ -z "$hash" ]] && continue
        
        local category=$(categorize_commit "$message" "$hash")
        local formatted_line="  - **$message** (\`$hash\`)"
        
        # Add detailed body if available and detailed output is enabled
        if [[ "$DETAILED_OUTPUT" == "true" && -n "$body" ]]; then
            # Clean up body text and indent it
            local clean_body=$(echo "$body" | sed '/^$/d' | sed 's/^/    /')
            if [[ -n "$clean_body" ]]; then
                formatted_line="$formatted_line"$'\n'"$clean_body"
            fi
        fi
        
        case "$category" in
            "💥 **Breaking Changes**") echo "$formatted_line" >> "$breaking_file" ;;
            "🚀 **Features**") echo "$formatted_line" >> "$features_file" ;;
            "✨ **Enhancements**") echo "$formatted_line" >> "$enhancements_file" ;;
            "🐛 **Bug Fixes**") echo "$formatted_line" >> "$bugfixes_file" ;;
            "⚡ **Performance**") echo "$formatted_line" >> "$performance_file" ;;
            "♻️ **Refactoring**") echo "$formatted_line" >> "$refactoring_file" ;;
            "📦 **Updates**") echo "$formatted_line" >> "$updates_file" ;;
            "🧪 **Tests**") echo "$formatted_line" >> "$tests_file" ;;
            "📚 **Documentation**") echo "$formatted_line" >> "$docs_file" ;;
            "🔧 **Maintenance**") echo "$formatted_line" >> "$maintenance_file" ;;
            *) echo "$formatted_line" >> "$other_file" ;;
        esac
    done <<< "$commits"
    
    # Output categories in preferred order
    local categories=(
        "💥 **Breaking Changes**:$breaking_file"
        "🚀 **Features**:$features_file"
        "✨ **Enhancements**:$enhancements_file"
        "🐛 **Bug Fixes**:$bugfixes_file"
        "⚡ **Performance**:$performance_file"
        "♻️ **Refactoring**:$refactoring_file"
        "📦 **Updates**:$updates_file"
        "🧪 **Tests**:$tests_file"
        "📚 **Documentation**:$docs_file"
        "🔧 **Maintenance**:$maintenance_file"
        "📝 **Other Changes**:$other_file"
    )
    
    for cat_file in "${categories[@]}"; do
        local category_name="${cat_file%:*}"
        local file_path="${cat_file#*:}"
        
        if [[ -f "$file_path" && -s "$file_path" ]]; then
            echo ""
            echo "$category_name"
            cat "$file_path"
        fi
    done
    
    # Clean up temporary files
    rm -rf "$temp_dir"
}

# Function to format commits (legacy format) 
format_commits_legacy() {
    local commits=$1
    echo "${commits}" | awk -F'|' '
    {
        if (NF >= 2) {
            if (length($3) > 0 && detailed == "true") {
                print "  - **" $2 "** (`" $1 "`)"
                gsub(/\\n/, "\n    ", $3)
                if (length($3) > 0) print "    " $3
            } else {
                print "  - **" $2 "** (`" $1 "`)"
            }
        }
    }' detailed="$DETAILED_OUTPUT"
}

# Function to output release notes in different formats
output_release_notes() {
    local all_tags=( $(get_all_tags) )
    
    case "$OUTPUT_FORMAT" in
        "json")
            output_json_format "${all_tags[@]}"
            ;;
        "plain")
            output_plain_format "${all_tags[@]}"
            ;;
        *)
            output_markdown_format "${all_tags[@]}"
            ;;
    esac
}

# Markdown format output
output_markdown_format() {
    local all_tags=("$@")
    
    cat <<EOF
# Kdeps Schema

This is the schema definitions used by [kdeps](https://kdeps.com).
See the [schema documentation](https://kdeps.github.io/schema).

## What is Kdeps?

Kdeps is an AI Agent framework for building self-hosted RAG AI Agents powered by open-source LLMs.

## Release Notes
EOF

    if [[ ${#all_tags[@]} -gt 0 ]]; then
        local latest_tag=${all_tags[0]}
        local latest_date=$(get_tag_date "$latest_tag")
        
        echo ""
        echo "### Latest Release: ${latest_tag}"
        [[ -n "$latest_date" ]] && echo "*Released: $latest_date*"
        echo ""
        
        if [[ ${#all_tags[@]} -gt 1 ]]; then
            local previous_tag=${all_tags[1]}
            local latest_commits=$(get_commits_between_tags "$previous_tag" "$latest_tag")
            
            # Show statistics if requested
            if [[ "$SHOW_STATS" == "true" ]]; then
                local commit_count=$(get_commit_count "$previous_tag" "$latest_tag")
                echo "**📊 Release Statistics:**"
                echo "- Commits: $commit_count"
                echo ""
            fi
            
            # Show contributors if requested
            if [[ "$SHOW_CONTRIBUTORS" == "true" ]]; then
                echo "**👥 Contributors:**"
                get_contributors "$previous_tag" "$latest_tag"
                echo ""
            fi
            
            # Show categorized commits
            format_commits_categorized "$latest_commits"
        else
            echo "  - Initial release"
        fi

        if [[ ${#all_tags[@]} -gt 1 ]]; then
            echo ""
            echo "### Complete Release History"
            echo ""
            echo "*Detailed changelog showing all changes from the beginning of the project*"
            echo ""
            
            for ((i=1; i<${#all_tags[@]} && i<MAX_RELEASES; i++)); do
                local current_tag=${all_tags[$i]}
                local tag_date=$(get_tag_date "$current_tag")
                local date_str=""
                [[ -n "$tag_date" ]] && date_str=" (*$tag_date*)"
                
                echo ""
                echo "## ${current_tag}${date_str}"
                
                if [[ $((i+1)) -lt ${#all_tags[@]} ]]; then
                    local prev_tag=${all_tags[$((i+1))]}
                    local tag_commits=$(get_commits_between_tags "$prev_tag" "$current_tag")
                    
                    if [[ "$SHOW_STATS" == "true" ]]; then
                        local commit_count=$(get_commit_count "$prev_tag" "$current_tag")
                        echo ""
                        echo "**📊 Release Statistics:**"
                        echo "- Commits: $commit_count"
                        echo ""
                    fi
                    
                    if [[ "$SHOW_CONTRIBUTORS" == "true" ]]; then
                        echo "**👥 Contributors:**"
                        get_contributors "$prev_tag" "$current_tag"
                        echo ""
                    fi
                    
                    if [[ "$DETAILED_OUTPUT" == "true" ]]; then
                        format_commits_categorized "$tag_commits"
                    else
                        format_commits_legacy "$tag_commits"
                    fi
                else
                    # This is the first release - show all commits from beginning
                    echo ""
                    echo "**📊 Initial Release Statistics:**"
                    local all_commits=$(get_all_commits_from_beginning "$current_tag")
                    local total_commits=$(echo "$all_commits" | wc -l | xargs)
                    echo "- Total commits: $total_commits"
                    echo "- Project inception"
                    echo ""
                    
                    if [[ "$SHOW_CONTRIBUTORS" == "true" ]]; then
                        echo "**👥 All Contributors:**"
                        git shortlog -sn "$current_tag" 2>/dev/null | sed 's/^[[:space:]]*[0-9]*[[:space:]]*/  - /'
                        echo ""
                    fi
                    
                    echo "**📝 All Changes Since Project Start:**"
                    if [[ "$DETAILED_OUTPUT" == "true" ]]; then
                        format_commits_categorized "$all_commits"
                    else
                        format_commits_legacy "$all_commits"
                    fi
                fi
            done
        fi
    else
        echo ""
        echo "No tags found in the repository."
    fi
    
    # Add generation timestamp
    echo ""
    echo "---"
    echo "*Generated on $(date '+%Y-%m-%d %H:%M:%S') by [Enhanced Release Notes Generator](scripts/generate_release_notes.sh)*"
}

# JSON format output
output_json_format() {
    local all_tags=("$@")
    
    echo "{"
    echo '  "project": "Kdeps Schema",'
    echo '  "generated": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",'
    echo '  "releases": ['
    
    local first=true
    for ((i=0; i<${#all_tags[@]} && i<MAX_RELEASES; i++)); do
        [[ "$first" != "true" ]] && echo ","
        first=false
        
        local tag=${all_tags[$i]}
        local tag_date=$(get_tag_date "$tag")
        
        echo "    {"
        echo '      "version": "'$tag'",'
        echo '      "date": "'$tag_date'",'
        
        if [[ $((i+1)) -lt ${#all_tags[@]} ]]; then
            local prev_tag=${all_tags[$((i+1))]}
            local commit_count=$(get_commit_count "$prev_tag" "$tag")
            echo '      "commits": '$commit_count','
            
            if [[ "$SHOW_CONTRIBUTORS" == "true" ]]; then
                echo '      "contributors": ['
                get_contributors "$prev_tag" "$tag" | sed 's/^  - /        "/' | sed 's/$/",/' | sed '$s/,$//'
                echo "      ],"
            fi
        fi
        
        echo '      "latest": '$([[ $i -eq 0 ]] && echo "true" || echo "false")
        echo -n "    }"
    done
    
    echo ""
    echo "  ]"
    echo "}"
}

# Plain text format output
output_plain_format() {
    local all_tags=("$@")
    
    echo "KDEPS SCHEMA RELEASE NOTES"
    echo "========================="
    echo ""
    
    for ((i=0; i<${#all_tags[@]} && i<MAX_RELEASES; i++)); do
        local tag=${all_tags[$i]}
        local tag_date=$(get_tag_date "$tag")
        
        echo "Version: $tag"
        [[ -n "$tag_date" ]] && echo "Date: $tag_date"
        
        if [[ $((i+1)) -lt ${#all_tags[@]} ]]; then
            local prev_tag=${all_tags[$((i+1))]}
            local commits=$(get_commits_between_tags "$prev_tag" "$tag")
            echo "Changes:"
            format_commits_legacy "$commits"
        else
            echo "Changes: Initial release"
        fi
        
        echo ""
        echo "----------------------------------------"
        echo ""
    done
}

# Main execution
main() {
    log "Starting enhanced release notes generation..."
    log "Configuration: stats=$SHOW_STATS, contributors=$SHOW_CONTRIBUTORS, format=$OUTPUT_FORMAT"
    
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${GREEN}Enhanced Release Notes Generator${NC}" >&2
        echo -e "${YELLOW}Repository: $(git remote get-url origin 2>/dev/null || echo 'Local repository')${NC}" >&2
        echo -e "${YELLOW}Branch: $(git branch --show-current)${NC}" >&2
        echo "" >&2
    fi
    
    output_release_notes
    
    log "Release notes generation completed successfully!"
}

# Run the main function
main
