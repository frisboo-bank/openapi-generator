package sqlx

import (
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"sync"
)

var _ contracts.SQLXRows = (*sqlxRowsTelemetry)(nil)

type sqlxRowsTelemetry struct {
	contracts.SQLXRows
	span tracercontracts.TracerSpan
	once sync.Once
}

func wrapSQLXRowsForTelemetry(delegate contracts.SQLXRows, span tracercontracts.TracerSpan) contracts.SQLXRows {
	return &sqlxRowsTelemetry{
		SQLXRows: delegate,
		span:     span,
	}
}

func (s *sqlxRowsTelemetry) Close() error {
	err := s.SQLXRows.Close()
	s.end(err)
	return err
}

func (r *sqlxRowsTelemetry) end(err error) {
	r.once.Do(func() {
		if err != nil {
			r.span.RecordError(err)
		}
		r.span.End()
	})
}
