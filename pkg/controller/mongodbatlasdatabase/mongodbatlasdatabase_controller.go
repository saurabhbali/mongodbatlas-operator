package mongodbatlasdatabase

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"reflect"

	knappekv1alpha1 "github.com/saurabhbali/mongodbatlas-operator/pkg/apis/knappek/v1alpha1"
	//corev1 "k8s.io/api/core/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	//"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log = logf.Log.WithName("controller_mongodbatlasdatabase")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

 // ReconciliationConfig let us customize reconcilitation
type ReconciliationConfig struct {
	Time time.Duration
}

// GetReconcilitationConfig gives us default values
func GetReconcilitationConfig() *ReconciliationConfig {
	// default reconciliation loop time is 2 minutes
	timeString := "120"
	timeInt, _ := strconv.Atoi(timeString)
	reconciliationTime := time.Second * time.Duration(timeInt)
	return &ReconciliationConfig{
		Time: reconciliationTime,
	}
}

// Add creates a new MongoDBAtlasDatabase Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMongoDBAtlasDatabase{client: mgr.GetClient(), scheme: mgr.GetScheme(), reconciliationConfig: GetReconcilitationConfig()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mongodbatlasdatabase-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MongoDBAtlasDatabase
	err = c.Watch(&source.Kind{Type: &knappekv1alpha1.MongoDBAtlasDatabase{}}, &handler.EnqueueRequestForObject{})
// err = c.Watch(&source.Kind{Type: &MongoDBAtlasDatabase{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileMongoDBAtlasDatabase implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMongoDBAtlasDatabase{}

// ReconcileMongoDBAtlasDatabase reconciles a MongoDBAtlasDatabase object
type ReconcileMongoDBAtlasDatabase struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	reconciliationConfig *ReconciliationConfig
}

// Reconcile reads that state of the cluster for a MongoDBAtlasDatabase object and makes changes based on the state read
// and what is in the MongoDBAtlasDatabase.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMongoDBAtlasDatabase) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MongoDBAtlasDatabase")

	// Fetch the MongoDBAtlasDatabase instance
	instance := &knappekv1alpha1.MongoDBAtlasDatabase{}
	// instance := &MongoDBAtlasDatabase{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Create a new MongoDBAtlasDatabase
	isMongoDBAtlasDatabaseToBeCreated := reflect.DeepEqual(instance.Status, knappekv1alpha1.MongoDBAtlasDatabaseStatus{})
	if isMongoDBAtlasDatabaseToBeCreated {
		// err = createMongoDBAtlasDatabaseUser(reqLogger, r.atlasClient, atlasDatabaseUser, atlasProject)
		err = createMongoDBAtlasDatabase(reqLogger, instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	// // Check if this Pod already exists
	// found := &corev1.Pod{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	// 	err = r.client.Create(context.TODO(), pod)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Pod created successfully - don't requeue
	// 	return reconcile.Result{}, nil
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }
	// // Pod already exists - don't requeue
	// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{RequeueAfter: r.reconciliationConfig.Time}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
// func newPodForCR(cr *knappekv1alpha1.MongoDBAtlasDatabase) *corev1.Pod {
// 	labels := map[string]string{
// 		"app": cr.Name,
// 	}
// 	return &corev1.Pod{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      cr.Name + "-pod",
// 			Namespace: cr.Namespace,
// 			Labels:    labels,
// 		},
// 		Spec: corev1.PodSpec{
// 			Containers: []corev1.Container{
// 				{
// 					Name:    "busybox",
// 					Image:   "busybox",
// 					Command: []string{"sleep", "3600"},
// 				},
// 			},
// 		},
// 	}
// }

//func createMongoDBAtlasDatabase(reqLogger logr.Logger, cr *knappekv1alpha1.MongoDBAtlasDatabase) error {
// func createMongoDBAtlasDatabase(reqLogger logr.Logger, cr *MongoDBAtlasDatabase) error {
func createMongoDBAtlasDatabase(reqLogger logr.Logger, cr *knappekv1alpha1.MongoDBAtlasDatabase) error {
	hosta := cr.Spec.Host
	clientdb, err := mongo.NewClient(options.Client().ApplyURI(hosta))
	if err != nil {
		return fmt.Errorf("Error creating new client %s", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = clientdb.Connect(ctx)
	if err != nil {
		return fmt.Errorf("Error connecting Atlas cluster %s", err)
	}
	defer clientdb.Disconnect(ctx)
	quickstartDatabase := clientdb.Database(cr.Spec.Database)
	coll := quickstartDatabase.Collection("test")
	podcastResult, err := coll.InsertOne(ctx, bson.D{
		{Key: "Sample", Value: "Sample Collection"},
	})
        if err != nil {
                return fmt.Errorf("Error creating sample collection %s", err)
        }
	_ = podcastResult
	reqLogger.Info("Database created.")
	return updateCRStatus(reqLogger, cr)
}

func updateCRStatus(reqLogger logr.Logger, cr *knappekv1alpha1.MongoDBAtlasDatabase) error {
	// update status field in CR
	cr.Status.Created = "created"
	cr.Status.Database = cr.Spec.Database
	return nil
}
