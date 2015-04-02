set -eo pipefail

main() {
  # set the number of placement groups for the default pools - they come up with defaults that are too low
  if ! etcdctl --no-sync -C $ETCD get /deis/store/defaultPoolsConfigured >/dev/null 2>&1 ; then
    echo "store-gateway: setting pg_num values for default pools..."
    function set_until_success {
      set +e

      echo "store-gateway: checking pool $1..."
      if ! ceph osd pool get $1 pg_num | grep "pg_num: $2" ; then
        ceph osd pool set $1 pg_num $2 2>/dev/null
        PG_SET=$?
        until [[ $PG_SET -eq 0 ]]; do
          sleep 5
          ceph osd pool set $1 pg_num $2 2>/dev/null
          PG_SET=$?
        done
      fi

      if ! ceph osd pool get $1 pgp_num | grep "pgp_num: $2" ; then
        ceph osd pool set $1 pgp_num $2 2>/dev/null
        PGP_SET=$?
        until [[ $PGP_SET -eq 0 ]]; do
          sleep 5
          ceph osd pool set $1 pgp_num $2 2>/dev/null
          PGP_SET=$?
        done
      fi

      set -e
    }

    PG_NUM=`etcdctl --no-sync -C $ETCD get /deis/store/pgNum`

    set_until_success data ${PG_NUM}
    set_until_success rbd ${PG_NUM}
    set_until_success metadata ${PG_NUM}

    etcdctl --no-sync -C $ETCD set /deis/store/defaultPoolsConfigured youBetcha >/dev/null
  fi

  # we generate a key for the gateway. we can do this because we have the client key templated out
  if ! etcdctl --no-sync -C $ETCD get /deis/store/gatewayKeyring >/dev/null 2>&1 ; then
    ceph-authtool --create-keyring /etc/ceph/ceph.client.radosgw.keyring
    chmod +r /etc/ceph/ceph.client.radosgw.keyring
    ceph-authtool /etc/ceph/ceph.client.radosgw.keyring -n client.radosgw.gateway --gen-key
    ceph-authtool -n client.radosgw.gateway --cap osd 'allow rwx' --cap mon 'allow rwx' /etc/ceph/ceph.client.radosgw.keyring
    ceph -k /etc/ceph/ceph.client.admin.keyring auth add client.radosgw.gateway -i /etc/ceph/ceph.client.radosgw.keyring
    etcdctl --no-sync -C $ETCD set /deis/store/gatewayKeyring < /etc/ceph/ceph.client.radosgw.keyring >/dev/null
  else
    etcdctl --no-sync -C $ETCD get /deis/store/gatewayKeyring > /etc/ceph/ceph.client.radosgw.keyring
    chmod +r /etc/ceph/ceph.client.radosgw.keyring
  fi

  if ! radosgw-admin user info --uid=deis >/dev/null 2>&1 ; then
    radosgw-admin user create --uid=deis --display-name="Deis" >/dev/null
  fi

  radosgw-admin user info --uid=deis >/etc/ceph/user.json
  # store the access key and secret key for consumption by other services
  ACCESS_KEY=`cat /etc/ceph/user.json | python -c 'import json,sys;obj=json.load(sys.stdin);print json.dumps(obj["keys"][0]["access_key"]);' | tr -d '"'`
  SECRET_KEY=`cat /etc/ceph/user.json | python -c 'import json,sys;obj=json.load(sys.stdin);print json.dumps(obj["keys"][0]["secret_key"]);' | tr -d '"'`
  etcdctl --no-sync -C $ETCD set $ETCD_PATH/accessKey ${ACCESS_KEY} >/dev/null
  etcdctl --no-sync -C $ETCD set $ETCD_PATH/secretKey ${SECRET_KEY} >/dev/null
}
