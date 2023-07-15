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

// Get360Hotspot is the resolver for the Get360Hotspot field.
func (r *queryResolver) Get360Hotspot(ctx context.Context, institutionID *string, hotspotID *string, directedFor []*string, externalFormat *bool) ([]*model.Hotspot, error) {
	list := mongoDB.Get360images(*institutionID, *hotspotID, directedFor)

	hotspotsList := []*model.Hotspot{}

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
				Transform: &model.HotspotPointTransform{
					Position: &model.Position{
						X: &mappingValue.Transform.Position.X,
						Y: &mappingValue.Transform.Position.Y,
						Z: &mappingValue.Transform.Position.Z,
					},
					Scale: &model.Scale{
						Width:  &mappingValue.Transform.Scale.Width,
						Height: &mappingValue.Transform.Scale.Height,
					},
				},
				Data: &mappingData,
			})

		}

		hotspotsList = append(hotspotsList, &model.Hotspot{
			UUID:  &value.UUID,
			Label: &value.Label,
			Meta: &model.HotspotMeta{
				CreatedAt: &value.Meta.CreatedAt,
				UpdatedAt: &value.Meta.UpdatedAt,
				CreatedBy: &creator,
			},
			Mapping: mapping,
		})

	}

	return hotspotsList, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
