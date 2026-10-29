package scheduler

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestComputePodResourceRequestCountsRestartableInitContainers(t *testing.T) {
	always := v1.ContainerRestartPolicyAlways
	pod := &v1.Pod{Spec: v1.PodSpec{
		InitContainers: []v1.Container{
			{
				Name:          "sidecar-a",
				RestartPolicy: &always,
				Resources: v1.ResourceRequirements{Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("500m"),
					v1.ResourceMemory: resource.MustParse("64Mi"),
				}},
			},
			{
				Name: "setup",
				Resources: v1.ResourceRequirements{Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("1"),
					v1.ResourceMemory: resource.MustParse("256Mi"),
				}},
			},
			{
				Name:          "sidecar-b",
				RestartPolicy: &always,
				Resources: v1.ResourceRequirements{Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("250m"),
					v1.ResourceMemory: resource.MustParse("128Mi"),
				}},
			},
		},
		Containers: []v1.Container{{
			Name: "app",
			Resources: v1.ResourceRequirements{Requests: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("2"),
				v1.ResourceMemory: resource.MustParse("64Mi"),
			}},
		}},
		Overhead: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("50m"),
			v1.ResourceMemory: resource.MustParse("16Mi"),
		},
	}}

	got := ComputePodResourceRequest(pod)
	if got.MilliCPU != 2800 {
		t.Errorf("CPU request = %dm, want 2800m", got.MilliCPU)
	}
	wantMemoryQuantity := resource.MustParse("336Mi")
	wantMemory := wantMemoryQuantity.Value()
	if got.Memory != wantMemory {
		t.Errorf("memory request = %d bytes, want %d bytes", got.Memory, wantMemory)
	}
}
