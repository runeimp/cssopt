#
# CSS Optimizer
#
PROJECT_NAME := 'CSS Optimizer'

APP_NAME := 'CSSopt'
CLI_CODE := 'cmd/cssopt/cssopt.go'
CLI_NAME := 'cssopt'
MAIN_CODE := 'cssopt.go'

alias ver := version

# set positional-arguments
# set shell := ["pwsh", "-c"] # PowerShell Core (Multi-Platform) # BROKEN ?!
set windows-shell := ["powershell", "-c"] # To use PowerShell Desktop instead of Core on Windows
shebang := if os() == 'windows' { 'powershell' } else { '/usr/bin/env pwsh' } # Shebang for PS Desktop on Windows and PS Core everywhere else

@_default: _term-wipe
	just --list

@args *args:
	echo "\$# = $#"
	echo "\$@ = $@"


# Run code with (optional) arguments
run *args='css/main.css': _term-wipe
	go run {{CLI_CODE}} {{args}}


# Wipe Terminal Buffer and Scrollback Buffer
_term-wipe:
	#!{{shebang}}
	$host.UI.RawUI.WindowTitle = "{{PROJECT_NAME}}"
	Clear-Host

# Unit Test Code
test: _term-wipe
	@# go test ./...
	go test .

tester *args: _term-wipe
	#!{{shebang}}
	Write-Host

	Write-Host "==> Tester with args: {{args}}"
	go run {{CLI_CODE}} {{args}}

	Write-Host


# Display version of app
@version:
	((Get-Content {{MAIN_CODE}} | Select-String 'AppVersion') -Split "'")[1]

