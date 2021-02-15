package config

import (
	"flag"
	"github.com/golang/glog"
)

const (
	MaxSlotsPerLevel uint = 65535
	DefaultSlotsPerLevel uint = 100
	VehicleTypeCar = "car"
	DefaultCache = "local"
	DefaultMQDriver = "rabbitmq"
)

var (
	ClusterMode bool

	DatabaseDriver string
	DatabaseHost string
	DatabaseUser string
	DatabasePassword string

	CacheType string = DefaultCache

	MQDriver string
	MQHost string
	QueueName string

	MlcpLevelCount uint16
	SlotsPerLevel uint
	VechileTypes []string = []string{VehicleTypeCar}
)

func init() {
	flag.BoolVar(&ClusterMode, "clusterMode", false, "Run MLCP server in cluster mode")
	flag.StringVar(&DatabaseDriver, "databaseDriver", "mariadb", "database driver to use")
	flag.StringVar(&DatabaseHost, "databaseHost", "", "database host")
	flag.StringVar(&DatabaseUser, "databaseUser", "", "database user")
	flag.StringVar(&DatabasePassword, "databasePassword", "", "database password")
	flag.UintVar(&SlotsPerLevel, "slotsPerLevel", DefaultSlotsPerLevel, "number of parking slots per level")
	flag.StringVar(&MQDriver, "mqDriver", DefaultMQDriver, "message queue driver")
	flag.StringVar(&MQHost, "mqurl", "", "messageQueue url")
	flag.StringVar(&QueueName, "queueName", "mlcp", "queue name")

	if SlotsPerLevel > MaxSlotsPerLevel {
		glog.Warningf("Slot count set per level exceeds maximum slots supported; defaulting to %d", DefaultSlotsPerLevel)
		SlotsPerLevel = DefaultSlotsPerLevel
	}
}
func Init () {
	flag.Set("logtostderr", "true")
	flag.Parse()
}
