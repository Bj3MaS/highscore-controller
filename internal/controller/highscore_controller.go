/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"slices"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	azuretnnovaiov1 "tutorial.kubebuilder.io/project/api/v1"
)

// HighscoreReconciler reconciles a Highscore object
type HighscoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=azure.tnnova.io.azure.tnnova.io,resources=highscores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=azure.tnnova.io.azure.tnnova.io,resources=highscores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=azure.tnnova.io.azure.tnnova.io,resources=highscores/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Highscore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *HighscoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info(req.Name)

	var player azuretnnovaiov1.Player

	if err := r.Get(ctx, req.NamespacedName, &player); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)

	}

	var highscore azuretnnovaiov1.Highscore

	if err := r.Get(ctx, client.ObjectKey{Name: "dx-olympics"}, &highscore); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to fetch Highscore")
			return ctrl.Result{}, err
		}

		highscore = azuretnnovaiov1.Highscore{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dx-olympics",
				Namespace: req.Namespace,
			},
			Spec: azuretnnovaiov1.HighscoreSpec{
				Scoreboard: []azuretnnovaiov1.Score{
					{
						Name:   player.Spec.Name,
						Points: player.Spec.Points,
					},
				},
			},
		}

		if err := r.Create(ctx, &highscore); err != nil {
			log.Error(err, "unable to create Highscore")
			return ctrl.Result{}, err
		}
	}

	const playerFinalizer = "highscore.example.com/finalizer"

	if player.ObjectMeta.DeletionTimestamp.IsZero() && !controllerutil.ContainsFinalizer(&player, playerFinalizer) {
		// Update only when necessary to avoid loops
		controllerutil.AddFinalizer(&player, playerFinalizer)
		if err := r.Update(ctx, &player); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	log.Info("before deleetion")
	// Handle resource deletion
	if !player.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("Delete player")
		log.Info(player.Name)
		// Perform cleanup before removing the finalizer
		if err := r.finalizeplayer(&player, &highscore); err != nil {
			return ctrl.Result{}, err
		}

		log.Info(fmt.Sprintf("%v", highscore.Spec.Scoreboard))

		// Remove finalizer to allow deletion to proceed
		controllerutil.RemoveFinalizer(&player, playerFinalizer)
		if err := r.Update(ctx, &player); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.Update(ctx, &highscore); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// var highscore azuretnnovaiov1.Highscore

	// if err := r.Get(ctx, client.ObjectKey{Name: "dx-olympics"}, &highscore); err != nil {
	// 	if client.IgnoreNotFound(err) != nil {
	// 		log.Error(err, "unable to fetch Highscore")
	// 		return ctrl.Result{}, err
	// 	}

	// 	highscore = azuretnnovaiov1.Highscore{
	// 		ObjectMeta: metav1.ObjectMeta{
	// 			Name:      "dx-olympics",
	// 			Namespace: req.Namespace,
	// 		},
	// 		Spec: azuretnnovaiov1.HighscoreSpec{
	// 			Scoreboard: []azuretnnovaiov1.Score{
	// 				{
	// 					Name:   player.Spec.Name,
	// 					Points: player.Spec.Points,
	// 				},
	// 			},
	// 		},
	// 	}

	// 	if err := r.Create(ctx, &highscore); err != nil {
	// 		log.Error(err, "unable to create Highscore")
	// 		return ctrl.Result{}, err
	// 	}
	// }

	if containsName(highscore, player.Spec.Name) {
		log.Info("Updating player")

		index := slices.IndexFunc(highscore.Spec.Scoreboard, func(score azuretnnovaiov1.Score) bool {
			return score.Name == player.Spec.Name
		})

		highscore.Spec.Scoreboard[index].Points = player.Spec.Points

	} else {
		highscore.Spec.Scoreboard = append(highscore.Spec.Scoreboard, azuretnnovaiov1.Score{
			Name:   player.Spec.Name,
			Points: player.Spec.Points,
		})
	}

	newLeader := setLeader(highscore)

	highscore.Spec.Leader = newLeader
	log.Info("updsting")
	log.Info(fmt.Sprintf("%v", highscore.Spec.Scoreboard))
	if err := r.Update(ctx, &highscore); err != nil {
		log.Error(err, "unable to update Highscore")
	}

	highscore.Status.Message = fmt.Sprintf("Player %s has taken the lead", newLeader)

	if err := r.Status().Update(ctx, &highscore); err != nil {
		log.Error(err, "Failed to update Highscore status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func containsName(highscore azuretnnovaiov1.Highscore, name string) bool {
	for _, item := range highscore.Spec.Scoreboard {
		if item.Name == name {
			return true
		}
	}
	return false
}

func setLeader(highscore azuretnnovaiov1.Highscore) string {
	var leader string
	var max int
	for _, item := range highscore.Spec.Scoreboard {
		if item.Points > max {
			max = item.Points
			leader = item.Name
		}
	}
	return leader
}

// SetupWithManager sets up the controller with the Manager.
func (r *HighscoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&azuretnnovaiov1.Player{}).
		Named("highscore").
		Complete(r)
}

func (r *HighscoreReconciler) finalizeplayer(player *azuretnnovaiov1.Player, highscore *azuretnnovaiov1.Highscore) error {
	index := slices.IndexFunc(highscore.Spec.Scoreboard, func(score azuretnnovaiov1.Score) bool {
		return score.Name == player.Spec.Name
	})

	highscore.Spec.Scoreboard = append(highscore.Spec.Scoreboard[:index], highscore.Spec.Scoreboard[index+1:]...)
	//
	// delete any external resources associated with the cronJob
	//
	// Ensure that delete implementation is idempotent and safe to invoke
	// multiple times for same object.
	return nil
}
