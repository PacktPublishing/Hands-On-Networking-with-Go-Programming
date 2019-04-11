package schema

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/graphql/gqlid"
	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/grpc/interface/order"
	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/client"
	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/resources/company"
)

// These constants define the ID types for the remote services.
// The GraphQL endpoint uses two remote services, one for "Orders" and one for "Companies".
const companyRestService = "companyrest"
const companyRestResource = "company"

// NewResolver creates the new GraphQL resolver, with the relevant REST and gRPC clients which provide access
// to data. These could, in fact, access data directly.
func NewResolver(oc order.OrdersClient, cc *client.CompanyClient) *Resolver {
	return &Resolver{
		OrderClient:   oc,
		CompanyClient: cc,
	}
}

// A Resolver resolves GraphQL requests: both Mutations and Queries.
type Resolver struct {
	//TODO: The Resolver here accepts real clients, but this makes it hard to test.
	//TODO: It would be better to replace them with interfaces so that they can be mocked out.
	OrderClient   order.OrdersClient
	CompanyClient *client.CompanyClient
}

// Mutation returns the MutationResolver which handles any mutation requests.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{
		resolver:      r,
		companyClient: r.CompanyClient,
		orderClient:   r.OrderClient,
	}
}

// Query returns resolvers for the GraphQL queries.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{
		resolver:      r,
		companyClient: r.CompanyClient,
		orderClient:   r.OrderClient,
	}
}

// Order returns the Order resolver. This isn't a query that clients can use, but is actually
// a way for any Order types retned by a query to lookup additional data, potentially from other
// services. The orderResolver here uses DataLoader middleware to batch up individual requests
// together into less backend operations, reducing latency.
func (r *Resolver) Order() OrderResolver {
	return &orderResolver{r}
}

// The mutationResolver handles mutation requests.
type mutationResolver struct {
	resolver      *Resolver
	companyClient *client.CompanyClient
	orderClient   order.OrdersClient
}

// CreateCompany is the createCompany mutation on the GraphQL schema. It uses the company client to
// create a new company.
func (r *mutationResolver) CreateCompany(ctx context.Context, new NewCompany) (created *Company, err error) {
	fmt.Println("mutationResolver: CreateCompany")
	id, err := r.companyClient.Post(company.Company{
		Name:      new.Name,
		VATNumber: new.VatNumber,
		Address: company.Address{
			Address1: new.Address.Address1,
			Address2: new.Address.Address2,
			Address3: new.Address.Address3,
			Address4: new.Address.Address4,
			Country:  new.Address.Country,
			Postcode: new.Address.Postcode,
		},
	})
	if err != nil {
		err = fmt.Errorf("failed to post company: %v", err)
		return
	}
	created = &Company{
		ID:        gqlid.Int64(companyRestService, companyRestResource, int64(id.ID)).Encoded(),
		Name:      new.Name,
		VatNumber: new.VatNumber,
		Address: Address{
			Address1: new.Address.Address1,
			Address2: new.Address.Address2,
			Address3: new.Address.Address3,
			Address4: new.Address.Address4,
			Country:  new.Address.Country,
			Postcode: new.Address.Postcode,
		},
	}
	return
}

// CreateOrder is the createOrder mutation on the GraphQL schema.
func (r *mutationResolver) CreateOrder(ctx context.Context, companyID string, new NewOrder) (o *Order, err error) {
	fmt.Println("mutationResolver: CreateOrder")
	id, err := gqlid.ParseInt64For(companyRestService, companyRestResource, companyID)
	if err != nil {
		err = fmt.Errorf("company ID was not in expected format: %v", err)
		return
	}
	fmt.Println("parsed id")
	items := make([]*order.Item, len(new.Items))
	for i, itm := range new.Items {
		items[i] = &order.Item{
			Id: itm,
		}
		//TODO: Look up the price and description of the items.
	}
	fmt.Println("adding order")
	ar, err := r.orderClient.Add(ctx, &order.AddRequest{
		Fields: &order.Fields{
			CompanyId: id,
			Items:     items,
			UserId:    "abc",
		},
	})
	if err != nil {
		return
	}
	fmt.Println("sorting out")
	o.CompanyID = id
	o.ID = ar.Id
	o.Items = make([]Item, len(new.Items))
	for i, itm := range new.Items {
		o.Items[i].ID = itm
		//TODO: Fill out the price and description of the items.
	}
	return
}

