package config

type InternalAPI struct {
	Bind                      string `env:"INTERNAL_API_GRPC_SERVER_BIND" envDefault:":11000"`
	DatasourceSnapshotAddress string `env:"INTERNAL_API_DATASOURCE_SNAPSHOT_ADDRESS" envDefault:"localhost:11100"`
	EnsResolverAddress        string `env:"INTERNAL_API_ENS_RESOLVER_ADDRESS" envDefault:":20200"`
}
