package schema

import (
	"context"
	"net/http"
	"time"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/graphql/gqlid"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/restexample/client"
)

type dataLoaderMiddlewareKey string

const companyLoaderKey = dataLoaderMiddlewareKey("dataloaderCompany")

// WithCompanyDataloaderMiddleware populates the Data Loader middleware for loading company data.
func WithCompanyDataloaderMiddleware(companyClient *client.CompanyClient, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := NewCompanyLoader(CompanyLoaderConfig{
			Fetch: func(ids []int64) (companies []*Company, errs []error) {
				// Load the data from the client.
				cc, err := companyClient.GetMany(ids)
				if err != nil {
					errs = append(errs, err)
					return
				}
				// Map from the client returned by the REST client to the client GraphQL type.
				companies = make([]*Company, len(cc))
				for i, c := range cc {
					companies[i] = &Company{
						ID:        gqlid.Int64(companyRestService, companyRestService, c.ID).Encoded(),
						Name:      c.Name,
						VatNumber: "Unknown",
					}
				}
				return
			},
			MaxBatch: 100,
			Wait:     10 * time.Millisecond,
		})
		ctx := context.WithValue(r.Context(), companyLoaderKey, l)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
