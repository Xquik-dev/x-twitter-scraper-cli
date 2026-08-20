Register-ArgumentCompleter -Native -CommandName __APPNAME__ -ScriptBlock {
  param($wordToComplete, $commandAst, $cursorPosition)

  $elements = $commandAst.CommandElements
  $completionArgs = @()

  # Collect command arguments.
  for ($i = 0; $i -lt $elements.Count; $i++) {
    $completionArgs += $elements[$i].Extent.Text
  }

  # Preserve trailing space to distinguish namespaced commands from subcommands.
  if ($wordToComplete.Length -eq 0 -and $elements.Count -gt 0) {
    $completionArgs += ""
  }

  $output = & {
    $env:COMPLETION_STYLE = 'pwsh'
    __APPNAME__ __complete @completionArgs 2>&1
  }
  $exitCode = $LASTEXITCODE

  # Detect file references anywhere in the token.
  $prefix = ""
  $filePart = $wordToComplete
  $forceFileCompletion = $false

  # Strip quotes for matching, but preserve them in completion results.
  $wordContent = $wordToComplete
  $leadingQuote = ""
  if ($wordToComplete -match '^([''"])(.*)(\1)$') {
    # Handle a fully quoted token.
    $leadingQuote = $Matches[1]
    $wordContent = $Matches[2]
  } elseif ($wordToComplete -match '^([''"])(.*)$') {
    # Handle an opening quote.
    $leadingQuote = $Matches[1]
    $wordContent = $Matches[2]
  }

  if ($wordContent -match '^(.*)@(file://|data://)?(.*)$') {
    $prefix = $leadingQuote + $Matches[1] + '@' + $Matches[2]
    $filePart = $Matches[3]
    $forceFileCompletion = $true
  }

  if ($forceFileCompletion) {
    # List the current directory for an empty file path.
    $items = if ([string]::IsNullOrEmpty($filePart)) {
      Get-ChildItem -ErrorAction SilentlyContinue
    } else {
      Get-ChildItem -Path "$filePart*" -ErrorAction SilentlyContinue
    }
    $items | ForEach-Object {
      $completionText = if ($_.PSIsContainer) { $prefix + $_.Name + "/" } else { $prefix + $_.Name }
      [System.Management.Automation.CompletionResult]::new(
        $completionText,
        $completionText,
        'ProviderItem',
        $completionText
      )
    }
  } else {
    switch ($exitCode) {
      10 {
        # Complete files.
        $items = if ([string]::IsNullOrEmpty($wordToComplete)) {
          Get-ChildItem -ErrorAction SilentlyContinue
        } else {
          Get-ChildItem -Path "$wordToComplete*" -ErrorAction SilentlyContinue
        }
        $items | ForEach-Object {
          $completionText = if ($_.PSIsContainer) { $_.Name + "/" } else { $_.Name }
          [System.Management.Automation.CompletionResult]::new(
            $completionText,
            $completionText,
            'ProviderItem',
            $completionText
          )
        }
      }
      11 {
        # Disable completion.
        [System.Management.Automation.CompletionResult]::new(' ', ' ', 'ParameterValue', ' ')
      }
      default {
        # Use command completions.
        $output | ForEach-Object {
          [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
      }
    }
  }
}
