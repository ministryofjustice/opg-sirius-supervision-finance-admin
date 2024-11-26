package db

//type mockConn struct {
//	rows pgx.Rows
//	err  error
//}
//
//func (m *mockConn) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
//	return m.rows, m.err
//}

//func testConvertRowsToStrings(t *testing.T) {
//	tests := []struct {
//		name    string
//		rows    pgx.Rows
//		want    [][]string
//		wantErr error
//	}{
//		{
//			name:    "success",
//			rows:    pgx.Rows{},
//			want:    [][]string{},
//			wantErr: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := convertRowsToString(tt.rows)
//
//			assert.Equal(t, tt.want, got)
//			assert.Equal(t, tt.wantErr, err)
//		})
//	}
//}
