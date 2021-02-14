package webproject

import (
	wp "github.com/chaunceyt/webproject-operator/pkg/apis/wp/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	kubectlContainer = "lachlanevenson/k8s-kubectl:v1.20.2"
)

// backupCronJob - responsible for backing up the database to a remote storage solution.
func (r *ReconcileWebproject) backupCronJob(cr *wp.WebProject) *v1beta1.CronJob {
	// TODO:
	// - Add configmap that contains the script to backup the database
	backupCommand := `
	echo 'Starting DB Backup'  &&  \
		mysqlshow -h$DATABASE_HOST -u$MYSQL_USER -p$MYSQL_PASSWORD && \
		mysqldump -h$DATABASE_HOST --opt $MYSQL_DATABASE > /var/lib/mysql/database-backup-drupal_db.sql -uroot -p$MYSQL_ROOT_PASSWORD && \
		cd /var/lib/mysql/ && \
		gzip database-backup-drupal_db.sql && \
		ls -ltr /var/lib/mysql/
`
	cron := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "backup-cron"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "backup"),
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: cr.Spec.DatabaseSidecar.Backup.BackupSchedule,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Volumes: []corev1.Volume{
								{
									Name: "database-backup-storage",
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{},
									},
								},
							},
							Containers: []corev1.Container{
								{
									Name:    cr.Name,
									Image:   cr.Spec.DatabaseSidecar.DatabaseImage,
									Command: []string{"/bin/sh", "-c"},
									Args:    []string{backupCommand},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "database-backup-storage",
											MountPath: "/var/lib/mysql",
										},
									},
									Env: []corev1.EnvVar{
										{
											Name:  "MYSQL_ROOT_PASSWORD",
											Value: cr.Spec.DatabaseSidecar.DatabaseRootPassword,
										},
										{
											Name:  "MYSQL_DATABASE",
											Value: cr.Spec.DatabaseSidecar.DatabaseName,
										},
										{
											Name:  "MYSQL_USER",
											Value: cr.Spec.DatabaseSidecar.DatabaseUser,
										},
										{
											Name:  "MYSQL_PASSWORD",
											Value: cr.Spec.DatabaseSidecar.DatabaseUserPassword,
										},
										{
											Name:  "DATABASE_HOST",
											Value: workloadName(cr, "backup-svc") + "." + cr.Namespace,
										},
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, cron, r.scheme)
	return cron
}

// databaseContainerCronJob - cronjob that will execute a script in webcontainer if enabled
func (r *ReconcileWebproject) databaseContainerCronJob(cr *wp.WebProject) *v1beta1.CronJob {
	runCommand := sidecarRunCommand(cr, "database")
	cron := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "database-cron"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "cron"),
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: cr.Spec.DatabaseSidecar.CronJob.Schedule,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							ServiceAccountName: "cronjob-executor",
							Containers: []corev1.Container{
								{
									Name:    cr.Name,
									Image:   kubectlContainer,
									Command: []string{"/bin/sh", "-c"},
									Args:    []string{runCommand},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, cron, r.scheme)
	return cron
}

// webContainerCronJob - cronjob that will execute a script in webcontainer if enabled
func (r *ReconcileWebproject) webContainerCronJob(cr *wp.WebProject) *v1beta1.CronJob {
	runCommand := sidecarRunCommand(cr, "web")
	cron := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "web-cron"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "cron"),
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: cr.Spec.WebContainer.CronJob.Schedule,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							ServiceAccountName: "cronjob-executor",
							Containers: []corev1.Container{
								{
									Name:    cr.Name,
									Image:   kubectlContainer,
									Command: []string{"/bin/sh", "-c"},
									Args:    []string{runCommand},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, cron, r.scheme)
	return cron
}

// searchContainerCronJob - cronjob that will execute a script in searchcontainer if enabled
func (r *ReconcileWebproject) searchContainerCronJob(cr *wp.WebProject) *v1beta1.CronJob {
	runCommand := sidecarRunCommand(cr, "search")
	cron := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "search-cron"),
			Namespace: cr.Namespace,
			Labels:    webprojectlabels(cr, "cron"),
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: "0 * * * *",
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							ServiceAccountName: "cronjob-executor",
							Containers: []corev1.Container{
								{
									Name:    cr.Name,
									Image:   kubectlContainer,
									Command: []string{"/bin/sh", "-c"},
									Args:    []string{runCommand},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(cr, cron, r.scheme)
	return cron
}

func sidecarRunCommand(cr *wp.WebProject, container string) string {
	runCommand := "export podName=$(kubectl get po -l app.kubernetes.io/name=" + cr.Name + " -o name | cut -f2 -d/) && kubectl exec $podName -c " + container + " -- /opt/script/script.sh"
	return runCommand
}
