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

package v1

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	azuretnnovaiov1 "tutorial.kubebuilder.io/project/api/v1"
)

// nolint:unused
// log is for logging in this package.
var playerlog = logf.Log.WithName("player-resource")

// SetupPlayerWebhookWithManager registers the webhook for Player in the manager.
func SetupPlayerWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&azuretnnovaiov1.Player{}).
		WithValidator(&PlayerCustomValidator{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-azure-tnnova-io-azure-tnnova-io-v1-player,mutating=false,failurePolicy=fail,sideEffects=None,groups=azure.tnnova.io.azure.tnnova.io,resources=players,verbs=create;update,versions=v1,name=vplayer-v1.kb.io,admissionReviewVersions=v1

// PlayerCustomValidator struct is responsible for validating the Player resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type PlayerCustomValidator struct {
	//TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &PlayerCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Player.
func (v *PlayerCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	player, ok := obj.(*azuretnnovaiov1.Player)
	if !ok {
		return nil, fmt.Errorf("expected a Player object but got %T", obj)
	}
	playerlog.Info("Validation for Player upon creation", "name", player.GetName())

	return nil, validatePlayer(player)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Player.
func (v *PlayerCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	player, ok := newObj.(*azuretnnovaiov1.Player)
	if !ok {
		return nil, fmt.Errorf("expected a Player object for the newObj but got %T", newObj)
	}
	playerlog.Info("Validation for Player upon update", "name", player.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, validatePlayer(player)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Player.
func (v *PlayerCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	player, ok := obj.(*azuretnnovaiov1.Player)
	if !ok {
		return nil, fmt.Errorf("expected a Player object but got %T", obj)
	}
	playerlog.Info("Validation for Player upon deletion", "name", player.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

func validatePlayer(player *azuretnnovaiov1.Player) error {
	var allErrs field.ErrorList
	if err := validatePlayerName(player); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "azure.tnnova.io", Kind: "Player"},
		player.Name, allErrs)
}

func validatePlayerName(player *azuretnnovaiov1.Player) *field.Error {
	if len(player.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		return field.Invalid(field.NewPath("metadata").Child("name"), player.ObjectMeta.Name, "must be no more than 52 characters")
	}

	if player.ObjectMeta.Name != player.Spec.Name {
		return field.Invalid(field.NewPath("spec").Child("name"), player.Spec.Name, "Name must be equal to metadata name")
	}
	return nil
}

//func validatePlayer(player *azuretnnovaiov1.Player) error {
//	// nolint:unused
//	return nil, nil
//}
