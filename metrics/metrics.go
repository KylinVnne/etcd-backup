package metrics

import (
	"github.com/giantswarm/etcd-backup/config"
)

func Send(prometheusConfig *config.PrometheusConfig, metrics *BackupMetrics, tenantClusterName string) (bool, error) {

	// RecordMetrics(tenantClusterName, metrics)

	return true, nil

	// prometheus URL might be empty, in that case we can't push any metric
	//if prometheusConfig.Url != "" {
	//	registry := prometheus.NewRegistry()
	//
	//	labels := prometheus.Labels{
	//		labelTenantClusterId: tenantClusterName,
	//	}
	//
	//	if metrics.Successful {
	//		// successful backup
	//		registry.MustRegister(creationTime, encryptionTime, uploadTime, backupSize, successCounter, attemptsCounter)
	//		pusher := push.New(prometheusConfig.Url, prometheusConfig.Job).Gatherer(registry)
	//
	//		creationTime.With(labels).Set(float64(metrics.CreationTimeMeasurement))
	//		encryptionTime.With(labels).Set(float64(metrics.EncryptionTimeMeasurement))
	//		uploadTime.With(labels).Set(float64(metrics.UploadTimeMeasurement))
	//		backupSize.With(labels).Set(float64(metrics.BackupSizeMeasurement))
	//		successCounter.With(labels).Inc()
	//		attemptsCounter.With(labels).Inc()
	//
	//		if err := pusher.Add(); err != nil {
	//			return true, err
	//		}
	//	} else {
	//		// failed backup
	//		registry.MustRegister(failureCounter, attemptsCounter)
	//		pusher := push.New(prometheusConfig.Url, prometheusConfig.Job).Gatherer(registry)
	//
	//		failureCounter.With(labels).Inc()
	//		attemptsCounter.With(labels).Inc()
	//
	//		if err := pusher.Add(); err != nil {
	//			return true, err
	//		}
	//	}
	//	return true, nil
	//}
	//
	//return false, nil
}
