# Simple test script to verify the service is running
$url = "http://localhost:9190/roll"
$body = @{
    expression = "2d6 + 1d4 - d8 + 5"
} | ConvertTo-Json

Write-Host "Testing dice service at $url" -ForegroundColor Cyan
Write-Host "Request body: $body" -ForegroundColor Gray
Write-Host ""

try {
    $response = Invoke-RestMethod -Uri $url -Method Post -Body $body -ContentType "application/json"
    Write-Host "Success! Response:" -ForegroundColor Green
    $response | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Make sure the server is running:" -ForegroundColor Yellow
    Write-Host "  go run ./cmd/server" -ForegroundColor Yellow
    Write-Host "  or" -ForegroundColor Yellow
    Write-Host "  .\start-server.bat" -ForegroundColor Yellow
}








