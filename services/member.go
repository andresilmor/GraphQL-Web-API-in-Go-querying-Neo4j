package services

import (
	"CareXR_WebService/fixtures"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"fmt"
)

type Member = map[string]interface{}

type MemberService interface {
	MemberLogin(username string) (Member, error)
}

type neo4jMemberService struct {
	loader *fixtures.FixtureLoader
	driver neo4j.Driver
}

func NewMemberService(loader *fixtures.FixtureLoader, driver neo4j.Driver) MemberService {
	return &neo4jMemberService{loader: loader, driver: driver}
}

func (ms *neo4jMemberService) MemberLogin(username string) (_ Member, err error) {
	session := ms.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(fmt.Sprintf(`
		MATCH (c:Member {username: $username})
		CALL {
			OPTIONAL MATCH (i:Institution)<-[w:WORKS_IN]-(c) 
				WITH w{ .* , institution : PROPERTIES(i)} AS Institution
				WITH collect(Institution) AS Institutions
				RETURN Institutions
		}
		WITH Institutions, c AS Member
		RETURN {member: Member, institutions: Institutions} as member 
		`), map[string]interface{}{
			"username": username,
		})
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			return nil, err
		}

		member, _ := record.Get("member")
		return member.(map[string]interface{}), nil

		/*	FOR MULTIPLE RECORDS
			records, err := result.Collect()
			if err != nil {
				return nil, err
			}
			var results []map[string]interface{}
			for _, record := range records {
				movie, _ := record.Get("Member")

				results = append(results, movie.(map[string]interface{}))
				println(len(results))

			}
			return results, nil
		*/

	})

	session.Close()

	if results == nil {
		return nil, err
	}

	return results.(Member), err

}
