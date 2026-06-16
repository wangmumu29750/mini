package service

import (
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"
)

func TestAdminService_ListStations(t *testing.T) {
	type args struct {
		query dto.StationListQuery
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.StationListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ListStations(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.ListStations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.ListStations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_PublicStations(t *testing.T) {
	type args struct {
		query dto.StationListQuery
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.PublicStations(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.PublicStations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.PublicStations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_CreateStation(t *testing.T) {
	type args struct {
		req dto.SaveStationRequest
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.AdminStationResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CreateStation(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.CreateStation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.CreateStation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_UpdateStation(t *testing.T) {
	type args struct {
		id  uint64
		req dto.SaveStationRequest
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.AdminStationResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.UpdateStation(tt.args.id, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.UpdateStation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.UpdateStation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_DisableStation(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.AdminStationResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.DisableStation(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.DisableStation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.DisableStation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stationFromRequest(t *testing.T) {
	type args struct {
		req dto.SaveStationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    model.Station
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stationFromRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("stationFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stationFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
