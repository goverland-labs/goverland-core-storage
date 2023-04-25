package dao

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// todo: add units for converting from and to internal models

func TestUnitCompare(t *testing.T) {
	for name, tc := range map[string]struct {
		d1       Dao
		d2       Dao
		expected bool
	}{
		"equal": {
			d1:       Dao{ID: "id", Name: "name", Strategies: Strategies{{Name: "name1"}, {Name: "name2"}}},
			d2:       Dao{ID: "id", Name: "name", Strategies: Strategies{{Name: "name1"}, {Name: "name2"}}},
			expected: true,
		},
		"different ID": {
			d1:       Dao{ID: "id-1", Name: "name", Strategies: Strategies{{Name: "name1"}, {Name: "name2"}}},
			d2:       Dao{ID: "id-2", Name: "name", Strategies: Strategies{{Name: "name1"}, {Name: "name2"}}},
			expected: false,
		},
		"different created at": {
			d1:       Dao{ID: "id-1", CreatedAt: time.Now()},
			d2:       Dao{ID: "id-1", CreatedAt: time.Now().Add(time.Second)},
			expected: true,
		},
		"different updated at": {
			d1:       Dao{ID: "id-1", UpdatedAt: time.Now()},
			d2:       Dao{ID: "id-1", UpdatedAt: time.Now().Add(time.Second)},
			expected: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, compare(tc.d1, tc.d2))
		})
	}
}

func TestUnitHandleDao(t *testing.T) {
	for name, tc := range map[string]struct {
		dp       func(ctrl *gomock.Controller) DataProvider
		p        func(ctrl *gomock.Controller) Publisher
		event    Dao
		expected error
	}{
		"correct creating": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, gorm.ErrRecordNotFound)
				m.EXPECT().Create(gomock.Any()).Times(1).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				return m
			},
			event:    Dao{ID: "id-1"},
			expected: nil,
		},
		"correct updating": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Dao{ID: "id-1", Name: "updated"}, nil)
				m.EXPECT().Update(gomock.Any()).Times(1).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				return m
			},
			event:    Dao{ID: "id-1", Name: "name"},
			expected: nil,
		},
		"do not update for equal objects": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Dao{ID: "id-1"}, nil)
				m.EXPECT().Update(gomock.Any()).Times(0).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				return NewMockPublisher(ctrl)
			},
			event:    Dao{ID: "id-1"},
			expected: nil,
		},
		"raise err on problems with reading from DB": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("unexpected error"))
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				return NewMockPublisher(ctrl)
			},
			event:    Dao{ID: "id-1"},
			expected: errors.New("unexpected error"),
		},
		"raise err on problems with creating in DB": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, gorm.ErrRecordNotFound)
				m.EXPECT().Create(gomock.Any()).Times(1).Return(errors.New("unexpected error"))
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				return NewMockPublisher(ctrl)
			},
			event:    Dao{ID: "id-1"},
			expected: errors.New("unexpected error"),
		},
		"allow do not send event after creating": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, gorm.ErrRecordNotFound)
				m.EXPECT().Create(gomock.Any()).Times(1).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(errors.New("unexpected error"))
				return m
			},
			event:    Dao{ID: "id-1"},
			expected: nil,
		},
		"raise err on problems with updating in DB": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Dao{ID: "id-1", Name: "name"}, nil)
				m.EXPECT().Update(gomock.Any()).Times(1).Return(errors.New("unexpected error"))
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				return NewMockPublisher(ctrl)
			},
			event:    Dao{ID: "id-1"},
			expected: errors.New("unexpected error"),
		},
		"allow do not send event after updating": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Dao{ID: "id-1", Name: "name"}, nil)
				m.EXPECT().Update(gomock.Any()).Times(1).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(errors.New("unexpected error"))
				return m
			},
			event:    Dao{ID: "id-1"},
			expected: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s, err := NewService(tc.dp(ctrl), tc.p(ctrl))
			require.Nil(t, err)

			err = s.HandleDao(context.Background(), tc.event)
			if tc.expected == nil {
				require.Nil(t, err)
				return
			}

			require.ErrorContains(t, err, tc.expected.Error())
		})
	}
}
