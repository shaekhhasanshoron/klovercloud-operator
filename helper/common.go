package helper

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type common struct{}

func Common() *common {
	return new(common)
}

func (sv *common) CheckSubstrings(str string, subs ...string) bool {
	fmt.Printf("String: \"%s\", Substrings: %s\n", str, subs)

	for _, sub := range subs {
		if strings.Contains(str, sub) == false {
			return false
		}
	}
	return true
}

func (sv *common) KlovercloudNamespace() *corev1.Namespace {
	labels := make(map[string]string)
	labels["role"] = "klovercloud"
	ns := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.SchemeGroupVersion.String(), Kind: "Namespace"},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "klovercloud",
			Labels: labels,
		},
	}
	return ns
}

func (sv *common) Int32(v int32) *int32 {
	return &v
}
