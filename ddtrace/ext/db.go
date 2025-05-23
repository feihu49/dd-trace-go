// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package ext

const (
	// DBApplication indicates the application using the database.
	DBApplication = "db.application"
	// DBName indicates the database name.
	DBName = "db.name"
	// DBType indicates the type of Database.
	DBType = "db.type"
	// DBInstance indicates the instance name of Database.
	DBInstance = "db.instance"
	// DBUser indicates the user name of Database, e.g. "readonly_user" or "reporting_user".
	DBUser = "db.user"
	// DBStatement records a database statement for the given database type.
	DBStatement = "db.statement"
	// DBSystem indicates the database management system (DBMS) product being used.
	DBSystem = "db.system"
)

// Available values for db.system.
const (
	DBSystemMemcached          = "memcached"
	DBSystemMySQL              = "mysql"
	DBSystemPostgreSQL         = "postgresql"
	DBSystemMicrosoftSQLServer = "mssql"
	// DBSystemOtherSQL is used for other SQL databases not listed above.
	DBSystemOtherSQL      = "other_sql"
	DBSystemElasticsearch = "elasticsearch"
	DBSystemRedis         = "redis"
	DBSystemValkey        = "valkey"
	DBSystemMongoDB       = "mongodb"
	DBSystemCassandra     = "cassandra"
	DBSystemConsulKV      = "consul"
	DBSystemLevelDB       = "leveldb"
	DBSystemBuntDB        = "buntdb"
)

// MicrosoftSQLServer tags.
const (
	// MicrosoftSQLServerInstanceName indicates the Microsoft SQL Server instance name connecting to.
	MicrosoftSQLServerInstanceName = "db.mssql.instance_name"
)

// MongoDB tags.
const (
	// MongoDBCollection indicates the collection being accessed.
	MongoDBCollection = "db.mongodb.collection"
)

// Redis tags.
const (
	// RedisDatabaseIndex indicates the Redis database index connected to.
	RedisDatabaseIndex = "db.redis.database_index"

	// RedisRawCommand allows to set the raw command for tags.
	RedisRawCommand = "redis.raw_command"

	// RedisClientCacheHit is the remaining TTL in seconds of client side cache.
	RedisClientCacheHit = "db.redis.client.cache.hit"

	// RedisClientCacheTTL captures the Time-To-Live (TTL) of a cached entry in the client.
	RedisClientCacheTTL = "db.redis.client.cache.ttl"

	// RedisClientCachePTTL is the remaining PTTL in seconds of client side cache.
	RedisClientCachePTTL = "db.redis.client.cache.pttl"

	// RedisClientCachePXAT is the remaining PXAT in seconds of client side cache.
	RedisClientCachePXAT = "db.redis.client.cache.pxat"
)

// Valkey tags.
const (
	// ValkeyRawCommand allows to set the raw command for tags.
	ValkeyRawCommand = "valkey.raw_command"

	// ValkeyClientCacheHit is the remaining TTL in seconds of client side cache.
	ValkeyClientCacheHit = "db.valkey.client.cache.hit"

	// ValkeyClientCacheTTL captures the Time-To-Live (TTL) of a cached entry in the client.
	ValkeyClientCacheTTL = "db.valkey.client.cache.ttl"

	// ValkeyClientCachePTTL is the remaining PTTL in seconds of client side cache.
	ValkeyClientCachePTTL = "db.valkey.client.cache.pttl"

	// ValkeyClientCachePXAT is the remaining PXAT in seconds of client side cache.
	ValkeyClientCachePXAT = "db.valkey.client.cache.pxat"
)

// Cassandra tags.
const (
	// CassandraConsistencyLevel is the tag name to set for consistency level.
	CassandraConsistencyLevel = "cassandra.consistency_level"

	// CassandraCluster specifies the tag name that is used to set the cluster.
	CassandraCluster = "cassandra.cluster"

	// CassandraDatacenter specifies the tag name that is used to set the datacenter.
	CassandraDatacenter = "cassandra.datacenter"

	// CassandraRowCount specifies the tag name to use when settings the row count.
	CassandraRowCount = "cassandra.row_count"

	// CassandraKeyspace is used as tag name for setting the key space.
	CassandraKeyspace = "cassandra.keyspace"

	// CassandraPaginated specifies the tag name for paginated queries.
	CassandraPaginated = "cassandra.paginated"

	// CassandraContactPoints holds the list of cassandra initial seed nodes used to discover the cluster.
	CassandraContactPoints = "db.cassandra.contact.points"

	// CassandraHostID represents the host ID for this operation.
	CassandraHostID = "db.cassandra.host.id"
)
