package main

import (
	"github.com/robfig/cron"

	"github.com/deis/deis/database/bindata"

	"github.com/deis/deis/pkg/boot"
	"github.com/deis/deis/pkg/commons"
	"github.com/deis/deis/pkg/logger"
)

func main() {
	externalPort := commons.Getopt("EXTERNAL_PORT", "5432")
	etcdPath := commons.Getopt("ETCD_PATH", "/deis/database")
	adminUser := commons.Getopt("PG_ADMIN_USER", "postgres")
	adminPass := commons.Getopt("PG_ADMIN_PASS", "changeme123")
	user := commons.Getopt("PG_USER_NAME", "deis")
	password := commons.Getopt("PG_USER_PASS", "changeme123")
	name := commons.Getopt("PG_USER_DB", "deis")
	bucketName := commons.Getopt("BUCKET_NAME", "db_wal")
	backupsToRetain := commons.Getopt("BACKUPS_TO_RETAIN", "5")
	backupFrequency := commons.Getopt("BACKUP_FREQUENCY", "3h")
	pgConfig := commons.Getopt("PG_CONFIG", "/etc/postgresql/9.3/main/postgresql.conf")
	listenAddress := commons.Getopt("PG_LISTEN", "*")

	bootProcess := boot.New("tcp", etcdPath, externalPort)
	logger.Log.Debug("creating required defaults in etcd...")

	commons.MkdirEtcd(bootProcess.Etcd, etcdPath)
	commons.SetDefaultEtcd(bootProcess.Etcd, etcdPath+"/engine", "postgresql_psycopg2")
	commons.SetDefaultEtcd(bootProcess.Etcd, etcdPath+"/adminUser", adminUser)
	commons.SetDefaultEtcd(bootProcess.Etcd, etcdPath+"/adminPass", adminPass)
	commons.SetDefaultEtcd(bootProcess.Etcd, etcdPath+"/user", user)
	commons.SetDefaultEtcd(bootProcess.Etcd, etcdPath+"/password", password)
	commons.SetDefaultEtcd(bootProcess.Etcd, etcdPath+"/name", name)
	commons.SetDefaultEtcd(bootProcess.Etcd, etcdPath+"/bucketName", bucketName)

	logger.Log.Info("starting deis-database...")

	bootProcess.Start()

	bootProcess.RunBashScript("bash/postgres-init.bash", nil, bindata.Asset)

	postgresCommand := "sudo -i -u postgres /usr/lib/postgresql/9.3/bin/postgres" +
		" -c config-file=" + pgConfig +
		" -c listen-addresses=" + listenAddress
	bootProcess.StartProcessAsChild(commons.BuildCommandFromString(postgresCommand))

	bootProcess.WaitForLocalConnection("5432")

	bootProcess.Publish()

	params := make(map[string]string)
	params["BUCKET_NAME"] = bucketName
	bootProcess.RunBashScript("bash/postgres.bash", params, bindata.Asset)

	// schedule periodic backups using wal-e
	scheduleBackup := cron.New()
	scheduleBackup.AddFunc("@every "+backupFrequency,
		func() {
			logger.Log.Debug("creating database backup with wal-e...")
			params := make(map[string]string)
			params["BACKUPS_TO_RETAIN"] = backupsToRetain
			bootProcess.RunBashScript("bash/backup.bash", params, bindata.Asset)
		})
	scheduleBackup.Start()

	logger.Log.Info("database: postgres is running...")
	bootProcess.Wait()
}
