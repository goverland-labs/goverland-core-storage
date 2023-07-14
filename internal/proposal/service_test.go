package proposal

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	coreevents "github.com/goverland-labs/platform-events/events/core"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// todo: add units for converting from and to internal models

var defaultPublisher = func(ctrl *gomock.Controller) Publisher {
	m := NewMockPublisher(ctrl)
	m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	return m
}

var defaultDaoProvider = func(ctrl *gomock.Controller) DaoProvider {
	m := NewMockDaoProvider(ctrl)
	m.EXPECT().GetIDByOriginalID(gomock.Any()).AnyTimes().Return(uuid.New(), nil)
	return m
}

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
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(nil)
				return m
			},
			event:    Proposal{ID: "id-1"},
			expected: nil,
		},
		"correct updating": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Proposal{ID: "id-1", Title: "updated", Quorum: 50}, nil)
				m.EXPECT().Update(gomock.Any()).Times(1).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(nil)
				return m
			},
			event:    Proposal{ID: "id-1", Title: "name", Quorum: 50},
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
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(errors.New("unexpected error"))
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
				m.EXPECT().GetByID(gomock.Any()).Times(1).Return(&Proposal{ID: "id-1", Title: "name", Quorum: 50}, nil)
				m.EXPECT().Update(gomock.Any()).Times(1).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(errors.New("unexpected error"))
				return m
			},
			event:    Proposal{ID: "id-1", Quorum: 50},
			expected: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer func() {
				<-time.After(10 * time.Millisecond)
				ctrl.Finish()
			}()

			s, err := NewService(tc.dp(ctrl), tc.p(ctrl), NewMockEventRegistered(ctrl), defaultDaoProvider(ctrl))
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

func TestUnitProcessAvailableForVoting(t *testing.T) {
	for name, tc := range map[string]struct {
		dp       func(ctrl *gomock.Controller) DataProvider
		er       func(ctrl *gomock.Controller) EventRegistered
		expected error
	}{
		"error on getting data": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetAvailableForVoting(gomock.Any()).MaxTimes(1).Return(nil, gorm.ErrRecordNotFound)
				return m
			},
			er: func(ctrl *gomock.Controller) EventRegistered {
				return NewMockEventRegistered(ctrl)
			},
			expected: gorm.ErrRecordNotFound,
		},
		"voting has started": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetAvailableForVoting(gomock.Any()).MaxTimes(1).Return([]*Proposal{
					{
						CreatedAt: time.Now().Add(-time.Hour * 24),
						Start:     int(time.Now().Add(-time.Hour * 2).Unix()),
						End:       int(time.Now().Add(time.Hour * 24).Unix()),
					},
				}, nil)
				return m
			},
			er: func(ctrl *gomock.Controller) EventRegistered {
				m := NewMockEventRegistered(ctrl)
				m.EXPECT().EventExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				m.EXPECT().RegisterEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(_ context.Context, _, _, event string) error {
						if coreevents.SubjectProposalVotingStarted != event {
							ctrl.T.Errorf("wrong subject event: %s instead of %s", event, coreevents.SubjectProposalVotingStarted)
						}

						return nil
					})
				return m
			},
			expected: nil,
		},
		"voting has ended": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetAvailableForVoting(gomock.Any()).MaxTimes(1).Return([]*Proposal{
					{
						CreatedAt: time.Now().Add(-time.Hour * 24),
						Start:     int(time.Now().Add(-time.Hour * 25).Unix()),
						End:       int(time.Now().Add(-time.Hour * 1).Unix()),
					},
				}, nil)
				return m
			},
			er: func(ctrl *gomock.Controller) EventRegistered {
				m := NewMockEventRegistered(ctrl)
				m.EXPECT().EventExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				m.EXPECT().RegisterEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(_ context.Context, _, _, event string) error {
						if coreevents.SubjectProposalVotingEnded != event {
							ctrl.T.Errorf("wrong subject event: %s instead of %s", event, coreevents.SubjectProposalVotingEnded)
						}

						return nil
					})
				return m
			},
			expected: nil,
		},
		"voting is coming": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetAvailableForVoting(gomock.Any()).MaxTimes(1).Return([]*Proposal{
					{
						CreatedAt: time.Now().Add(-time.Hour * 24),
						Start:     int(time.Now().Add(time.Minute * 25).Unix()),
						End:       int(time.Now().Add(time.Hour * 24 * 7).Unix()),
					},
				}, nil)
				return m
			},
			er: func(ctrl *gomock.Controller) EventRegistered {
				m := NewMockEventRegistered(ctrl)
				m.EXPECT().EventExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				m.EXPECT().RegisterEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(_ context.Context, _, _, event string) error {
						if coreevents.SubjectProposalVotingStartsSoon != event {
							ctrl.T.Errorf("wrong subject event: %s instead of %s", event, coreevents.SubjectProposalVotingStartsSoon)
						}

						return nil
					})
				return m
			},
			expected: nil,
		},
		"send single event": {
			dp: func(ctrl *gomock.Controller) DataProvider {
				m := NewMockDataProvider(ctrl)
				m.EXPECT().GetAvailableForVoting(gomock.Any()).MaxTimes(1).Return([]*Proposal{
					{
						CreatedAt: time.Now().Add(-time.Hour * 24),
						Start:     int(time.Now().Add(-time.Minute * 30).Unix()),
						End:       int(time.Now().Add(time.Hour * 24).Unix()),
					},
				}, nil)
				return m
			},
			er: func(ctrl *gomock.Controller) EventRegistered {
				m := NewMockEventRegistered(ctrl)
				m.EXPECT().EventExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				m.EXPECT().RegisterEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0).Return(nil)
				return m
			},
			expected: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer func() {
				<-time.After(10 * time.Millisecond)
				ctrl.Finish()
			}()
			s, err := NewService(tc.dp(ctrl), defaultPublisher(ctrl), tc.er(ctrl), nil)
			require.Nil(t, err)

			err = s.processAvailableForVoting(context.TODO())
			if tc.expected == nil {
				require.Nil(t, err)
				return
			}

			require.ErrorContains(t, err, tc.expected.Error())
		})
	}
}

