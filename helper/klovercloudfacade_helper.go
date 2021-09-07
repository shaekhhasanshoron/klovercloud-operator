package helper

import (
	servicev1alpha1 "bitbucket.org/klovercloud/klovercloud-operator/api/v1alpha1"
	"bitbucket.org/klovercloud/klovercloud-operator/constant"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type kc_facade_helper struct {
	labels          map[string]string
	matchLabels     map[string]string
	annotations     map[string]string
	deployment      *appsv1.Deployment
	svc             *corev1.Service
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

	cpuLimit, _ := resource.ParseQuantity("1000m")
	memoryLimit, _ := resource.ParseQuantity("1024Mi")
	cpuRequest, _ := resource.ParseQuantity("500m")
	memoryRequest, _ := resource.ParseQuantity("1024Mi")

	sv.facadeContainer.Resources = corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    cpuLimit,
			corev1.ResourceMemory: memoryLimit,
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    cpuRequest,
			corev1.ResourceMemory: memoryRequest,
		},
	}

	// TODO image pull policy
	sv.deployment.Spec.Template.Spec.Containers = []corev1.Container{sv.facadeContainer}
	return sv.deployment
}

func (sv *kc_facade_helper) Service() *corev1.Service {
	sv.labels = make(map[string]string)
	sv.labels["app"] = constant.KlovercloudFacadeDeployment

	sv.matchLabels = make(map[string]string)
	sv.matchLabels["app"] = constant.KlovercloudFacadeDeployment

	sv.svc = &corev1.Service{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.SchemeGroupVersion.String(), Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      constant.KlovercloudFacadeDeployment,
			Namespace: constant.KlovercloudNamespace,
			Labels:    sv.labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "http-rest", Port: 80, Protocol: "TCP", TargetPort: intstr.FromString("8080")},
				{Name: "http-metrics", Port: 8081, Protocol: "TCP", TargetPort: intstr.FromString("8081")},
			},
			Selector: sv.matchLabels,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
	return sv.svc
}
