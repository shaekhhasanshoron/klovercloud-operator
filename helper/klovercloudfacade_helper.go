package helper

import (
	servicev1alpha1 "bitbucket.org/klovercloud/klovercloud-operator/api/v1alpha1"
	"bitbucket.org/klovercloud/klovercloud-operator/constant"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type kc_facade_helper struct {
	labels          map[string]string
	matchLabels     map[string]string
	annotations     map[string]string
	deployment      *appsv1.Deployment
	facadeContainer corev1.Container
}

func KlovercloudFacade() *kc_facade_helper {
	return new(kc_facade_helper)
}

func (sv *kc_facade_helper) CustomResourceV1Alpha1() *servicev1alpha1.KlovercloudFacade {
	sv.labels = make(map[string]string)
	sv.labels["app"] = constant.KlovercloudFacadeDeployment

	facade := &servicev1alpha1.KlovercloudFacade{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1alpha1", Kind: constant.KlovercloudFacadeKind},
		ObjectMeta: metav1.ObjectMeta{
			Name:      constant.KlovercloudFacadeDeployment,
			Namespace: constant.KlovercloudNamespace,
			Labels:    sv.labels,
		},
		Spec:   servicev1alpha1.KlovercloudFacadeSpec{},
		Status: servicev1alpha1.KlovercloudFacadeStatus{},
	}
	return facade
}

func (sv *kc_facade_helper) Deployment() *appsv1.Deployment {
	sv.labels = make(map[string]string)
	sv.labels["app"] = constant.KlovercloudFacadeDeployment

	sv.annotations = make(map[string]string)

	sv.matchLabels = make(map[string]string)
	sv.matchLabels["app"] = constant.KlovercloudFacadeDeployment

	sv.deployment = &appsv1.Deployment{}
	sv.deployment.TypeMeta = metav1.TypeMeta{APIVersion: appsv1.SchemeGroupVersion.String(), Kind: "Deployment"}
	sv.deployment.ObjectMeta = metav1.ObjectMeta{
		Name:        constant.KlovercloudFacadeDeployment,
		Namespace:   constant.KlovercloudNamespace,
		Labels:      sv.labels,
		Annotations: sv.annotations,
	}

	sv.deployment.Spec = appsv1.DeploymentSpec{}
	sv.deployment.Spec.Replicas = Common().Int32(1)
	sv.deployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: sv.matchLabels,
	}
	sv.deployment.Spec.Template = corev1.PodTemplateSpec{}
	sv.deployment.Spec.Template.ObjectMeta = metav1.ObjectMeta{
		Labels: sv.labels,
	}
	sv.deployment.Spec.Template.Spec = corev1.PodSpec{}

	sv.facadeContainer = corev1.Container{}
	sv.facadeContainer.Name = constant.KlovercloudFacadeDeployment
	sv.facadeContainer.Image = constant.KlovercloudFacadeImage
	sv.facadeContainer.ImagePullPolicy = corev1.PullAlways
	sv.facadeContainer.Ports = []corev1.ContainerPort{
		corev1.ContainerPort{
			ContainerPort: 8080,
			Protocol:      corev1.ProtocolTCP,
		}, corev1.ContainerPort{
			ContainerPort: 8081,
			Protocol:      corev1.ProtocolTCP,
		}}

	sv.deployment.Spec.Template.Spec.Containers = []corev1.Container{sv.facadeContainer}
	return sv.deployment
}
