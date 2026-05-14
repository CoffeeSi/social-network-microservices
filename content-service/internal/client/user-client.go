package client

import "context"

type MockUserClient struct{}

func NewMockUserClient() *MockUserClient {
	return &MockUserClient{}
}

func (c *MockUserClient) UserExists(ctx context.Context, userId string) (bool, error) {
	if userId == "" {
		return false, nil
	}
	return true, nil
}
