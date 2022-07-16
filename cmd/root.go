package cmd

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/numary/ledger/cmd/internal"
	"github.com/numary/ledger/pkg/redis"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DebugFlag                           = "debug"
	StorageDirFlag                      = "storage.dir"
	StorageDriverFlag                   = "storage.driver"
	StorageSQLiteDBNameFlag             = "storage.sqlite.db_name"
	StoragePostgresConnectionStringFlag = "storage.postgres.conn_string"
	StorageCacheFlag                    = "storage.cache"
	ServerHttpBindAddressFlag           = "server.http.bind_address"
	UiHttpBindAddressFlag               = "ui.http.bind_address"
	// Deprecated

	LockStrategyFlag                     = "lock-strategy"
	LockStrategyRedisUrlFlag             = "lock-strategy-redis-url"
	LockStrategyRedisDurationFlag        = "lock-strategy-redis-duration"
	LockStrategyRedisRetryFlag           = "lock-strategy-redis-retry"
	LockStrategyRedisTLSEnabledFlag      = "lock-strategy-redis-tls-enabled"
	LockStrategyRedisTLSInsecureFlag     = "lock-strategy-redis-tls-insecure"
	OtelTracesFlag                       = "otel-traces"
	OtelTracesBatchFlag                  = "otel-traces-batch"
	OtelTracesExporterFlag               = "otel-traces-exporter"
	OtelTracesExporterJaegerEndpointFlag = "otel-traces-exporter-jaeger-endpoint"
	OtelTracesExporterJaegerUserFlag     = "otel-traces-exporter-jaeger-user"
	OtelTracesExporterJaegerPasswordFlag = "otel-traces-exporter-jaeger-password"
	OtelTracesExporterOTLPModeFlag       = "otel-traces-exporter-otlp-mode"
	OtelTracesExporterOTLPEndpointFlag   = "otel-traces-exporter-otlp-endpoint"
	OtelTracesExporterOTLPInsecureFlag   = "otel-traces-exporter-otlp-insecure"
	OtelMetricsFlag                      = "otel-metrics"
	OtelMetricsExporterFlag              = "otel-metrics-exporter"
	OtelMetricsExporterOTLPModeFlag      = "otel-metrics-exporter-otlp-mode"
	OtelMetricsExporterOTLPEndpointFlag  = "otel-metrics-exporter-otlp-endpoint"
	OtelMetricsExporterOTLPInsecureFlag  = "otel-metrics-exporter-otlp-insecure"
	PublisherKafkaEnabledFlag            = "publisher-kafka-enabled"
	PublisherKafkaBrokerFlag             = "publisher-kafka-broker"
	PublisherKafkaSASLEnabled            = "publisher-kafka-sasl-enabled"
	PublisherKafkaSASLUsername           = "publisher-kafka-sasl-username"
	PublisherKafkaSASLPassword           = "publisher-kafka-sasl-password"
	PublisherKafkaSASLMechanism          = "publisher-kafka-sasl-mechanism"
	PublisherKafkaSASLScramSHASize       = "publisher-kafka-sasl-scram-sha-size"
	PublisherKafkaTLSEnabled             = "publisher-kafka-tls-enabled"
	PublisherTopicMappingFlag            = "publisher-topic-mapping"
	PublisherHttpEnabledFlag             = "publisher-http-enabled"

	AuthBearerEnabledFlag           = "auth-bearer-enabled"
	AuthBearerIntrospectUrlFlag     = "auth-bearer-introspect-url"
	AuthBearerAudienceFlag          = "auth-bearer-audience"
	AuthBearerAudiencesWildcardFlag = "auth-bearer-audiences-wildcard"
	AuthBearerUseScopesFlag         = "auth-bearer-use-scopes"

	SegmentEnabledFlag       = "segment-enabled"
	SegmentWriteKey          = "segment-write-key"
	SegmentApplicationId     = "segment-application-id"
	SegmentHeartbeatInterval = "segment-heartbeat-interval"
)

