package metrics

import (
	"bytes"
	"encoding/json"
	"github.com/giantswarm/etcd-backup/config"
	"github.com/giantswarm/microerror"
	"net/http"
)

func Send(prometheusConfig *config.PrometheusConfig, metrics *BackupMetrics, tenantClusterName string) (bool, error) {
	type bodyType struct {
		Cluster                   string
		Successful                bool
		BackupSizeMeasurement     int64
		CreationTimeMeasurement   int64
		EncryptionTimeMeasurement int64
		UploadTimeMeasurement     int64
	}

	body := bodyType{
		Cluster:                   tenantClusterName,
		Successful:                metrics.Successful,
		BackupSizeMeasurement:     metrics.BackupSizeMeasurement,
		CreationTimeMeasurement:   metrics.CreationTimeMeasurement,
		EncryptionTimeMeasurement: metrics.EncryptionTimeMeasurement,
		UploadTimeMeasurement:     metrics.UploadTimeMeasurement,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return true, microerror.Mask(err)
	}

	_, err = http.Post("http://etcd-backup-metrics-collector:8080/", "application/json", bytes.NewBuffer(data))

	return true, err
}
