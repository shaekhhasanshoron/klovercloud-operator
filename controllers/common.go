package controllers

import (
	"bitbucket.org/klovercloud/klovercloud-operator/constant"
	"context"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

func InitiateFacadeController(mgr manager.Manager) error {
	facadeController := &KlovercloudFacadeReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Err:     nil,
		Context: context.Background(),
	}

	if err := facadeController.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", constant.KlovercloudFacadeKind)
		return err
	}

	go facadeController.Init()
	return nil
}

func InitiateManagementController(mgr manager.Manager) error {
	managementController := &KlovercloudManagementReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName(constant.KlovercloudManagementController),
	}

	if err := managementController.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", constant.KlovercloudManagementKind)
		return err
	}
	return nil
}

func InitiatePipelineController(mgr manager.Manager) error {
	pipelineController := &KlovercloudPipelineReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName(constant.KlovercloudPipelineController),
	}

	if err := pipelineController.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", constant.KlovercloudPipelineKind)
		return err
	}
	return nil
}