func TestUnitCheckSpecificUpdate(t *testing.T) {
	for name, tc := range map[string]struct {
		existed Proposal
		new     Proposal
		er      func(ctrl *gomock.Controller) EventRegistered
		p       func(ctrl *gomock.Controller) Publisher
	}{
		"on voting reached": {
			existed: Proposal{ScoresTotal: 20, Quorum: 50, State: "state-1"},
			new:     Proposal{ScoresTotal: 50, Quorum: 50, State: "state-1"},
			er: func(ctrl *gomock.Controller) EventRegistered {
				m := NewMockEventRegistered(ctrl)
				m.EXPECT().EventExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				m.EXPECT().RegisterEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(_ context.Context, _, _, event string) error {
						if coreevents.SubjectProposalVotingQuorumReached != event {
							ctrl.T.Errorf("wrong subject event: %s instead of %s", event, coreevents.SubjectProposalVotingQuorumReached)
						}

						return nil
					})
				return m
			},
			p: defaultPublisher,
		},
		"on state update": {
			existed: Proposal{ScoresTotal: 0, Quorum: 50, State: "state-1"},
			new:     Proposal{ScoresTotal: 0, Quorum: 50, State: "state-2"},
			er: func(ctrl *gomock.Controller) EventRegistered {
				return NewMockEventRegistered(ctrl)
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(_ context.Context, event string, _ any) error {
						if coreevents.SubjectProposalUpdatedState != event {
							ctrl.T.Errorf("wrong subject event: %s instead of %s", event, coreevents.SubjectProposalUpdatedState)
						}

						return nil
					})
				return m
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer func() {
				<-time.After(10 * time.Millisecond)
				ctrl.Finish()
			}()
			s, err := NewService(nil, tc.p(ctrl), tc.er(ctrl), nil)
			require.Nil(t, err)

			s.checkSpecificUpdate(context.TODO(), tc.new, tc.existed)
		})
	}
}

func TestUnitRegisterEventOnce(t *testing.T) {
	for name, tc := range map[string]struct {
		er func(ctrl *gomock.Controller) EventRegistered
		p  func(ctrl *gomock.Controller) Publisher
	}{
		"send event if not registered before": {
			er: func(ctrl *gomock.Controller) EventRegistered {
				m := NewMockEventRegistered(ctrl)
				m.EXPECT().EventExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(false, nil)
				m.EXPECT().RegisterEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				return m
			},
		},
		"do not send event twice": {
			er: func(ctrl *gomock.Controller) EventRegistered {
				m := NewMockEventRegistered(ctrl)
				m.EXPECT().EventExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				m.EXPECT().RegisterEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(0).Return(nil)
				return m
			},
		},
		"do not send event on getting error": {
			er: func(ctrl *gomock.Controller) EventRegistered {
				m := NewMockEventRegistered(ctrl)
				m.EXPECT().EventExist(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(false, errors.New("unspecified error"))
				m.EXPECT().RegisterEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0).Return(nil)
				return m
			},
			p: func(ctrl *gomock.Controller) Publisher {
				m := NewMockPublisher(ctrl)
				m.EXPECT().PublishJSON(gomock.Any(), gomock.Any(), gomock.Any()).Times(0).Return(nil)
				return m
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer func() {
				<-time.After(10 * time.Millisecond)
				ctrl.Finish()
			}()
			s, err := NewService(nil, tc.p(ctrl), tc.er(ctrl), nil)
			require.Nil(t, err)

			s.registerEventOnce(context.TODO(), Proposal{}, "group", "subject")
		})
	}
}
