package repository

import (
	"context"
	"fmt"
	"testing"

	learn_go_database "github.com/izzanzahrial/learn-go-database"
	"github.com/izzanzahrial/learn-go-database/entity"
)

func TestCommentInsert(t *testing.T) {
	commentRepository := NewCommentRepository(learn_go_database.GetConnection())

	ctx := context.Background()
	comment := entity.Comment{
		Email:   "izzanzahrial@test.com",
		Comment: "Test Repository",
	}
	result, err := commentRepository.Insert(ctx, comment)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestFindById(t *testing.T) {
	commentRepository := NewCommentRepository(learn_go_database.GetConnection())
	ctx := context.Background()
	comment, err := commentRepository.FindById(ctx, 37)
	if err != nil {
		panic(err)
	}
	fmt.Println(comment)
}

func TestFindAll(t *testing.T) {
	commentRepository := NewCommentRepository(learn_go_database.GetConnection())
	ctx := context.Background()
	comments, err := commentRepository.FindAll(ctx)
	if err != nil {
		panic(err)
	}
	for _, comment := range comments {
		fmt.Println(comment)
	}
}
