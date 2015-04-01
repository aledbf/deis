set -eo pipefail

main() {
  echo '#!/bin/sh' > /usr/local/bin/reload
  chmod 0755 /usr/local/bin/reload

  touch /etc/postgresql/9.3/main/pg_hba.conf
  touch /etc/postgresql/9.3/main/postgresql.conf
}
