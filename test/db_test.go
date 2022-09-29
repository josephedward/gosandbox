package test

import (
	"gosandbox/acloud"
	// "gosandbox/cli"
	// "gosandbox/aws"
	"fmt"
	"testing"
	"os"
	// "errors"
	// "log"
	"database/sql"
)

func TestDb(t *testing.T) {

	const fileName = "sqlite.db"

	os.Remove(fileName)

	db, err := sql.Open("sqlite3", fileName)
	fmt.Println("db : ", db)
	// if err != nil {
		
	// }

	sandboxRepository := acloud.NewSQLiteRepository(db)
	

	// if err := 
	sandboxRepository.Migrate();
	//  err != nil {
	// 	cli.PrintIfErr(err)
	// }

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
	// if err != nil {
	// 	
	// // }
	fmt.Println("createdExample : ", createdExample)

	gotExample, err := sandboxRepository.GetByName("testUser")
	// if err != nil {
	// 	
	// }

	fmt.Printf("get by User: %+v\n", gotExample)

	createdExample.URL = "newURL"
	if _, err := sandboxRepository.Update(createdExample.ID, *createdExample); err != nil {
		
	}

	all, err := sandboxRepository.All()
	if err != nil {
		
	}

	fmt.Printf("\nAll sandboxes:\n")
	for _, sandbox := range all {
		fmt.Printf("sandbox: %+v\n", sandbox)
	}

	if err := sandboxRepository.Delete(createdExample.ID); err != nil {
		
	}

	all, err = sandboxRepository.All()
	if err != nil {
		
	}
	fmt.Printf("\nAll sandboxes:\n")
	for _, sandbox := range all {
		fmt.Printf("sandbox: %+v\n", sandbox)
	}

}
