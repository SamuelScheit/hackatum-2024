echo "building checkmate"
/usr/local/go/bin/go build . -tags "sqlite_math_functions"

echo "restarting systemd service"
systemctl restart checkmate.service

sleep 2
journalctl -u checkmate.service --no-pager -n 20
