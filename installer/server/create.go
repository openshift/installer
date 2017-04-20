package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/dghubble/sessions"
	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/server/ctxh"
)

// CreateOperation defines a cluster creation request.
type CreateOperation struct {
	// Kind of cluster which should be created
	ClusterKind string `json:"clusterKind"`
	// Cluster properties
	ClusterData json.RawMessage `json:"cluster"`
	// If true, don't actually create cluster. Just generate assets.
	DryRun bool `json:"dryRun"`
}

// Cluster parses cluster kind and data to return a Cluster.
func (o *CreateOperation) Cluster() (Cluster, error) {
	var cluster Cluster
	switch o.ClusterKind {
	case "tectonic-metal":
		cluster = new(TectonicMetalCluster)
	case "tectonic-aws":
		cluster = new(TectonicAWSCluster)
	default:
		return nil, fmt.Errorf("installer: invalid cluster kind %s", o.ClusterKind)
	}
	err := json.Unmarshal(o.ClusterData, cluster)
	return cluster, err
}

func createHandler(sessionProvider sessions.Store) ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		if req.Method != "POST" {
			return ctxh.NewAppError(nil, "POST method required", http.StatusMethodNotAllowed)
		}

		op, err := decodeCreateOp(req)
		if err != nil {
			return ctxh.NewAppError(err, "failed to parse body", http.StatusBadRequest)
		}

		cluster, err := op.Cluster()
		if err != nil {
			return ctxh.NewAppError(err, "failed to parse cluster data", http.StatusBadRequest)
		}

		// validate cluster data and set default values
		if err := cluster.Initialize(); err != nil {
			return ctxh.NewAppError(err, err.Error(), http.StatusBadRequest)
		}

		// generate assets for the cluster
		assets, err := cluster.GenerateAssets()
		if err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to submit definition: %v", err.Error()), http.StatusBadRequest)
		}

		// push cluster data to a remote provisioning service
		if op.DryRun {
			log.Info("Dry run requested. Provisioner will not be contacted.")
		} else {
			err := cluster.Publish(ctx)
			if err != nil {
				return ctxh.NewAppError(err, fmt.Sprintf("Error pushing cluster data to provisioner: %v", err.Error()), http.StatusBadRequest)
			}
		}

		// store the cluster kind and StatusChecker in the session
		checker, err := cluster.StatusChecker()
		if err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to create a status checker: %v", err.Error()), http.StatusInternalServerError)
		}
		session := sessionProvider.New(installerSessionName)
		session.Values["kind"] = cluster.Kind()
		session.Values["checker"] = checker
		if err := sessionProvider.Save(w, session); err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to save session: %v", err.Error()), http.StatusInternalServerError)
		}

		// write assets.zip in response
		if assets != nil {
			zipAssetHandler(assets).ServeHTTP(ctx, w, req)
		}
		return nil
	}
	return ctxh.ContextHandlerFuncWithError(fn)
}

func decodeCreateOp(req *http.Request) (*CreateOperation, error) {
	op := new(CreateOperation)
	err := json.NewDecoder(req.Body).Decode(op)
	return op, err
}
