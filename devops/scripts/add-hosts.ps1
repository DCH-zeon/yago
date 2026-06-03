param([string]$Domain)

$entries = @(
    "127.0.0.1 $Domain",
    "127.0.0.1 admin.$Domain",
    "127.0.0.1 storage.$Domain",
    "127.0.0.1 api.$Domain"
)

foreach ($entry in $entries) {
    "$entry" | Out-File -FilePath C:\Windows\System32\drivers\etc\hosts -Append -Encoding ASCII
}