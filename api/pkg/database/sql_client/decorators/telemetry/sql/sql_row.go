package sql

import (
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"sync"
)

var _ contracts.SQLRow = (*sqlRowTelemetry)(nil)

type sqlRowTelemetry struct {
	delegate contracts.SQLRow
	span     tracercontracts.TracerSpan
	once     sync.Once
}

func wrapSQLRowForTelemetry(delegate contracts.SQLRow, span tracercontracts.TracerSpan) contracts.SQLRow {
	return &sqlRowTelemetry{
		delegate: delegate,
		span:     span,
	}
}

func (s *sqlRowTelemetry) Err() error {
	err := s.delegate.Err()
	s.end(err)
	return err
}

func (s *sqlRowTelemetry) Scan(dest ...any) error {
	err := s.delegate.Scan(dest...)
	s.end(err)
	return err
}

func (s *sqlRowTelemetry) end(err error) {
	s.once.Do(func() {
		if err != nil {
			s.span.RecordError(err)
		}
		s.span.End()
	})
}
