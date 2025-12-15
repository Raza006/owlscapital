param(
    [Parameter(Position = 0)]
    [ValidateSet(
        "help", "up", "down", "up-postgres", "restart", "restart-postgres",
        "logs", "logs-postgres", "build", "ps", "pull", "clean", "recreate"
    )]
    [string]$Command = "help"
)

$composeFile = Join-Path $PSScriptRoot "services\docker-compose.yml"

function Invoke-Compose {
    param([string[]]$ComposeArgs)
    & docker compose -f $composeFile @ComposeArgs
}

switch ($Command) {
    "help" {
        Write-Host "Usage: .\make.ps1 <command>"
        Write-Host ""
        Write-Host "Available commands:"
        Write-Host "  up                Build and start all services"
        Write-Host "  down              Stop and remove services"
        Write-Host "  up-postgres       Start only the postgres service"
        Write-Host "  restart           Restart the discordbot service"
        Write-Host "  restart-postgres  Restart the postgres service"
        Write-Host "  logs              Tail discordbot logs"
        Write-Host "  logs-postgres     Tail postgres logs"
        Write-Host "  build             Build the discordbot image"
        Write-Host "  ps                Show service status"
        Write-Host "  pull              Pull remote images"
        Write-Host "  clean             Remove containers, networks, and volumes"
        Write-Host "  recreate          Force recreate containers"
    }
    "up" { Invoke-Compose @("up", "-d", "--build") }
    "down" { Invoke-Compose @("down") }
    "up-postgres" { Invoke-Compose @("up", "-d", "postgres") }
    "restart" { Invoke-Compose @("restart", "discordbot") }
    "restart-postgres" { Invoke-Compose @("restart", "postgres") }
    "logs" { Invoke-Compose @("logs", "-f", "discordbot") }
    "logs-postgres" { Invoke-Compose @("logs", "-f", "postgres") }
    "build" { Invoke-Compose @("build", "discordbot") }
    "ps" { Invoke-Compose @("ps") }
    "pull" { Invoke-Compose @("pull") }
    "clean" { Invoke-Compose @("down", "-v") }
    "recreate" { Invoke-Compose @("up", "-d", "--force-recreate") }
    default {
        Write-Host "Unknown command: $Command"
        exit 1
    }
}

