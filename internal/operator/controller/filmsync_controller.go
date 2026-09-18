package controller

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	syncv1alpha1 "demo/internal/operator/api/v1alpha1"
	"demo/internal/tmdbsync"
)

// FilmSyncReconciler réconcilie l'état réel (film synchronisé en DB) avec l'état
// désiré décrit par une ressource FilmSync.
type FilmSyncReconciler struct {
	client.Client
	SyncClient *tmdbsync.Client
}

func (r *FilmSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var fs syncv1alpha1.FilmSync
	if err := r.Get(ctx, req.NamespacedName, &fs); err != nil {
		// La ressource a été supprimée entre-temps : rien à faire (pas de finalizer en v1).
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Rien à refaire si le spec n'a pas changé depuis le dernier reconcile réussi.
	if fs.Status.ObservedGeneration == fs.Generation && fs.Status.Phase != "" {
		return ctrl.Result{}, nil
	}

	film, err := r.SyncClient.Sync(fs.Spec.TMDBID)
	if err != nil {
		fs.Status.Phase = syncv1alpha1.PhaseFailed
		fs.Status.Message = err.Error()
		fs.Status.ObservedGeneration = fs.Generation
		if statusErr := r.Status().Update(ctx, &fs); statusErr != nil {
			logger.Error(statusErr, "échec de la mise à jour du status après un échec de sync")
		}
		// Erreur retournée => controller-runtime requeue automatiquement avec backoff exponentiel.
		return ctrl.Result{}, fmt.Errorf("sync film tmdb_id=%d: %w", fs.Spec.TMDBID, err)
	}

	now := metav1.Now()
	fs.Status.Phase = syncv1alpha1.PhaseSynced
	fs.Status.FilmID = &film.ID
	fs.Status.Titre = film.Titre
	fs.Status.Message = ""
	fs.Status.LastSyncTime = &now
	fs.Status.ObservedGeneration = fs.Generation

	if err := r.Status().Update(ctx, &fs); err != nil {
		return ctrl.Result{}, fmt.Errorf("mise à jour status: %w", err)
	}

	logger.Info("film synchronisé", "tmdbID", fs.Spec.TMDBID, "filmID", film.ID)
	return ctrl.Result{}, nil
}

func (r *FilmSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&syncv1alpha1.FilmSync{}).
		Complete(r)
}
