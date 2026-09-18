package main

import (
	"log"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	syncv1alpha1 "demo/internal/operator/api/v1alpha1"
	"demo/internal/operator/controller"
	"demo/internal/tmdbsync"
)

func main() {
	tmdbSyncURL := os.Getenv("TMDB_SYNC_URL")
	if tmdbSyncURL == "" {
		tmdbSyncURL = "http://localhost:8082"
	}

	scheme := runtime.NewScheme()
	if err := syncv1alpha1.AddToScheme(scheme); err != nil {
		log.Fatalf("enregistrement du scheme: %v", err)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: ":8081",
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
	})
	if err != nil {
		log.Fatalf("création du manager: %v", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Fatalf("healthz check: %v", err)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		log.Fatalf("readyz check: %v", err)
	}

	reconciler := &controller.FilmSyncReconciler{
		Client:     mgr.GetClient(),
		SyncClient: tmdbsync.NewClient(tmdbSyncURL),
	}
	if err := reconciler.SetupWithManager(mgr); err != nil {
		log.Fatalf("setup du controller: %v", err)
	}

	log.Println("démarrage du manager filmsync-operator")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatalf("erreur manager: %v", err)
	}
}
