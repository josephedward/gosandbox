package test

import (
	"database/sql"
	"fmt"
	"gosandbox/acloud"
	"os"
	"testing"
)

func TestDb(t *testing.T) {
	const fileName = "sqlite.db"

	os.Remove(fileName)

	db, err := sql.Open("sqlite3", fileName)
	fmt.Println("db : ", db)

	sandboxRepository := acloud.NewSQLiteRepository(db)

	sandboxRepository.Migrate()

	fmt.Println("sandboxRepository : ", sandboxRepository)

	example := acloud.SandboxCredential{
		User:      "testUser",
		Password:  "testPassword",
		URL:       "testURL",
		KeyID:     "testKeyID",
		AccessKey: "testAccessKey",
	}

	fmt.Println("example : ", example)

	createdExample, err := sandboxRepository.Create(example)
	fmt.Println("createdExample : ", createdExample)

	gotExample, err := sandboxRepository.GetByName("testUser")
	fmt.Printf("get by User: %+v\n", gotExample)

	createdExample.URL = "newURL"
	sandboxRepository.Update(createdExample.ID, *createdExample)

	all, err := sandboxRepository.All()
	fmt.Printf("\nAll sandboxes:\n")
	for _, sandbox := range all {
		fmt.Printf("sandbox: %+v\n", sandbox)
	}

	sandboxRepository.Delete(createdExample.ID)

	all, err = sandboxRepository.All()
	fmt.Printf("\nAll sandboxes:\n")
	for _, sandbox := range all {
		fmt.Printf("sandbox: %+v\n", sandbox)
	}

}
