package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goops-top/ent-pro/ent"
	"github.com/goops-top/ent-pro/ent/car"
	"github.com/goops-top/ent-pro/ent/group"
	"github.com/goops-top/ent-pro/ent/user"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// create the user
	u, uErr := CreateUser(context.Background(), client)
	if uErr != nil {
		log.Fatalf("failed creating the user:%v\n", uErr)
	}

	fmt.Println(u)

	// query the user
	qu, quErr := QueryUser(context.Background(), client)
	if quErr != nil {
		log.Fatalf("failed quering the user:%v\n", quErr)
	}

	fmt.Println(qu)

	fmt.Println("create the cars and add the cars to a user")
	carsUser, carErr := CreateCars(context.Background(), client)
	if carErr != nil {
		log.Fatalf("failed add the cars:%v\n", carErr)
	}

	fmt.Println("query the cars")
	queryErr := QueryCars(context.Background(), carsUser)
	fmt.Println(queryErr)

	fmt.Println("query the user of the car.")
	queryCarErr := QueryCarUsers(context.Background(), carsUser)
	fmt.Println(queryCarErr)

	fmt.Println("create a data graph.")
	CreateGraph(context.Background(), client)

	fmt.Println("query group with `family` within the data graph")
	QueryFamily(context.Background(), client)

	fmt.Println("query user with `bgbiao` within the data graph")
	QueryBGBiaoCars(context.Background(), client)

	fmt.Println("query group-users within the data graph")
	QueryGroupWithUsers(context.Background(), client)

}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %v", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.NameEQ("a8m")).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %v", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

// create 2 cars and adding them to a user
func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	// create a new car with model "Tesla."
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating car:%v", err)
	}

	// create a new car with model "Ford".
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating car:%v", err)
	}
	log.Println("car was created: ", ford)

	// create a new user,and add it the 2cars.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating user: %v", err)
	}

	log.Println("user was created: ", a8m)

	return a8m, nil

}

// query the cars edge.
func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %v", err)
	}

	log.Println("returned cars:", cars)

	// filter the specific cars
	ford, err := a8m.QueryCars().
		Where(car.ModelEQ("Ford")).
		Only(ctx)

	if err != nil {
		return fmt.Errorf("failed quering user cars: %v", err)
	}

	log.Println(ford)
	return nil

}

// quering the inverse edge.
func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed quering user cars: %v", err)
	}

	// Query the inverse edge.
	for _, ca := range cars {
		owner, err := ca.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %v", ca.Model, err)
		}
		log.Printf("car %q owner: %q\n", ca.Model, owner.Name)
	}

	return nil

}

// create the graph and data with the node and edge.
func CreateGraph(ctx context.Context, client *ent.Client) error {
	// first , create the users.
	bgbiao, err := client.User.
		Create().
		SetAge(28).
		SetName("BGBiao").
		Save(ctx)

	if err != nil {
		return nil
	}

	xxbandy, err := client.User.
		Create().
		SetAge(28).
		SetName("xxbandy").
		Save(ctx)

	if err != nil {
		return err
	}

	// then, create the cars, and attach them to the users in the creation.
	_, err = client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		SetOwner(bgbiao).
		Save(ctx)

	if err != nil {
		return err
	}

	_, err = client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()).
		SetOwner(xxbandy).
		Save(ctx)

	if err != nil {
		return err
	}

	_, err = client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		SetOwner(bgbiao).
		Save(ctx)
	if err != nil {
		return err
	}

	// create the groups, and add their users in the creation.
	_, err = client.Group.
		Create().
		SetName("family").
		AddUsers(bgbiao, xxbandy).
		Save(ctx)

	if err != nil {
		return err
	}

	_, err = client.Group.
		Create().
		SetName("myself").
		AddUsers(bgbiao).
		Save(ctx)

	if err != nil {
		return err
	}

	log.Println("the graph was created successfully.")

	return nil

}

// run a few queries on the graph.

// get all user's cars within the group named "family".

func QueryFamily(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.
		Query().
		Where(group.Name("family")).
		QueryUsers().
		QueryCars().
		All(ctx)

	if err != nil {
		return fmt.Errorf("failed getting cars: %v", err)
	}

	log.Println("cars returned:", cars)
	return nil
}

// get a specifi user and some source.
func QueryBGBiaoCars(ctx context.Context, client *ent.Client) error {
	// get "bgbiao" from previous steps.
	bgbiao := client.User.
		Query().
		Where(
			user.HasCars(),
			user.Name("BGBiao"),
		).
		OnlyX(ctx)

	cars, err := bgbiao.
		QueryGroups().
		QueryUsers().
		QueryCars().
		Where(
			car.Not(
				car.ModelEQ("Mazda"),
			),
		).All(ctx)

	if err != nil {
		return err
	}

	log.Println("cars returned:", cars)

	return nil
}

// get all groups that have users(query with a look-aside predicate)
func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
	groups, err := client.Group.
		Query().
		Where(group.HasUsers()).
		All(ctx)

	if err != nil {
		return fmt.Errorf("failed getting groups: %v", err)
	}

	log.Println("groups returned:", groups)

	return nil
}
