/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package v1api20200601

// TODO: it doesn't really matter where these are (as long as they're in 'apis', where is where we run controller-gen).
// These are the permissions required by the generic_controller. They're here because they can't go outside the 'apis'
// directory.

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
