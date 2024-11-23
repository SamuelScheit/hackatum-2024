echo "building checkmate"
/usr/local/go/bin/go build .

echo "restarging systemd service"
systemctl restart checkmate.service

sleep 2
journalctl -u checkmate.service --no-pager
