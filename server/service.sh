
cat > /usr/lib/systemd/system/checkmate.service <<EOL
[Unit]
Description=checkmate
Name=checkmate

[Service]
Restart=always
ExecStart=/usr/bin/heimdall-aggregator

[Install]
WantedBy=multi-user.target
EOL

echo "building checkmate"
go build .

echo "enabling systemd service"
systemctl enable heimdall.service

echo "starting systemd service"
systemctl start heimdall.service
