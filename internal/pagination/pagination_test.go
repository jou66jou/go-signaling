package pagination

import (
	"reflect"
	"testing"
)

func TestPagination_SetTotalCountAndPage(t *testing.T) {
	type args struct {
		total int32
	}
	tests := []struct {
		name               string
		p                  *Pagination
		args               args
		expectedPagination *Pagination
	}{
		{
			name:               "total 100,set page=0;perPage=33",
			p:                  &Pagination{Page: 0, PerPage: 33},
			args:               args{total: 100},
			expectedPagination: &Pagination{Page: 1, PerPage: 33, TotalCount: 100, TotalPage: 4},
		},
		{
			name:               "total 100,set page=1;perPage=0",
			p:                  &Pagination{Page: 0, PerPage: 0},
			args:               args{total: 100},
			expectedPagination: &Pagination{Page: 1, PerPage: globalDefaultPerPage, TotalCount: 100, TotalPage: 5},
		},
		{
			name:               "total 1,set page=1;perPage=4",
			p:                  &Pagination{Page: 1, PerPage: 4},
			args:               args{total: 1},
			expectedPagination: &Pagination{Page: 1, PerPage: 4, TotalCount: 1, TotalPage: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.p
			p.SetTotalCountAndPage(tt.args.total)
			if p.Page != tt.expectedPagination.Page {
				t.Errorf("SetTotalCountAndPage error, Page '%v', expected '%v'", p.Page, tt.expectedPagination.Page)
				return
			}
			if p.PerPage != tt.expectedPagination.PerPage {
				t.Errorf("SetTotalCountAndPage error, PerPage '%v', expected '%v'", p.PerPage, tt.expectedPagination.PerPage)
				return
			}
			if p.TotalCount != tt.expectedPagination.TotalCount {
				t.Errorf("SetTotalCountAndPage error, TotalCount '%v', expected '%v'", p.TotalCount, tt.expectedPagination.TotalCount)
				return
			}
			if p.TotalPage != tt.expectedPagination.TotalPage {
				t.Errorf("SetTotalCountAndPage error, TotalPage '%v', expected '%v'", p.TotalPage, tt.expectedPagination.TotalPage)
				return
			}
		})
	}
}

func TestPagination_LimitAndOffset(t *testing.T) {
	type fields struct {
		Page       int32
		PerPage    int32
		TotalCount int32
		TotalPage  int32
	}
	tests := []struct {
		name       string
		fields     fields
		wantLimit  uint32
		wantOffset uint32
	}{
		{
			name: "total=100;page=3;perPage=33",
			fields: fields{
				Page:       3,
				PerPage:    33,
				TotalCount: 100,
				TotalPage:  4,
			},
			wantLimit:  33,
			wantOffset: 66,
		},
		{
			name: "total=1;page=1;perPage=33",
			fields: fields{
				Page:       1,
				PerPage:    33,
				TotalCount: 1,
				TotalPage:  1,
			},
			wantLimit:  33,
			wantOffset: 0,
		},
		{
			name: "total=100;page=4;perPage=25",
			fields: fields{
				Page:       4,
				PerPage:    25,
				TotalCount: 100,
				TotalPage:  4,
			},
			wantLimit:  25,
			wantOffset: 75,
		},
		{
			name: "total=0;page=0;perPage=30",
			fields: fields{
				Page:       0,
				PerPage:    30,
				TotalCount: 0,
				TotalPage:  0,
			},
			wantLimit:  30,
			wantOffset: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Page:       tt.fields.Page,
				PerPage:    tt.fields.PerPage,
				TotalCount: tt.fields.TotalCount,
				TotalPage:  tt.fields.TotalPage,
			}
			got, got1 := p.LimitAndOffset()
			if got != tt.wantLimit {
				t.Errorf("Pagination.LimitAndOffset() got = %v, want %v", got, tt.wantLimit)
			}
			if got1 != tt.wantOffset {
				t.Errorf("Pagination.LimitAndOffset() got1 = %v, want %v", got1, tt.wantOffset)
			}
		})
	}
}

func TestPagination_CheckOrSetDefault(t *testing.T) {
	type fields struct {
		Page    int32
		PerPage int32
	}
	type args struct {
		params []int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Pagination
	}{
		{
			name:   "set with specified default perPage",
			fields: fields{Page: 0, PerPage: 0},
			args:   args{params: []int32{20}},
			want:   &Pagination{Page: 1, PerPage: 20},
		},
		{
			name:   "set with global default perPage",
			fields: fields{Page: 0, PerPage: 0},
			args:   args{params: []int32{}},
			want:   &Pagination{Page: 1, PerPage: globalDefaultPerPage},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Page:    tt.fields.Page,
				PerPage: tt.fields.PerPage,
			}
			if got := p.CheckOrSetDefault(tt.args.params...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pagination.CheckOrSetDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
