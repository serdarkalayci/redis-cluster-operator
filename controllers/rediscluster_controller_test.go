/*
Copyright 2022.

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

package controllers

import (
	"context"
	cachev1alpha1 "github.com/containersolutions/redis-cluster-operator/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestRedisClusterReconciler_Reconcile_ReturnsIfRedisClusterIsNotFound(t *testing.T) {
	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(cachev1alpha1.GroupVersion)
	clientBuilder := fake.ClientBuilder{}
	// Create a ReconcileMemcached object with the scheme and fake client.
	r := &RedisClusterReconciler{
		Client: clientBuilder.Build(),
		Scheme: s,
	}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "redis-cluster",
			Namespace: "default",
		},
	}
	_, err := r.Reconcile(context.TODO(), req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
}

func TestRedisClusterReconciler_Reconcile_ReturnsErrorIfCannotGetStatefulset(t *testing.T) {
	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(cachev1alpha1.GroupVersion)
	clientBuilder := fake.NewClientBuilder()
	client := clientBuilder.Build()

	// Create a ReconcileMemcached object with the scheme and fake client.
	r := &RedisClusterReconciler{
		Client: client,
		Scheme: s,
	}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "redis-cluster",
			Namespace: "default",
		},
	}
	_, err := r.Reconcile(context.TODO(), req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
}

func TestRedisClusterReconciler_Reconcile_CreatesStatefulsetIfDoesntExist(t *testing.T) {
	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	_ = cachev1alpha1.AddToScheme(s)

	clientBuilder := fake.NewClientBuilder()
	clientBuilder.WithObjects(&cachev1alpha1.RedisCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis-cluster",
			Namespace: "default",
		},
	})
	client := clientBuilder.Build()

	// Create a ReconcileMemcached object with the scheme and fake client.
	r := &RedisClusterReconciler{
		Client: client,
		Scheme: s,
	}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "redis-cluster",
			Namespace: "default",
		},
	}
	_, err := r.Reconcile(context.TODO(), req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	sts := &v1.StatefulSet{}
	err = client.Get(context.TODO(), types.NamespacedName{
		Name:      "redis-cluster",
		Namespace: "default",
	}, sts)
	if err != nil {
		t.Fatalf("Failed to fetch created Statefulset %v", err)
	}
}

func TestRedisClusterReconciler_Reconcile_DoesNotFailIfStatefulsetExists(t *testing.T) {
	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	_ = cachev1alpha1.AddToScheme(s)

	clientBuilder := fake.NewClientBuilder()
	clientBuilder.WithObjects(&cachev1alpha1.RedisCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis-cluster",
			Namespace: "default",
		},
	}, &v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis-cluster",
			Namespace: "default",
		},
	})
	client := clientBuilder.Build()

	// Create a ReconcileMemcached object with the scheme and fake client.
	r := &RedisClusterReconciler{
		Client: client,
		Scheme: s,
	}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "redis-cluster",
			Namespace: "default",
		},
	}
	_, err := r.Reconcile(context.TODO(), req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	sts := &v1.StatefulSet{}
	err = client.Get(context.TODO(), types.NamespacedName{
		Name:      "redis-cluster",
		Namespace: "default",
	}, sts)
	if err != nil {
		t.Fatalf("Failed to fetch created Statefulset %v", err)
	}
}
