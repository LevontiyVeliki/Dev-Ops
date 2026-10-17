package implementation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"

	"example.com/restaurant-service/implementation"
	"example.com/restaurant-service/repository"
	"example.com/restaurant-service/restaurant"
)

func newService() restaurant.Service {
	return implementation.NewService(repository.NewMemoryRepository(), log.NewNopLogger())
}

func TestCreateAndGet(t *testing.T) {
	ctx := context.Background()
	svc := newService()
	id, err := svc.Create(ctx, restaurant.Restaurant{Name: "  Шашлычная ", Address: "ул. Мира, 1", Rating: 4.5, Open: true})
	if err != nil || id == "" {
		t.Fatalf("Create: %q, %v", id, err)
	}
	got, err := svc.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	want := restaurant.Restaurant{ID: id, Name: "Шашлычная", Address: "ул. Мира, 1", Rating: 4.5, Open: true}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestCreateValidation(t *testing.T) {
	svc := newService()
	cases := []struct {
		r    restaurant.Restaurant
		want error
	}{
		{restaurant.Restaurant{Name: ""}, restaurant.ErrInvalidName},
		{restaurant.Restaurant{Name: "   "}, restaurant.ErrInvalidName},
		{restaurant.Restaurant{Name: "A", Rating: -0.1}, restaurant.ErrInvalidRating},
		{restaurant.Restaurant{Name: "A", Rating: 5.1}, restaurant.ErrInvalidRating},
	}
	for _, c := range cases {
		if _, err := svc.Create(context.Background(), c.r); !errors.Is(err, c.want) {
			t.Errorf("Create(%+v) = %v, want %v", c.r, err, c.want)
		}
	}
}

func TestChangeRating(t *testing.T) {
	ctx := context.Background()
	svc := newService()
	id, _ := svc.Create(ctx, restaurant.Restaurant{Name: "Пиццерия", Rating: 3})
	for _, r := range []float32{0, 5, 4.2} {
		if err := svc.ChangeRating(ctx, id, r); err != nil {
			t.Fatalf("ChangeRating(%v): %v", r, err)
		}
	}
	if got, _ := svc.GetByID(ctx, id); got.Rating != 4.2 {
		t.Errorf("rating = %v, want 4.2", got.Rating)
	}
	if err := svc.ChangeRating(ctx, id, 6); !errors.Is(err, restaurant.ErrInvalidRating) {
		t.Errorf("rating 6: %v", err)
	}
	if err := svc.ChangeRating(ctx, "nope", 4); !errors.Is(err, restaurant.ErrRestaurantNotFound) {
		t.Errorf("unknown id: %v", err)
	}
	if _, err := svc.GetByID(ctx, "nope"); !errors.Is(err, restaurant.ErrRestaurantNotFound) {
		t.Errorf("GetByID unknown: %v", err)
	}
}
