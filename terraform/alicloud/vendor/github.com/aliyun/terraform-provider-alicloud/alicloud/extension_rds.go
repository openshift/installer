package alicloud

type Engine string

const (
	MySQL      = Engine("MySQL")
	SQLServer  = Engine("SQLServer")
	PPAS       = Engine("PPAS")
	PostgreSQL = Engine("PostgreSQL")
)

var WEEK_ENUM = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

var BACKUP_TIME = []string{
	"00:00Z-01:00Z", "01:00Z-02:00Z", "02:00Z-03:00Z", "03:00Z-04:00Z", "04:00Z-05:00Z",
	"05:00Z-06:00Z", "06:00Z-07:00Z", "07:00Z-08:00Z", "08:00Z-09:00Z", "09:00Z-10:00Z",
	"10:00Z-11:00Z", "11:00Z-12:00Z", "12:00Z-13:00Z", "13:00Z-14:00Z", "14:00Z-15:00Z",
	"15:00Z-16:00Z", "16:00Z-17:00Z", "17:00Z-18:00Z", "18:00Z-19:00Z", "19:00Z-20:00Z",
	"20:00Z-21:00Z", "21:00Z-22:00Z", "22:00Z-23:00Z", "23:00Z-24:00Z",
}

type KVStoreInstanceType string

const (
	KVStoreRedis    = KVStoreInstanceType("Redis")
	KVStoreMemcache = KVStoreInstanceType("Memcache")
)

type KVStoreEngineVersion string

const (
	KVStore2Dot8 = KVStoreEngineVersion("2.8")
	KVStore4Dot0 = KVStoreEngineVersion("4.0")
)
