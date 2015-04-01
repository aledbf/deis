set -eo pipefail

main() {
  if etcdctl --no-sync -C $ETCD mk ${ETCD_PATH}/masterLock $HOSTNAME --ttl $ETCD_TTL >/dev/null 2>&1 \
  || [[ `etcdctl --no-sync -C $ETCD get ${ETCD_PATH}/masterLock` == "$HOSTNAME" ]] ; then
    etcdctl --no-sync -C $ETCD set $ETCD_PATH/host $HOST --ttl $ETCD_TTL >/dev/null
    etcdctl --no-sync -C $ETCD set $ETCD_PATH/port $EXTERNAL_PORT --ttl $ETCD_TTL >/dev/null
    etcdctl --no-sync -C $ETCD update ${ETCD_PATH}/masterLock $HOSTNAME --ttl $ETCD_TTL >/dev/null
  fi
}
