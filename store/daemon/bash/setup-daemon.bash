set -eo pipefail

main() {
  HOSTNAME=`hostname`
  OSD_ID=''

  until etcdctl --no-sync -C $ETCD get /deis/store/monSetupComplete >/dev/null 2>&1 ; do
    echo "store-daemon: waiting for monitor setup to complete..."
    sleep 5
  done

  if ! etcdctl --no-sync -C $ETCD get /deis/store/osds/$HOST >/dev/null 2>&1 ; then
    echo "store-daemon: creating OSD..."

    OSD_ID=`ceph osd create 2>/dev/null`

    if ! [[ "${OSD_ID}" =~ ^-?[0-9]+$ ]] ; then
      echo "store-daemon: FATAL - We have an OSD ID that isn't an integer"
      echo "store-daemon: FATAL - This likely means the monitor we tried to connect to isn't up, but others may be."
      echo "store-daemon: FATAL - We can't proceed because we don't know if an OSD was created or not."
      exit 1
    fi

    echo "store-daemon: created OSD ${OSD_ID}"
    etcdctl --no-sync -C $ETCD set /deis/store/osds/$HOST ${OSD_ID} >/dev/null
  fi

  if [ -z "${OSD_ID}" ]; then
    OSD_ID=`etcdctl --no-sync -C $ETCD get /deis/store/osds/${HOST}`
  fi

  # Make sure osd directory exists
  mkdir -p /var/lib/ceph/osd/ceph-${OSD_ID}

  # Check to see if our OSD has been initialized
  if [ ! -e /var/lib/ceph/osd/ceph-${OSD_ID}/keyring ]; then
    echo "store-daemon: OSD not yet initialized. Initializing..."
    ceph-osd -i $OSD_ID --mkfs --mkjournal --osd-journal /var/lib/ceph/osd/ceph-${OSD_ID}/journal
    ceph auth get-or-create osd.${OSD_ID} osd 'allow *' mon 'allow profile osd' -o /var/lib/ceph/osd/ceph-${OSD_ID}/keyring
    ceph osd crush add ${OSD_ID} 1.0 root=default host=${HOSTNAME}
  fi
}