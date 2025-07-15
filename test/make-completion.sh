# Shell completion for Makefile targets in this directory
# Usage: source ./make-completion.sh

_makefile_targets() {
    local makefile
    makefile=Makefile
    if [[ ! -f $makefile ]]; then
        return
    fi
    # Extract targets: lines that start with a word, followed by a colon, not indented, and not special
    grep -E '^[a-zA-Z0-9][^$#\t=:+]*:' "$makefile" | \
        grep -vE '^\.' | \
        sed 's/:.*//' | \
        grep -vE '^(ifeq|else|endif|ifneq|define|endef|include|override|export|unexport|private|vpath|\s*)$' | \
        sort -u
}

_make_target_completion() {
    local cur prev words cword
    _init_completion || return
    local targets
    targets=$(_makefile_targets)
    COMPREPLY=( $(compgen -W "$targets" -- "$cur") )
}

# Bash completion
if [[ -n $BASH_VERSION ]]; then
    complete -F _make_target_completion make
fi

# Zsh completion
if [[ -n $ZSH_VERSION ]]; then
    _makefile_targets_zsh() {
        reply=($(_makefile_targets))
    }
    compctl -K _makefile_targets_zsh make
fi

# Usage message
echo "[make-completion] Loaded. Tab-completion for 'make' targets is now enabled in this shell." 