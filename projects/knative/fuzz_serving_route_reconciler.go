package route

import (
	"testing"
	nativeTesting "testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	fakeservingclient "knative.dev/serving/pkg/client/injection/client/fake"
	fakerevisioninformer "knative.dev/serving/pkg/client/injection/informers/serving/v1/revision/fake"
	fakerouteinformer "knative.dev/serving/pkg/client/injection/informers/serving/v1/route/fake"
	. "knative.dev/serving/pkg/testing/v1"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
	"k8s.io/client-go/tools/cache"
)

func FuzzRouteReconciler(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		ff := fuzz.NewConsumer(data)

		route := &servingv1.Route{}
		ff.GenerateStruct(route)
		errs := route.Validate(context.Background())
		if len(errs) != 0 {
			t.Skip()
		}
		rev := &servingv1
		ff.GenerateStruct(rev)
		errs = rev.Validate(context.Background())
		if len(errs) != 0 {
			t.Skip()
		}

		newT := &nativeTesting.T{}
		ctx, _, ctl, _, cf := newTestSetup(newT)
		defer cf()

		fakeservingclient.Get(ctx).ServingV1().Revisions(testNamespace).Create(ctx, rev, metav1.CreateOptions{})
		fakerevisioninformer.Get(ctx).Informer().GetIndexer().Add(rev)

		fakeservingclient.Get(ctx).ServingV1().Routes(testNamespace).Create(ctx, route, metav1.CreateOptions{})

		fakerouteinformer.Get(ctx).Informer().GetIndexer().Add(route)

		key, err := cache.MetaNamespaceKeyFunc(route)
		if err != nil {
			t.Skip()
		}
		ctl.Reconciler.Reconcile(ctx, key)
	})
}
