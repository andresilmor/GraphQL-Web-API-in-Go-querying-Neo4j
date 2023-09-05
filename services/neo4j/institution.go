package neo4j

import (
	"CareXR_WebService/fixtures"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Institution = map[string]interface{}

type InstitutionService interface {
	GetInstitutionPacients(institutionID string) ([]Patient, error)
}

type neo4jInstitutionService struct {
	loader *fixtures.FixtureLoader
	driver neo4j.Driver
}

func NewInstitutionService(loader *fixtures.FixtureLoader, driver neo4j.Driver) InstitutionService {
	return &neo4jInstitutionService{loader: loader, driver: driver}
}

type Patient = map[string]interface{}

func (ms *neo4jInstitutionService) GetInstitutionPacients(institutionID string) ([]Patient, error) {
	session := ms.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(fmt.Sprintf(`
		MATCH (n:Pacient)-[:UNDER_CARE_OF]->(Institution{uuid: $uuid})
		RETURN PROPERTIES(n) AS institution_patients
		`), map[string]interface{}{
			"uuid": institutionID,
		})
		if err != nil {
			return nil, err
		}

		records, err := result.Collect()

		fmt.Println(institutionID)
		if err != nil {
			return nil, err
		}
		var results []map[string]interface{}
		for _, record := range records {
			patient, _ := record.Get("institution_patients")

			results = append(results, patient.(map[string]interface{}))
		}
		return results, nil

	})

	session.Close()

	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, err
	}

	return results.([]Patient), nil
}
