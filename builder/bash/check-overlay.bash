set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  # remove any pre-existing docker.sock
  test -e /var/run/docker.sock && rm -f /var/run/docker.sock

  # create empty env file
  touch /etc/docker.env

  # force overlay if it's supported
  mkdir --parents --mode=0700 /
  fstype=$(findmnt --noheadings --output FSTYPE --target /)
  if [[ "$fstype" == "overlay" ]]; then
    echo 'DRIVER_OVERRIDE="--storage-driver=overlay"' > /etc/docker.env
  fi
}

