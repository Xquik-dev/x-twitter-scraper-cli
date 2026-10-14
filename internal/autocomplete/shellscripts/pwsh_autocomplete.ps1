Register-ArgumentCompleter -Native -CommandName __APPNAME__ -ScriptBlock {
  param($wordToComplete, $commandAst, $cursorPosition)

  $elements = $commandAst.CommandElements
  $completionArgs = @()

  for ($i = 0; $i -lt $elements.Count; $i++) {
    $completionArgs += $elements[$i].Extent.Text
  }

  if ($wordToComplete.Length -eq 0 -and $elements.Count -gt 0) {
    $completionArgs += ""
  }

  $output = & {
    $env:COMPLETION_STYLE = 'pwsh'
    __APPNAME__ __complete @completionArgs 2>&1
  }
  $exitCode = $LASTEXITCODE

  $prefix = ""
  $filePart = $wordToComplete
  $forceFileCompletion = $false

  $wordContent = $wordToComplete
  $leadingQuote = ""
  if ($wordToComplete -match '^([''"])(.*)(\1)$') {
    $leadingQuote = $Matches[1]
    $wordContent = $Matches[2]
  } elseif ($wordToComplete -match '^([''"])(.*)$') {
    $leadingQuote = $Matches[1]
    $wordContent = $Matches[2]
  }

  if ($wordContent -match '^(.*)@(file://|data://)?(.*)$') {
    $prefix = $leadingQuote + $Matches[1] + '@' + $Matches[2]
    $filePart = $Matches[3]
    $forceFileCompletion = $true
  }

  if ($forceFileCompletion) {
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
        [System.Management.Automation.CompletionResult]::new(' ', ' ', 'ParameterValue', ' ')
      }
      default {
        $output | ForEach-Object {
          [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
      }
    }
  }
}
