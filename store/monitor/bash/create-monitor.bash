set -eo pipefail

main() {
  # If we don't have a monitor keyring, this is a new monitor
  if [ ! -e /var/lib/ceph/mon/ceph-${HOSTNAME}/keyring ]; then
    if [ ! -f /etc/ceph/monmap ]; then
      ceph mon getmap -o /etc/ceph/monmap
    fi

    # Import the client.admin keyring and the monitor keyring into a new, temporary one
    ceph-authtool /tmp/ceph.mon.keyring --create-keyring --import-keyring /etc/ceph/ceph.client.admin.keyring
    ceph-authtool /tmp/ceph.mon.keyring --import-keyring /etc/ceph/ceph.mon.keyring

    # Make the monitor directory
    mkdir -p /var/lib/ceph/mon/ceph-${HOSTNAME}

    # Prepare the monitor daemon's directory with the map and keyring
    ceph-mon --mkfs -i ${HOSTNAME} --monmap /etc/ceph/monmap --keyring /tmp/ceph.mon.keyring

    # Clean up the temporary key
    rm /tmp/ceph.mon.keyring
  fi
}
