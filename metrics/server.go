package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
)

type Server struct {
	Logger micrologger.Logger
}

const labelTenantClusterId = "tenant_cluster_id"

var (
	labels = []string{
		labelTenantClusterId,
	}
	namespace    = "etcd_backup"
	creationTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: prometheus.BuildFQName(namespace, "", "creation_time_ms"),
		Help: "Gauge about the time in ms spent by the ETCD backup creation process.",
	}, labels)
	encryptionTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: prometheus.BuildFQName(namespace, "", "encryption_time_ms"),
		Help: "Gauge about the time in ms spent by the ETCD backup encryption process.",
	}, labels)
	uploadTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: prometheus.BuildFQName(namespace, "", "upload_time_ms"),
		Help: "Gauge about the time in ms spent by the ETCD backup upload process.",
	}, labels)
	backupSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: prometheus.BuildFQName(namespace, "", "size_bytes"),
		Help: "Gauge about the size of the backup file, as seen by S3.",
	}, labels)
	attemptsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "attempts_count"),
		Help: "Count of attempted backups",
	}, labels)
	successCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "success_count"),
		Help: "Count of successful backups",
	}, labels)
	failureCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: prometheus.BuildFQName(namespace, "", "failure_count"),
		Help: "Count of failed backups",
	}, labels)
)

func New(logger micrologger.Logger) *Server {
	return &Server{
		Logger: logger,
	}
}

func recordMetrics(tenantClusterName string, metrics *BackupMetrics) {
	labels := prometheus.Labels{
		labelTenantClusterId: tenantClusterName,
	}

	attemptsCounter.With(labels).Inc()

	if metrics.Successful {
		creationTime.With(labels).Set(float64(metrics.CreationTimeMeasurement))
		encryptionTime.With(labels).Set(float64(metrics.EncryptionTimeMeasurement))
		uploadTime.With(labels).Set(float64(metrics.UploadTimeMeasurement))
		backupSize.With(labels).Set(float64(metrics.BackupSizeMeasurement))
		successCounter.With(labels).Inc()
	} else {
		failureCounter.With(labels).Inc()
	}
}

func (ms Server) Listen() error {
	go func() {
		ms.listenHttp()
	}()
	return ms.listenPrometheus()
}

func (ms Server) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("qua\n")
	if r.Method == "GET" && r.RequestURI == "/healthz" {
		fmt.Fprint(w, "OK")
		return
	}

	body, _ := ioutil.ReadAll(r.Body)

	var metrics BackupMetrics

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&metrics)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		ms.Logger.Log("msg", err)
		return
	}

	var info struct{ Cluster string }

	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&info)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		ms.Logger.Log("msg", err)
		return
	}

	recordMetrics(info.Cluster, &metrics)

	fmt.Fprint(w, "OK")
}

func (ms Server) listenHttp() error {
	port := 8080
	ms.Logger.Log("level", "info", "msg", fmt.Sprintf("Starting metrics update listener on port %d", port))
	http.HandleFunc("/", ms.handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		ms.Logger.Log("level", "error", "msg", fmt.Sprintf("Error listening on port %d", port))
		return microerror.Mask(err)
	}

	return nil
}

func (ms Server) listenPrometheus() error {
	port := 2112
	ms.Logger.Log("level", "info", "msg", fmt.Sprintf("Starting prometheus metrics listener on port %d", port))
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		ms.Logger.Log("level", "error", "msg", fmt.Sprintf("Error listening on port %d", port))
		return microerror.Mask(err)
	}

	return nil
}
