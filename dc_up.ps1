Write-Host "Building go web..."
docker build -t docker-test-web .

New-Item -f -type d c:\Users\docker-test

Write-Host "Copying config files to shared directory..."
cp ./config/prometheus.yml c:\Users\docker-test\prometheus.yml

Write-Host "Calling docker-compose..."
docker-compose up
