package inject

import (
	"errors"
	"flag"
)

const (
	flagInjectDefault               = "inject-default"
	flagEnableIAMForServiceAccounts = "enable-iam-for-service-accounts"
	flagAWSRegion                   = "aws-region"
	flagEnvoyPreview                = "preview"
	flagLogLevel                    = "log-level"
	flagECRSecret                   = "ecr-secret"
	flagSidecarImage                = "sidecar-image"
	flagSidecarCpuRequests          = "sidecar-cpu-requests"
	flagSidecarMemoryRequests       = "sidecar-memory-requests"
	flagInitImage                   = "init-image"
	flagIgnoredIPs                  = "ignored-ips"
	flagEnableJaegerTracing         = "enable-jaeger-tracing"
	flagJaegerAddress               = "jaeger-address"
	flagJaegerPort                  = "jaeger-port"
	flagEnableDatadogTracing        = "enable-datadog-tracing"
	flagDatadogAddress              = "datadog-address"
	flagDatadogPort                 = "datadog-port"
	flagEnableXrayTracing           = "enable-xray-tracing"
	flagEnableStatsTags             = "enable-stats-tags"
	flagEnableStatsD                = "enable-statsd"
)

type Config struct {
	// Injection Settings
	InjectDefault bool

	// If enabled, an fsGroup: 1337 will be injected in the absence of it within pod securityContext
	// see https://github.com/aws/amazon-eks-pod-identity-webhook/issues/8 for more details
	EnableIAMForServiceAccounts bool

	// Sidecar settings
	SidecarImage  string
	SidecarCpu    string
	SidecarMemory string
	Region        string
	Preview       bool
	LogLevel      string
	EcrSecret     bool

	// Init container settings
	InitImage  string
	IgnoredIPs string

	// Observability settings
	EnableXrayTracing    bool
	EnableStatsTags      bool
	EnableStatsD         bool
	EnableJaegerTracing  bool
	JaegerAddress        string
	JaegerPort           string
	EnableDatadogTracing bool
	DatadogAddress       string
	DatadogPort          string
}

// MultipleTracer checks if more than one tracer is configured.
func MultipleTracer(config *Config) bool {
	j := config.EnableJaegerTracing
	d := config.EnableDatadogTracing
	x := config.EnableXrayTracing

	return (j && d) || (d && x) || (j && x)
}

func (cfg *Config) BindFlags() {
	flag.BoolVar(&cfg.InjectDefault, flagInjectDefault, true,
		`If enabled, sidecars will be injected in the absence of the corresponding pod annotation`)
	flag.BoolVar(&cfg.EnableIAMForServiceAccounts, flagEnableIAMForServiceAccounts, true,
		`If enabled, an fsGroup: 1337 will be injected in the absence of it within pod securityContext`)
	flag.StringVar(&cfg.Region, flagAWSRegion, "",
		"AWS App Mesh region")
	flag.BoolVar(&cfg.Preview, flagEnvoyPreview, false,
		"Enable preview channel")
	flag.StringVar(&cfg.LogLevel, flagLogLevel, "info",
		"AWS App Mesh envoy log level")
	flag.BoolVar(&cfg.EcrSecret, flagECRSecret, false,
		"Inject AWS app mesh pull secrets")
	flag.StringVar(&cfg.SidecarImage, flagSidecarImage, "840364872350.dkr.ecr.us-west-2.amazonaws.com/aws-appmesh-envoy:v1.12.3.0-prod",
		"Envoy sidecar container image.")
	flag.StringVar(&cfg.SidecarCpu, flagSidecarCpuRequests, "10m",
		"Envoy sidecar CPU resources requests.")
	flag.StringVar(&cfg.SidecarMemory, flagSidecarMemoryRequests, "32Mi",
		"Envoy sidecar memory resources requests.")
	flag.StringVar(&cfg.InitImage, flagInitImage, "111345817488.dkr.ecr.us-west-2.amazonaws.com/aws-appmesh-proxy-route-manager:v2",
		"Init container image.")
	flag.StringVar(&cfg.IgnoredIPs, flagIgnoredIPs, "169.254.169.254",
		"Init container ignored IPs.")
	flag.BoolVar(&cfg.EnableJaegerTracing, flagEnableJaegerTracing, false,
		"Enable Envoy Jaeger tracing")
	flag.StringVar(&cfg.JaegerAddress, flagJaegerAddress, "appmesh-jaeger.appmesh-system",
		"Jaeger address")
	flag.StringVar(&cfg.JaegerPort, flagJaegerPort, "9411",
		"Jaeger port")
	flag.BoolVar(&cfg.EnableDatadogTracing, flagEnableDatadogTracing, false,
		"Enable Envoy Datadog tracing")
	flag.StringVar(&cfg.DatadogAddress, flagDatadogAddress, "datadog.appmesh-system",
		"Datadog Agent address")
	flag.StringVar(&cfg.DatadogPort, flagDatadogPort, "8126",
		"Datadog Agent tracing port")
	flag.BoolVar(&cfg.EnableXrayTracing, flagEnableXrayTracing, false,
		"Enable Envoy X-Ray tracing integration and injects xray-daemon as sidecar")
	flag.BoolVar(&cfg.EnableStatsTags, flagEnableStatsTags, false,
		"Enable Envoy to tag stats")
	flag.BoolVar(&cfg.EnableStatsD, flagEnableStatsD, false,
		"If enabled, Envoy will send DogStatsD metrics to 127.0.0.1:8125")
}

func (cfg *Config) BindEnv() error {
	return nil
}

func (cfg *Config) Validate() error {
	if MultipleTracer(cfg) {
		return errors.New("Envoy only supports a single tracer instance. Please choose between Jaeger, Datadog or X-Ray.")
	}
	return nil
}