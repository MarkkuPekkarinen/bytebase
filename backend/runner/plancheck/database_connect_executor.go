package plancheck

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/component/dbfactory"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	"github.com/bytebase/bytebase/backend/plugin/db"
	"github.com/bytebase/bytebase/backend/store"
)

var _ Executor = (*DatabaseConnectExecutor)(nil)

// NewDatabaseConnectExecutor creates a task check database connect executor.
func NewDatabaseConnectExecutor(store *store.Store, dbFactory *dbfactory.DBFactory) Executor {
	return &DatabaseConnectExecutor{
		store:     store,
		dbFactory: dbFactory,
	}
}

// DatabaseConnectExecutor checks if the database connection is valid.
type DatabaseConnectExecutor struct {
	store     *store.Store
	dbFactory *dbfactory.DBFactory
}

// Run runs the executor.
func (e *DatabaseConnectExecutor) Run(ctx context.Context, config *storepb.PlanCheckRunConfig) ([]*storepb.PlanCheckRunResult_Result, error) {
	instance, err := e.store.GetInstanceV2(ctx, &store.FindInstanceMessage{ResourceID: &config.InstanceId})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get instance %s", config.InstanceId)
	}
	if instance == nil {
		return nil, errors.Errorf("instance %s not found ", config.InstanceId)
	}

	database, err := e.store.GetDatabaseV2(ctx, &store.FindDatabaseMessage{InstanceID: &instance.ResourceID, DatabaseName: &config.DatabaseName})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get database %q", config.DatabaseName)
	}
	if database == nil {
		return nil, errors.Errorf("database not found %q", config.DatabaseName)
	}

	return []*storepb.PlanCheckRunResult_Result{e.checkDatabaseConnection(ctx, instance, database)}, nil
}

func (e *DatabaseConnectExecutor) checkDatabaseConnection(ctx context.Context, instance *store.InstanceMessage, database *store.DatabaseMessage) *storepb.PlanCheckRunResult_Result {
	err := func() error {
		driver, err := e.dbFactory.GetAdminDatabaseDriver(ctx, instance, database, db.ConnectionContext{})
		if err != nil {
			return err
		}
		defer driver.Close(ctx)
		return driver.Ping(ctx)
	}()
	if err != nil {
		return &storepb.PlanCheckRunResult_Result{
			Status:  storepb.PlanCheckRunResult_Result_ERROR,
			Code:    common.DBConnectionFailure.Int32(),
			Title:   fmt.Sprintf("Failed to connect %q", database.DatabaseName),
			Content: err.Error(),
		}
	}
	return &storepb.PlanCheckRunResult_Result{
		Status:  storepb.PlanCheckRunResult_Result_SUCCESS,
		Code:    common.Ok.Int32(),
		Title:   "OK",
		Content: fmt.Sprintf("Successfully connected %q", database.DatabaseName),
	}
}
