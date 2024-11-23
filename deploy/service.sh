
cat > /usr/lib/systemd/system/checkmate.service <<EOL
[Unit]
Description=checkmate
Name=checkmate

[Service]
Restart=always
ExecStart=/home/challenger/hackatum-2024/server/checkmate
WorkingDirectory=/home/challenger/hackatum-2024/server

[Install]
WantedBy=multi-user.target
EOL


echo "building checkmate"
/usr/local/go/bin/go build .

echo "enabling systemd service"
systemctl enable checkmate.service

echo "starting systemd service"
systemctl start checkmate.service
