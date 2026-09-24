package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// GroupVersion identifie le groupe/version de l'API : sync.letterboxd.dev/v1alpha1.
var (
	GroupVersion  = schema.GroupVersion{Group: "sync.letterboxd.dev", Version: "v1alpha1"}
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
	AddToScheme   = SchemeBuilder.AddToScheme
)

func init() {
	SchemeBuilder.Register(&FilmSync{}, &FilmSyncList{})
}

// FilmSyncSpec décrit l'état désiré : quel film TMDB doit être synchronisé.
type FilmSyncSpec struct {
	TMDBID int `json:"tmdbID"`
}

// FilmSyncPhase résume l'état de synchronisation observé par le controller.
type FilmSyncPhase string

const (
	PhasePending FilmSyncPhase = "Pending"
	PhaseSynced  FilmSyncPhase = "Synced"
	PhaseFailed  FilmSyncPhase = "Failed"
)

// ConditionTypeSynced est le seul type de condition exposé pour l'instant :
// True = dernier sync réussi, False = dernier sync en échec.
const ConditionTypeSynced = "Synced"

// FilmSyncStatus décrit l'état observé, mis à jour uniquement par le controller.
type FilmSyncStatus struct {
	Phase              FilmSyncPhase      `json:"phase,omitempty"`
	FilmID             *uint              `json:"filmID,omitempty"`
	Titre              string             `json:"titre,omitempty"`
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

// FilmSync est la ressource custom représentant la synchronisation d'un film TMDB.
type FilmSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FilmSyncSpec   `json:"spec,omitempty"`
	Status FilmSyncStatus `json:"status,omitempty"`
}

// FilmSyncList contient une liste de FilmSync (nécessaire pour "kubectl get filmsyncs").
type FilmSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FilmSync `json:"items"`
}

// DeepCopyObject implémente runtime.Object : le client K8s (informers/cache) en a besoin
// pour cloner les objets qu'il met en cache, sans partager d'état muable entre deux reconciles.
// Normalement généré par controller-gen ; on l'écrit ici à la main pour bien voir ce qu'il fait.
func (in *FilmSync) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(FilmSync)
	in.DeepCopyInto(out)
	return out
}

func (in *FilmSync) DeepCopyInto(out *FilmSync) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *FilmSyncStatus) DeepCopyInto(out *FilmSyncStatus) {
	*out = *in
	if in.FilmID != nil {
		out.FilmID = new(uint)
		*out.FilmID = *in.FilmID
	}
	if in.Conditions != nil {
		out.Conditions = make([]metav1.Condition, len(in.Conditions))
		copy(out.Conditions, in.Conditions)
	}
}

func (in *FilmSyncList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(FilmSyncList)
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]FilmSync, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return out
}
