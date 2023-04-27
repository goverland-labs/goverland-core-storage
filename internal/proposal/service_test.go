package proposal

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
		p1       Proposal
		p2       Proposal
		expected bool
	}{
		"equal": {
			p1:       Proposal{ID: "id", Title: "title", Strategies: Strategies{{Name: "name1"}, {Name: "name2"}}},
			p2:       Proposal{ID: "id", Title: "title", Strategies: Strategies{{Name: "name1"}, {Name: "name2"}}},
			expected: true,
		},
		"different ID": {
			p1:       Proposal{ID: "id-1", Title: "title", Strategies: Strategies{{Name: "name1"}, {Name: "name2"}}},
			p2:       Proposal{ID: "id-2", Title: "title", Strategies: Strategies{{Name: "name1"}, {Name: "name2"}}},
			expected: false,
		},
		"different created at": {
			p1:       Proposal{ID: "id-1", CreatedAt: time.Now()},
			p2:       Proposal{ID: "id-1", CreatedAt: time.Now().Add(time.Second)},
			expected: true,
		},
		"different updated at": {
			p1:       Proposal{ID: "id-1", UpdatedAt: time.Now()},
			p2:       Proposal{ID: "id-1", UpdatedAt: time.Now().Add(time.Second)},
			expected: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, compare(tc.p1, tc.p2))
		})
	}
}

func TestUnitHandleProposal(t *testing.T) {
	for name, tc := range map[string]struct {
		dp       func(ctrl *gomock.Controller) DataProvider
		p        func(ctrl *gomock.Controller) Publisher
		event    Proposal
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
			event:    Proposal{ID: "id-1"},
			expected: nil,
		},
		"correct updating": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Proposal{ID: "id-1", Title: "updated"}, nil)
				m.EXPECT().Update(gomock.Any()).Times(1).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				return m
			},
			event:    Proposal{ID: "id-1", Title: "name"},
			expected: nil,
		},
		"do not update for equal objects": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Proposal{ID: "id-1"}, nil)
				m.EXPECT().Update(gomock.Any()).Times(0).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				return NewMockPublisher(ctrl)
			},
			event:    Proposal{ID: "id-1"},
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
			event:    Proposal{ID: "id-1"},
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
			event:    Proposal{ID: "id-1"},
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
			event:    Proposal{ID: "id-1"},
			expected: nil,
		},
		"raise err on problems with updating in DB": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Proposal{ID: "id-1", Title: "name"}, nil)
				m.EXPECT().Update(gomock.Any()).Times(1).Return(errors.New("unexpected error"))
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				return NewMockPublisher(ctrl)
			},
			event:    Proposal{ID: "id-1"},
			expected: errors.New("unexpected error"),
		},
		"allow do not send event after updating": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Proposal{ID: "id-1", Title: "name"}, nil)
				m.EXPECT().Update(gomock.Any()).Times(1).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(errors.New("unexpected error"))
				return m
			},
			event:    Proposal{ID: "id-1"},
			expected: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s, err := NewService(tc.dp(ctrl), tc.p(ctrl), NewMockEventRegistered(ctrl))
			require.Nil(t, err)

			err = s.HandleProposal(context.Background(), tc.event)
			if tc.expected == nil {
				require.Nil(t, err)
				return
			}

			require.ErrorContains(t, err, tc.expected.Error())
		})
	}
}
