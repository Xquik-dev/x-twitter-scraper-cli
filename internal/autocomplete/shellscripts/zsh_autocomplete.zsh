#compdef __APPNAME__

____APPNAME___zsh_autocomplete() {

  local -a opts
  local temp
  local exit_code

  temp=$(COMPLETION_STYLE=zsh "${words[1]}" __complete "${words[@]:1}")
  exit_code=$?

  local cur="${words[CURRENT]}"

  if [[ "$cur" = *'@'* ]]; then
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
      _files
      ;;
    11)
      return 1
      ;;
    0)
      opts=("${(@f)temp}")
      _describe 'values' opts
      ;;
  esac
}

if [[ "${funcstack[1]}" = "___APPNAME__" ]]; then
  ____APPNAME___zsh_autocomplete "$@"
else
  compdef ____APPNAME___zsh_autocomplete __APPNAME__
fi
