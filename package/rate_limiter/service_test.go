package rate_limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RateLimiterServiceImplTestSuite struct {
	suite.Suite

	sut                           *RateLimiterServiceImpl
	ipConfig                      RateLimiterConfig
	rateLimiterItemRepositoryStub RateLimiterItemRepository
}

func (s *RateLimiterServiceImplTestSuite) SetupTest() {
	s.rateLimiterItemRepositoryStub = NewRateLimiterItemRepositoryMemory()
	s.ipConfig = RateLimiterConfig{
		LimitPerInterval:  3,
		IntervalInSeconds: 2,
		BanTimeInSeconds:  2,
	}
	s.sut = NewRateLimiterServiceImpl(s.rateLimiterItemRepositoryStub)
}

func (s *RateLimiterServiceImplTestSuite) TestConfig() {
	s.sut.AddConfig("ip", s.ipConfig)

	config := s.sut.GetConfig("ip")

	s.Equal(s.ipConfig, config)
}

func (s *RateLimiterServiceImplTestSuite) TestExecute() {
	s.sut.AddConfig("ip", s.ipConfig)
	key := "any-key"
	configType := "ip"

	for i := 0; i < s.ipConfig.LimitPerInterval; i++ {
		s.True(s.sut.Execute(key, configType))
	}

	s.False(s.sut.Execute(key, configType))
	s.False(s.sut.Execute(key, configType))

	time.Sleep(time.Duration(s.ipConfig.BanTimeInSeconds) * time.Second)

	s.True(s.sut.Execute(key, configType))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterServiceImplTestSuite))
}