// orderResolver resolves any requirements that the order has when being loaded. In this case, the *Order type
// doesn't have a Company field on it, but within the GraphQL schema, the Order type has a Company child field.
// This resolver is responsible for loading the Company records and populating those child fields.
type orderResolver struct{ *Resolver }

// Company populates the Company field for the Order GraphQL type. In this case, it uses the GraphQL dataloader
// to group together multiple requests.
func (r *orderResolver) Company(ctx context.Context, obj *Order) (*Company, error) {
	fmt.Println("Order:Company")
	// This version doesn't use the data loader, but the CompanyID on the returned value has been converted to a string
	// as is usual for GraphQL, so it won't work any more.
	// return r.Query().Company(ctx, obj.CompanyID)
	// Use the data loader instead.
	return ctx.Value(companyLoaderKey).(*CompanyLoader).Load(obj.CompanyID)
}

// queryResolver implements the queries that can be made against the GraphQL endpoint.
type queryResolver struct {
	resolver      *Resolver
	companyClient *client.CompanyClient
	orderClient   order.OrdersClient
}

// Orders implements the Orders GraphQL query.
func (r *queryResolver) Orders(ctx context.Context, companyID string, ids []string) (orders []Order, err error) {
	fmt.Println("Query:Orders")
	//TODO: Implement a multi-get at the backend, or iterate through the ids.
	gr, err := r.orderClient.Get(ctx, &order.GetRequest{
		Id: companyID,
	})
	if err != nil {
		return
	}
	orders = append(orders, Order{
		ID:        base64.StdEncoding.EncodeToString([]byte(gr.Order.Id)),
		CompanyID: gr.Order.Fields.CompanyId,
		Items:     mapItemsToGraphQL(gr.Order.Fields.Items),
	})
	return
}

func mapItemsToGraphQL(input []*order.Item) (output []Item) {
	output = make([]Item, len(input))
	for i, in := range input {
		output[i].Desc = in.Desc
		output[i].ID = base64.StdEncoding.EncodeToString([]byte(in.Id))
		output[i].Price = int(in.Price)
	}
	return
}

// Company implements the Company GraphQL query.
func (r *queryResolver) Company(ctx context.Context, id string) (result *Company, err error) {
	fmt.Println("queryResolver: Company", id)
	int64id, err := gqlid.ParseInt64For(companyRestService, companyRestResource, id)
	if err != nil {
		return
	}
	c, err := r.companyClient.Get(int64id)
	if err != nil {
		return
	}
	result = &Company{
		ID:        id,
		Name:      c.Name,
		VatNumber: c.VATNumber,
		Address: Address{
			Address1: c.Address.Address1,
			Address2: c.Address.Address2,
			Address3: c.Address.Address3,
			Address4: c.Address.Address4,
			Country:  c.Address.Country,
			Postcode: c.Address.Postcode,
		},
	}
	return
}

// Companies implments the Companies GraphQL query.
func (r *queryResolver) Companies(ctx context.Context, ids []string) (op []Company, err error) {
	fmt.Println("queryResolver: Company", ids)
	// This is inefficient, because it simply calls the other resolver multiple times in a loop.
	// It would be more realistic to use an endpoint which supports retrieving multiple items in a
	// single operation.
	op = make([]Company, len(ids))
	for i, id := range ids {
		c, cerr := r.Company(ctx, id)
		if cerr != nil {
			err = cerr
			return
		}
		op[i] = *c
	}
	return
}
