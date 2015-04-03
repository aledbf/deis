set -eo pipefail

# set debug based on envvar
[[ $DEBUG ]] && set -x

main() {
  # HACK: load progrium/cedarish tarball for faster boot times
  # see https://github.com/deis/deis/issues/1027
  if ! docker history progrium/cedarish >/dev/null 2>/dev/null ; then
      echo "Loading cedarish..."
      cat /progrium_cedarish.tar.gz | docker import - progrium/cedarish
  else
      echo "Cedarish already loaded"
  fi

  # build required images
  docker build -t deis/slugbuilder /app/slugbuilder/
  docker build -t deis/slugrunner /app/slugrunner/
}

