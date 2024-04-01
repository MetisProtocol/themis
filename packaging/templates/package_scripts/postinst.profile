#!/bin/bash
# This is a postinstallation script so the service can be configured and started when requested
#
sudo -u themis themisd init --chain=testnet --home /var/lib/themis
sudo systemctl daemon-reload