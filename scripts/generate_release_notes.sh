#!/bin/bash

# Script to generate release notes based on git tags and commit messages

# Check for git in the environment
if ! command -v git &> /dev/null
then
    echo "Git is not installed. Please install Git to use this script."
    exit 1
fi

# Function to get the latest tag
get_latest_tag() {
    git describe --tags --abbrev=0
}

# Function to list tags in descending order
get_all_tags() {
    git tag --sort=-v:refname
}

# Function to get commit messages between two tags
get_commits_between_tags() {
    local previous_tag=$1
    local current_tag=$2
    git --no-pager log --pretty=format:"- %s" ${previous_tag}..${current_tag}
}

output_release_notes() {
    local all_tags=( $(get_all_tags) )
    
    echo "## Release Notes"

    if [[ ${#all_tags[@]} -gt 0 ]]; then
        local latest_tag=${all_tags[0]}
        echo -e "\n### Latest Release: ${latest_tag}"
        get_commits_between_tags ${all_tags[1]} ${latest_tag} | sed 's/^/  /' || echo "No commits found."

        if [[ ${#all_tags[@]} -gt 1 ]]; then
            echo -e "\n### Previous Highlights"
            for ((i=1; i<${#all_tags[@]}-1; i++)); do
                echo -n "- **${all_tags[$i]}**: "
                get_commits_between_tags ${all_tags[$((i+1))]} ${all_tags[$i]} | sed -e '1!b' -e 's/^/  /' -e 's/$/./' || echo "  No commits found."
                echo # Add a blank line between versions
            done
        fi
    else
        echo "No tags found in the repository."
    fi
}

# Run the release notes generation
output_release_notes
