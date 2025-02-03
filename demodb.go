package main

import (
    "fmt"
    "time"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

type Sim struct {
    ID               uint      `json:"id" gorm:"primaryKey"`
    Name             string    `json:"name" gorm:"not null"`
    Number           string    `json:"number" gorm:"not null"`
    LastRechargeDate time.Time `json:"lastRechargeDate" gorm:"not null"`
    RechargeValidity time.Time `json:"rechargeValidity" gorm:"not null"`
    IncomingValidity time.Time `json:"incomingValidity" gorm:"not null"`
    SimExpiry        time.Time `json:"simExpiry" gorm:"not null"`
}

func main() {
    // Database connection string
    dsn := ""  //write your postgresql 
    
    // Open database connection
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database: " + err.Error())
    }
    fmt.Println("Successfully connected to database!")

    // Auto migrate the schema
    err = db.AutoMigrate(&Sim{})
    if err != nil {
        panic("Failed to migrate database: " + err.Error())
    }
    fmt.Println("Successfully migrated database schema!")

    // Test inserting a record
    testSim := Sim{
        Name:             "Test SIM",
        Number:           "1234567890",
        LastRechargeDate: time.Now(),
        RechargeValidity: time.Now().AddDate(0, 0, 28),
        IncomingValidity: time.Now().AddDate(0, 0, 35),
        SimExpiry:        time.Now().AddDate(0, 0, 120),
    }

    result := db.Create(&testSim)
    if result.Error != nil {
        panic("Failed to insert test record: " + result.Error.Error())
    }
    fmt.Printf("Successfully inserted test record with ID: %d\n", testSim.ID)

    // Test retrieving records
    var sims []Sim
    result = db.Find(&sims)
    if result.Error != nil {
        panic("Failed to retrieve records: " + result.Error.Error())
    }
    fmt.Printf("Found %d records in database\n", len(sims))
    
    // Print all records
    for _, sim := range sims {
        fmt.Printf("ID: %d, Name: %s, Number: %s, Last Recharge: %s\n",
            sim.ID, sim.Name, sim.Number, sim.LastRechargeDate.Format("2006-01-02"))
    }
}
