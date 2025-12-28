// Package domain contains all logic pertaining to domains of the application
// a domain is broken into 4 major parts:
//
//  1. handler.go: This contains all handler function logic for declared route
//  2. route.go: This contains all specific domain routing declaration which will be later added to the main router see api.NewMux
//  3. model.go: This contains all logic pertaining to the model representing the domain
//     This includes the sql schema and the go struct definition for that domain
//  4. repo.go: This contains all specific database interactions such as CRUD operations
package domain
