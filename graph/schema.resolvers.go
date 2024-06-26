package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.30

import (
	"CareXR_WebService/config"
	"CareXR_WebService/fixtures"
	"CareXR_WebService/graph/model"
	"CareXR_WebService/pkg/bcrypt"
	"CareXR_WebService/pkg/jwt"
	"CareXR_WebService/services/mongoDB"
	"CareXR_WebService/services/neo4j"
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

// MemberLogin is the resolver for the MemberLogin field.
func (r *queryResolver) MemberLogin(ctx context.Context, loginCredentials *model.LoginCredentials) (model.MemberLoginResponse, error) {
	service := neo4j.NewMemberService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		config.Neo4jDriver)

	fmt.Println(loginCredentials)
	//service.Get360imagesByInstitution("test")

	response, err := service.MemberLogin(loginCredentials.Email)

	if err != nil {
		return &model.Error{
			Message: "Invalid Credentials",
		}, nil

	}

	memberData := response["member"].(dbtype.Node).Props

	isValid := bcrypt.CheckPasswordHash(loginCredentials.Password, memberData["password"].(string))
	if !isValid {
		return &model.Error{
			Message: "Invalid Credentials",
		}, nil
	}

	memberOf := []*model.MemberOf{}
	memberOfData := response["institutions"].([]interface{})
	for _, element := range memberOfData {
		role := element.(map[string]interface{})["role"].(string)

		institutionData := element.(map[string]interface{})["institution"].(map[string]interface{})
		institutionUUID := institutionData["uuid"].(string)
		institutionLabel := institutionData["label"].(string)
		institutionName := institutionData["name"].(string)

		memberOf = append(memberOf, &model.MemberOf{&role, &model.Institution{&institutionUUID, &institutionLabel, &institutionName}})

	}

	uuid := memberData["uuid"].(string)
	name := memberData["name"].(string)
	username := memberData["username"].(string)

	tokenContent := map[string]any{
		"iss": "CareXR",
		"sub": &loginCredentials.Email,
		"aud": []string{"member"},
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}

	token, _ := jwt.GenerateToken(tokenContent)

	return &model.Member{
		UUID:     &uuid,
		Name:     &name,
		Username: &username,
		Email:    &loginCredentials.Email,
		Token:    &token,
		MemberOf: memberOf,
	}, nil
}

// MedicationToTake is the resolver for the MedicationToTake field.
func (r *queryResolver) MedicationToTake(ctx context.Context, isAvailable bool, pacientID *string, memberID *string, institutionID *string) ([]*model.MedicationToTake, error) {
	panic(fmt.Errorf("not implemented: MedicationToTake - MedicationToTake"))
}

// GetPanoramicSessions is the resolver for the GetPanoramicSessions field.
func (r *queryResolver) GetPanoramicSessions(ctx context.Context, institutionID *string, panoramicID *string, directedFor []*string, externalFormat *bool) ([]*model.PanoramicSession, error) {
	institution := institutionID

	ids := []string{*institution}
	list := mongoDB.GetPanoramicImages(ids, "", nil)

	hotspotsList := []*model.PanoramicSession{}

	service := neo4j.NewMemberService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		config.Neo4jDriver)

	for _, value := range list {
		var creator string = value.Meta.CreatedBy
		if *externalFormat {
			creator, _ = service.GetNameByUUID(creator)
		}

		mapping := []*model.HotspotPoint{}

		//fmt.Println(reflect.TypeOf(&value.Mapping))

		for _, mappingValue := range value.Mapping {

			mappingData := model.HotspotPointData{
				Alias: &mappingValue.Data.Alias,
			}

			mapping = append(mapping, &model.HotspotPoint{
				BoundingBox: &model.BoundingBox{
					X:      &mappingValue.BoundingBox.X,
					Y:      &mappingValue.BoundingBox.Y,
					Width:  &mappingValue.BoundingBox.Width,
					Height: &mappingValue.BoundingBox.Height,
				},
				Data: &mappingData,
			})

		}

		hotspotsList = append(hotspotsList, &model.PanoramicSession{
			UUID:  &value.UUID,
			Label: &value.Label,
			Meta: &model.PanoramicSessionMeta{
				CreatedAt: &value.Meta.CreatedAt,
				UpdatedAt: &value.Meta.UpdatedAt,
				CreatedBy: &creator,
			},
			ImageWidth: &value.ImageWidth,
			Mapping:    mapping,
		})

	}

	return hotspotsList, nil
}

func mapToPacient(m map[string]interface{}) (*model.Pacient, error) {
	var pacient model.Pacient

	// Use type assertions to retrieve values from the map
	if uuid, ok := m["uuid"].(string); ok {
		pacient.UUID = &uuid
	}

	if label, ok := m["label"].(string); ok {
		pacient.Label = &label
	}

	if name, ok := m["name"].(string); ok {
		pacient.Name = &name
	}

	// Repeat this process for all fields

	return &pacient, nil
}

// GetInstitutionPacients is the resolver for the GetInstitutionPacients field.
func (r *queryResolver) GetInstitutionPacients(ctx context.Context, institutionID *string) ([]*model.Pacient, error) {
	institution := *institutionID

	pacients := []*model.Pacient{}

	service := neo4j.NewInstitutionService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		config.Neo4jDriver)

	// Call GetInstitutionPacients to get []map[string]interface{} and convert it
	result, _ := service.GetInstitutionPacients(institution)
	for _, item := range result {
		var pacient model.Pacient
		uuid := item["uuid"].(string)
		pacient.UUID = &uuid
		label := item["label"].(string)
		pacient.Label = &label
		name := item["name"].(string)
		pacient.Name = &name

		pacients = append(pacients, &pacient)
	}

	return pacients, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
