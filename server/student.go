package server

import (
	"context"
	"krlosaren/go/grpc/models"
	"krlosaren/go/grpc/repository"
	"krlosaren/go/grpc/studentpbf"
)

type Server struct {
	repo repository.Repository
	studentpbf.UnimplementedStudentServiceServer
}

func NewStudentServer(repo repository.Repository) *Server {
	return &Server{
		repo: repo,
	}
}

func (s *Server) GetStudent(ctx context.Context, req *studentpbf.GetStudentRequest) (*studentpbf.Student, error) {
	student, err := s.repo.GetStudent(ctx, req.GetId())

	if err != nil {
		return nil, err
	}

	return &studentpbf.Student{
		Id:   student.Id,
		Name: student.Name,
		Age:  student.Age,
	}, nil
}

func (s *Server) SetStudent(ctx context.Context, req *studentpbf.Student) (*studentpbf.SetStudentResponse, error) {
	student := &models.Student{
		Id:   req.GetId(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}

	err := s.repo.SetStudent(ctx, student)
	if err != nil {
		return nil, err
	}

	return &studentpbf.SetStudentResponse{
		Id: student.Id,
	}, nil
}
