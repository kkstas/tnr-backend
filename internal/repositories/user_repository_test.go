package repositories_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/kkstas/tnr-backend/internal/repositories"
	"github.com/kkstas/tnr-backend/internal/testutils"
)

func TestCreateAndFindAllUsers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db, cleanup := testutils.OpenTestDB(t)
	t.Cleanup(cleanup)
	userRepo := repositories.NewUserRepo(db)

	firstName := "John"
	lastName := "Doe"
	email := "john.doe@email.com"

	err := userRepo.CreateOne(ctx, firstName, lastName, email, "somepassword")
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}

	foundUsers, err := userRepo.FindAll(ctx)
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}

	testutils.AssertEqual(t, len(foundUsers), 1)
	testutils.AssertEqual(t, foundUsers[0].FirstName, firstName)
	testutils.AssertEqual(t, foundUsers[0].LastName, lastName)
	testutils.AssertEqual(t, foundUsers[0].Email, email)
	testutils.AssertValidDate(t, foundUsers[0].CreatedAt)
	if err := uuid.Validate(foundUsers[0].ID); err != nil {
		t.Errorf("expected id to be valid uuid, got error: %v", err)
	}
}

func TestFindOneUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db, cleanup := testutils.OpenTestDB(t)
	t.Cleanup(cleanup)
	userRepo := repositories.NewUserRepo(db)

	firstName := "John"
	lastName := "Doe"
	email := "john.doe@email.com"

	err := userRepo.CreateOne(ctx, firstName, lastName, email, "somepassword")
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}

	foundUsers, err := userRepo.FindAll(ctx)
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}

	foundUser, err := userRepo.FindOneByID(ctx, foundUsers[0].ID)
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}

	testutils.AssertEqual(t, foundUser.FirstName, firstName)
	testutils.AssertEqual(t, foundUser.LastName, lastName)
	testutils.AssertEqual(t, foundUser.Email, email)
	testutils.AssertValidDate(t, foundUser.CreatedAt)
	if err := uuid.Validate(foundUser.ID); err != nil {
		t.Errorf("expected id to be valid uuid, got error: %v", err)
	}
}