var (
	Version                = "develop"
	BuildDate              = "-"
	Commit                 = "-"
	DefaultSegmentWriteKey = ""
)

func NewRootCommand() *cobra.Command {
	viper.SetDefault("version", Version)

	root := &cobra.Command{
		Use:               "numary",
		Short:             "Numary",
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := os.MkdirAll(viper.GetString(StorageDirFlag), 0700)
			if err != nil {
				return errors.Wrap(err, "creating storage directory")
			}
			return nil
		},
	}

	server := NewServer()
	version := NewVersion()
	start := NewServerStart()
	server.AddCommand(start)

	conf := NewConfig()
	conf.AddCommand(NewConfigInit())
	store := NewStorage()
	store.AddCommand(NewStorageInit())
	store.AddCommand(NewStorageList())
	store.AddCommand(NewStorageUpgrade())
	store.AddCommand(NewStorageScan())
	store.AddCommand(NewStorageDelete())

	scriptExec := NewScriptExec()
	scriptCheck := NewScriptCheck()

	root.AddCommand(server)
	root.AddCommand(conf)
	root.AddCommand(UICmd)
	root.AddCommand(store)
	root.AddCommand(scriptExec)
	root.AddCommand(scriptCheck)
	root.AddCommand(version)
	root.AddCommand(stickersCmd)

	root.AddCommand(NewDocCommand())

	home, err := os.UserHomeDir()
	if err != nil {
		home = "/root"
	}

	root.PersistentFlags().Bool(DebugFlag, false, "Debug mode")
	root.PersistentFlags().String(StorageDriverFlag, "sqlite", "Storage driver")
	root.PersistentFlags().String(StorageDirFlag, path.Join(home, ".numary/data"), "Storage directory (for sqlite)")
	root.PersistentFlags().String(StorageSQLiteDBNameFlag, "numary", "SQLite database name")
	root.PersistentFlags().String(StoragePostgresConnectionStringFlag, "postgresql://localhost/postgres", "Postgre connection string")
	root.PersistentFlags().Bool(StorageCacheFlag, true, "Storage cache")
	root.PersistentFlags().String(ServerHttpBindAddressFlag, "localhost:3068", "API bind address")
	root.PersistentFlags().String(UiHttpBindAddressFlag, "localhost:3068", "UI bind address")
	root.PersistentFlags().Bool(OtelTracesFlag, false, "Enable OpenTelemetry traces support")
	root.PersistentFlags().Bool(OtelTracesBatchFlag, false, "Use OpenTelemetry batching")
	root.PersistentFlags().String(OtelTracesExporterFlag, "stdout", "OpenTelemetry traces exporter")
	root.PersistentFlags().String(OtelTracesExporterJaegerEndpointFlag, "", "OpenTelemetry traces Jaeger exporter endpoint")
	root.PersistentFlags().String(OtelTracesExporterJaegerUserFlag, "", "OpenTelemetry traces Jaeger exporter user")
	root.PersistentFlags().String(OtelTracesExporterJaegerPasswordFlag, "", "OpenTelemetry traces Jaeger exporter password")
	root.PersistentFlags().String(OtelTracesExporterOTLPModeFlag, "grpc", "OpenTelemetry traces OTLP exporter mode (grpc|http)")
	root.PersistentFlags().String(OtelTracesExporterOTLPEndpointFlag, "", "OpenTelemetry traces grpc endpoint")
	root.PersistentFlags().Bool(OtelTracesExporterOTLPInsecureFlag, false, "OpenTelemetry traces grpc insecure")
	root.PersistentFlags().Bool(OtelMetricsFlag, false, "Enable OpenTelemetry metrics support")
	root.PersistentFlags().String(OtelMetricsExporterFlag, "stdout", "OpenTelemetry metrics exporter")
	root.PersistentFlags().String(OtelMetricsExporterOTLPModeFlag, "grpc", "OpenTelemetry metrics OTLP exporter mode (grpc|http)")
	root.PersistentFlags().String(OtelMetricsExporterOTLPEndpointFlag, "", "OpenTelemetry metrics grpc endpoint")
	root.PersistentFlags().Bool(OtelMetricsExporterOTLPInsecureFlag, false, "OpenTelemetry metrics grpc insecure")
	root.PersistentFlags().String(LockStrategyFlag, "memory", "Lock strategy (memory, none, redis)")
	root.PersistentFlags().String(LockStrategyRedisUrlFlag, "", "Redis url when using redis locking strategy")
	root.PersistentFlags().Duration(LockStrategyRedisDurationFlag, redis.DefaultLockDuration, "Lock duration")
	root.PersistentFlags().Duration(LockStrategyRedisRetryFlag, redis.DefaultRetryInterval, "Retry lock period")
	root.PersistentFlags().Bool(LockStrategyRedisTLSEnabledFlag, false, "Use tls on redis")
	root.PersistentFlags().Bool(LockStrategyRedisTLSInsecureFlag, false, "Whether redis is trusted or not")
	root.PersistentFlags().Bool(PublisherKafkaEnabledFlag, false, "Publish write events to kafka")
	root.PersistentFlags().StringSlice(PublisherKafkaBrokerFlag, []string{}, "Kafka address is kafka enabled")
	root.PersistentFlags().StringSlice(PublisherTopicMappingFlag, []string{}, "Define mapping between internal event types and topics")
	root.PersistentFlags().Bool(PublisherHttpEnabledFlag, false, "Sent write event to http endpoint")
	root.PersistentFlags().Bool(PublisherKafkaSASLEnabled, false, "Enable SASL authentication on kafka publisher")
	root.PersistentFlags().String(PublisherKafkaSASLUsername, "", "SASL username")
	root.PersistentFlags().String(PublisherKafkaSASLPassword, "", "SASL password")
	root.PersistentFlags().String(PublisherKafkaSASLMechanism, "", "SASL authentication mechanism")
	root.PersistentFlags().Int(PublisherKafkaSASLScramSHASize, 512, "SASL SCRAM SHA size")
	root.PersistentFlags().Bool(PublisherKafkaTLSEnabled, false, "Enable TLS to connect on kafka")
	root.PersistentFlags().Bool(AuthBearerEnabledFlag, false, "Enable bearer auth")
	root.PersistentFlags().String(AuthBearerIntrospectUrlFlag, "", "OAuth2 introspect URL")
	root.PersistentFlags().StringSlice(AuthBearerAudienceFlag, []string{}, "Allowed audiences")
	root.PersistentFlags().Bool(AuthBearerAudiencesWildcardFlag, false, "Don't check audience")
	root.PersistentFlags().Bool(AuthBearerUseScopesFlag, false, "Use scopes as defined by rfc https://datatracker.ietf.org/doc/html/rfc8693")
	root.PersistentFlags().Bool(SegmentEnabledFlag, true, "Is segment enabled")
	root.PersistentFlags().String(SegmentApplicationId, "", "Segment application id")
	root.PersistentFlags().String(SegmentWriteKey, DefaultSegmentWriteKey, "Segment write key")
	root.PersistentFlags().Duration(SegmentHeartbeatInterval, 24*time.Hour, "Segment heartbeat interval")

	internal.InitHTTPBasicFlags(root)

	if err = viper.BindPFlags(root.PersistentFlags()); err != nil {
		panic(err)
	}

	viper.SetConfigName("numary")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.numary")
	viper.AddConfigPath("/etc/numary")
	if err = viper.ReadInConfig(); err != nil {
		// fmt.Printf("loading config file: %s\n", err)
	}

	internal.BindEnv(viper.GetViper())

	return root
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}
