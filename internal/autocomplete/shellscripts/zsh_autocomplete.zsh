#compdef __APPNAME__

____APPNAME___zsh_autocomplete() {

  local -a opts
  local temp
  local exit_code

  temp=$(COMPLETION_STYLE=zsh "${words[1]}" __complete "${words[@]:1}")
  exit_code=$?

  # Detect file references anywhere in the token.
  local cur="${words[CURRENT]}"

  if [[ "$cur" = *'@'* ]]; then
    # Read the path after the last @.
    local after_last_at="${cur##*@}"

    if [[ $after_last_at =~ ^(file://|data://) ]]; then
      compset -P "*$MATCH"
      _files
    else
      compset -P '*@'
      _files
    fi
    return
  fi

  case $exit_code in
    10)
      # Complete files.
      _files
      ;;
    11)
      # Disable completion.
      return 1
      ;;
    0)
      # Use command completions.
      opts=("${(@f)temp}")
      _describe 'values' opts
      ;;
  esac
}

# fpath autoloads this body as ___APPNAME__. Sourcing registers it with compdef.
if [[ "${funcstack[1]}" = "___APPNAME__" ]]; then
  ____APPNAME___zsh_autocomplete "$@"
else
  compdef ____APPNAME___zsh_autocomplete __APPNAME__
fi
