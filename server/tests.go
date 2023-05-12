package server

import (
	"context"
	"io"
	"krlosaren/go/grpc/models"
	"krlosaren/go/grpc/repository"
	"krlosaren/go/grpc/studentpbf"
	"krlosaren/go/grpc/testpbf"
	"log"
	"time"
)

type TestServer struct {
	repo repository.Repository
	testpbf.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{
		repo: repo,
	}
}

func (s *TestServer) GetTest(ctx context.Context, r *testpbf.GetTestRequest) (*testpbf.Test, error) {
	test, err := s.repo.GetTest(ctx, r.GetId())

	if err != nil {
		return nil, err
	}

	return &testpbf.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpbf.Test) (*testpbf.SetTestResponse, error) {
	test := &models.Test{
		Id:   req.GetId(),
		Name: req.GetName(),
	}

	err := s.repo.SetTest(ctx, test)

	if err != nil {
		return nil, err
	}

	return &testpbf.SetTestResponse{
		Id:   test.Id,
		Name: test.Name,
	}, nil

}

func (s *TestServer) SetQuestion(stream testpbf.TestService_SetQuestionServer) error {

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&testpbf.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
			return err
		}
		question := &models.Question{
			Id:       msg.GetId(),
			Answer:   msg.GetAnswer(),
			Question: msg.GetQuestion(),
			TestId:   msg.GetTestId(),
		}
		err = s.repo.SetQuestion(context.Background(), question)

		if err != nil {
			return stream.SendAndClose(&testpbf.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) SetEnrollStudents(stream testpbf.TestService_SetEnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&testpbf.SetQuestionResponse{
				Ok: true,
			})
		}

		if err != nil {
			return err
		}

		enrollment := &models.Enrollment{
			StudentId: msg.GetStudentId(),
			TestId:    msg.GetTestId(),
		}

		err = s.repo.SetEnrollment(context.Background(), enrollment)

		if err != nil {
			return stream.SendAndClose(&testpbf.SetQuestionResponse{
				Ok: false,
			})
		}

	}
}

func (s *TestServer) GetStudentsPerTest(req *testpbf.GetStudentsPerTestRequest, stream testpbf.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(context.Background(), req.GetTestId())

	if err != nil {
		return err
	}

	for _, student := range students {
		student := &studentpbf.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		}
		err = stream.Send(student)
		time.Sleep(2 * time.Second)

		if err != nil {
			return err
		}

	}
	return nil
}

func (s *TestServer) TakeTest(stream testpbf.TestService_TakeTestServer) error {
	questions, err := s.repo.GetQuestionsPerTest(context.Background(), "t1")

	if err != nil {
		return err
	}

	i := 0

	var currentQuestion = &models.Question{}

	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}

		if i <= len(questions) {
			questionToSend := &testpbf.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}

			err := stream.Send(questionToSend)

			if err != nil {
				return err
			}

			i++
		}

		answer, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		log.Println("Answer received from server: ", answer.GetAnswer())
	}

}
