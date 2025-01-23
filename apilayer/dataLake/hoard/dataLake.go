package hoard

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/earmuff-jam/fleetwise/db"
	"github.com/earmuff-jam/fleetwise/model"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type FItem struct {
	Name        string
	Description string
}

// IngestSvc ...
//
// main function to ingest data from the datalake
func IngestSvc() {

	// load environment variables
	err := godotenv.Load(filepath.Join("..", "..", ".env"))
	if err != nil {
		log.Printf("No env file detected. Using os system configuration.")
	}

	// setupDB
	currentUser := os.Getenv("CLIENT_USER")
	db, err := SetupDB(currentUser)
	if err != nil {
		log.Fatal("Unable to use data lake without a db source")
	}
	defer db.Close()

	emailAddress := "admin@gmail.com"
	// using the same user as test instance
	creatorID, err := GenerateFakeUser(currentUser, emailAddress)
	if err != nil {
		log.Fatalf("error during finding a test user... %+v", err)
	}
	GenerateFakeDataWithLimit(currentUser, 10, "note", creatorID)
	GenerateFakeDataWithLimit(currentUser, 10, "inventory", creatorID)
	GenerateFakeDataWithLimit(currentUser, 10, "category", creatorID)
	GenerateFakeDataWithLimit(currentUser, 10, "maintenance_plan", creatorID)
}

// GenerateFakeDataWithLimit ...
//
// GenerateFakeDataWithLimit retrieves fake data for mashed application with a limit of rows.
func GenerateFakeDataWithLimit(user string, limit int, typeOf string, creatorID string) error {

	if limit <= 0 {
		fmt.Println("unable to add fake data with empty limit")
		return errors.New("unable to add fake data")
	}

	// rowCounts are generated from limit. this is to generate dynamic information
	rowCounts := gofakeit.Number(1, limit)

	switch typeOf {
	case "note":
		log.Printf("loading %d of note rows", rowCounts)
		populateFakeNotes(user, rowCounts, creatorID)
	case "inventory":
		log.Printf("loading %d of inventory rows", rowCounts)
		populateFakeInventories(user, rowCounts, creatorID)
	case "category":
		log.Printf("loading %d of category rows", rowCounts)
		populateFakeCategories(user, rowCounts, creatorID)
	case "maintenance_plan":
		log.Printf("loading %d of maintenance plan rows", rowCounts)
		populateFakeMaintenancePlans(user, rowCounts, creatorID)
	default:
	}
	return nil
}

// populateFakeCategories ...
//
// generate fake category data to test in the application
func populateFakeCategories(user string, limit int, creatorID string) {
	fakeCategories := []FItem{
		{
			Name:        "Office Equipment",
			Description: "Track everything from desks and chairs to printers and photocopiers, ensuring your office runs smoothly and efficiently.",
		},
		{
			Name:        "IT Hardware",
			Description: "Monitor your tech inventory, including laptops, monitors, servers, and peripherals, to keep your digital workspace optimized.",
		},
		{
			Name:        "Vehicles & Fleet",
			Description: "Keep tabs on company vehicles, from delivery vans to executive cars, and never lose track of your fleet.",
		},
		{
			Name:        "Furniture & Fixtures",
			Description: "Manage office and workplace furnishings with ease, from ergonomic chairs to lighting setups.",
		},
		{
			Name:        "Tools & Machinery",
			Description: "Perfect for contractors and manufacturers, this category covers drills, saws, forklifts, and all the heavy-duty essentials.",
		},
		{
			Name:        "Inventory Supplies",
			Description: "Stay on top of stockroom items like packaging materials, spare parts, and consumables to avoid bottlenecks.",
		},
		{
			Name:        "Marketing Assets",
			Description: "From banners to branded merchandise, organize and track assets used for promotional and advertising campaigns.",
		},
		{
			Name:        "Maintenance Equipment",
			Description: "Track repair tools and cleaning supplies so you can keep your facilities in top shape.",
		},
		{
			Name:        "Software Licenses",
			Description: "Manage software subscriptions and licenses to avoid downtime and ensure compliance with usage agreements.",
		},
		{
			Name:        "Customer Loaned Items",
			Description: "Easily monitor items lent to clients, like demo products, rental gear, or loaner tools, to ensure timely returns.",
		},
	}

	for i := 1; i <= limit; i++ {
		now := time.Now()
		draftCategory := model.Category{}
		daysAgo := gofakeit.Number(1, 30)
		startDate := now.AddDate(0, 0, -daysAgo)

		randomIndex := gofakeit.Number(0, len(fakeCategories)-1)
		selectedCategory := fakeCategories[randomIndex]

		draftCategory.Name = selectedCategory.Name
		draftCategory.Description = selectedCategory.Description
		draftCategory.Color = gofakeit.RandString([]string{"#ffcc99", "#ff00ff", "#bb5588"})
		draftCategory.Status = "draft"
		draftCategory.CreatedAt = startDate
		draftCategory.UpdatedAt = startDate
		draftCategory.CreatedBy = creatorID
		draftCategory.Creator = gofakeit.FirstName()
		draftCategory.UpdatedBy = creatorID
		draftCategory.Updator = gofakeit.FirstName()
		draftCategory.SharableGroups = []string{creatorID}

		db.CreateCategory(user, &draftCategory)
	}
}

// populateFakeMaintenancePlans ...
//
// generate fake maintenance plan data to test in the application
func populateFakeMaintenancePlans(user string, limit int, creatorID string) {
	fakeMaintenancePlans := []FItem{
		{
			Name:        "Basic Equipment Checkup",
			Description: "Regular inspections for wear and tear, ensuring basic functionality of equipment.",
		},
		{
			Name:        "Preventive Maintenance",
			Description: "Scheduled maintenance to prevent equipment breakdown and extend lifespan.",
		},
		{
			Name:        "Seasonal Maintenance",
			Description: "Specialized checks and updates before seasonal usage or storage periods.",
		},
		{
			Name:        "Emergency Repair Plan",
			Description: "Rapid response services to address unexpected equipment failures.",
		},
		{
			Name:        "Software Update Plan",
			Description: "Regular updates to firmware and software to maintain optimal performance.",
		},
		{
			Name:        "Comprehensive Overhaul",
			Description: "Full equipment overhaul for restoring it to like-new condition.",
		},
		{
			Name:        "Warranty Support Plan",
			Description: "Maintenance services covered under manufacturer warranty agreements.",
		},
		{
			Name:        "Energy Efficiency Audit",
			Description: "Assess and optimize equipment energy usage to reduce operational costs.",
		},
		{
			Name:        "Safety Compliance Plan",
			Description: "Ensure all equipment meets regulatory safety and compliance standards.",
		},
		{
			Name:        "On-Demand Maintenance",
			Description: "Maintenance performed only as needed, based on real-time equipment conditions.",
		},
	}
	for i := 1; i <= limit; i++ {
		draftMaintenancePlan := model.MaintenancePlan{}
		// time is set for created / updated
		now := time.Now()
		daysAgo := gofakeit.Number(1, 30)
		startDate := now.AddDate(0, 0, -daysAgo)

		randomIndex := gofakeit.Number(0, len(fakeMaintenancePlans)-1)
		selectedCategory := fakeMaintenancePlans[randomIndex]

		draftMaintenancePlan.Name = selectedCategory.Name
		draftMaintenancePlan.Description = selectedCategory.Description
		draftMaintenancePlan.Color = gofakeit.RandString([]string{"#ffcc99", "#ff00ff", "#bb5588"})
		draftMaintenancePlan.Status = "draft"
		draftMaintenancePlan.PlanType = gofakeit.RandString([]string{"annual", "weekly", "daily"})
		draftMaintenancePlan.PlanDue = startDate
		draftMaintenancePlan.CreatedAt = startDate
		draftMaintenancePlan.UpdatedAt = startDate
		draftMaintenancePlan.CreatedBy = creatorID
		draftMaintenancePlan.Creator = gofakeit.FirstName()
		draftMaintenancePlan.UpdatedBy = creatorID
		draftMaintenancePlan.Updator = gofakeit.FirstName()
		draftMaintenancePlan.SharableGroups = []string{creatorID}

		db.CreateMaintenancePlan(user, &draftMaintenancePlan)
	}
}

// populateFakeNotes ...
//
// generate fake item data to test in the application
func populateFakeNotes(user string, limit int, creatorID string) {

	for i := 1; i <= limit; i++ {
		draftNote := model.Note{}

		// time is set for created / updated
		now := time.Now()
		daysAgo := gofakeit.Number(1, 30)
		startDate := now.AddDate(0, 0, -daysAgo)

		draftNote.Title = gofakeit.JobTitle()
		draftNote.Description = gofakeit.HipsterSentence(2)
		draftNote.CreatedAt = startDate
		draftNote.UpdatedAt = startDate
		draftNote.CreatedBy = creatorID
		draftNote.Creator = gofakeit.FirstName()
		draftNote.UpdatedBy = creatorID
		draftNote.Updator = gofakeit.FirstName()

		db.AddNewNote(user, creatorID, draftNote)
	}
}

// populateFakeInventories ...
//
// generate fake personal inventories data to test in the application
func populateFakeInventories(user string, limit int, creatorID string) {

	fakeAssets := []FItem{
		{
			Name:        "Advanced Diagnostics Tool",
			Description: "Comprehensive equipment diagnostics for identifying underlying issues and optimizing performance.",
		},
		{
			Name:        "Smart Monitoring System",
			Description: "Real-time monitoring and alerts for equipment status and performance.",
		},
		{
			Name:        "Portable Power Generator",
			Description: "Reliable backup power source for emergencies and remote operations.",
		},
		{
			Name:        "Heavy-Duty Hydraulic Press",
			Description: "Powerful and efficient press for industrial-grade applications.",
		},
		{
			Name:        "Precision Calibration Kit",
			Description: "Tools designed for accurate calibration of sensitive equipment.",
		},
		{
			Name:        "Industrial Cooling System",
			Description: "High-capacity cooling system to maintain optimal equipment temperatures.",
		},
		{
			Name:        "Ergonomic Workstation Setup",
			Description: "Customizable workstations designed for enhanced comfort and productivity.",
		},
		{
			Name:        "Wireless Communication Hub",
			Description: "Centralized hub for seamless wireless communication across devices.",
		},
		{
			Name:        "Autonomous Cleaning Robot",
			Description: "AI-driven cleaning robot for maintaining cleanliness in industrial spaces.",
		},
		{
			Name:        "High-Performance Filtration Unit",
			Description: "Advanced filtration system for removing contaminants and ensuring equipment longevity.",
		},
	}

	for i := 1; i <= limit; i++ {
		draftInventory := model.Inventory{}

		// time is set for created / updated
		now := time.Now()
		daysAgo := gofakeit.Number(1, 90)
		startDate := now.AddDate(0, 0, -daysAgo)
		isReturnableStatus := gofakeit.Bool()

		randomIndex := gofakeit.Number(0, len(fakeAssets)-1)
		selectedAsset := fakeAssets[randomIndex]

		draftInventory.Name = selectedAsset.Name
		draftInventory.Description = selectedAsset.Description
		draftInventory.Price = gofakeit.Price(2, 120)
		draftInventory.Status = gofakeit.RandString([]string{"COUPONS", "DRAFT", "HIDDEN", "ALL"})
		draftInventory.Barcode = fmt.Sprintf("%d,%v", gofakeit.Number(2, 20), gofakeit.BeerHop())
		draftInventory.SKU = fmt.Sprintf("%d,%v", gofakeit.Number(2, 20), gofakeit.BeerHop())
		draftInventory.Quantity = gofakeit.Number(10, 100)
		draftInventory.Location = gofakeit.RandString(
			[]string{"Master Bedroom Closet",
				"Garage",
				"Living Room Cabinet",
				"Bathroom Closet",
				"Dining Room Hutch",
				"Home Office Desk",
				"Basement Storage",
				"Kids'' Playroom",
				"Garage Workshop",
				"Guest Bedroom Closet",
				"Outdoor Shed"})
		draftInventory.IsReturnable = isReturnableStatus
		if isReturnableStatus {
			draftInventory.ReturnLocation = gofakeit.CarMaker()
		}
		draftInventory.MaxWeight = fmt.Sprintf("%d", gofakeit.Number(2, 5))
		draftInventory.MinWeight = fmt.Sprintf("%d", gofakeit.Number(2, 2))
		draftInventory.MaxHeight = fmt.Sprintf("%d", gofakeit.Number(2, 5))
		draftInventory.MinHeight = fmt.Sprintf("%d", gofakeit.Number(2, 2))
		draftInventory.CreatedAt = startDate
		draftInventory.UpdatedAt = startDate
		draftInventory.CreatedBy = creatorID
		draftInventory.CreatorName = gofakeit.FirstName()
		draftInventory.UpdatedBy = creatorID
		draftInventory.UpdaterName = gofakeit.FirstName()
		draftInventory.BoughtAt = gofakeit.CarMaker()

		db.AddInventory(user, creatorID, draftInventory)
	}
}

// GenerateFakeUser ...
//
// used to create a fake user for testing and ingest purpose only
func GenerateFakeUser(user string, emailAddress string) (string, error) {
	db, err := SetupDB(user)
	if err != nil {
		return "", err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	var updatedProfile model.Profile

	sqlStr := `
	UPDATE community.profiles SET 
	username = 'IngestSvcUser',
	full_name = 'John Doe',
	phone_number = '1234567890',
	about_me = 'I am an architect with a passion for creating functional and aesthetically pleasing spaces that inspire and serve their purpose. My approach to design is rooted in meticulous planning and a deep appreciation for organization.'
	WHERE email_address = $1
	RETURNING id;`

	row := tx.QueryRow(sqlStr, emailAddress)
	err = row.Scan(
		&updatedProfile.ID,
	)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	updatedProfile.Validate()
	return updatedProfile.ID.String(), nil
}

// SetupDB ...
//
// function is used to build the db for the data lake
func SetupDB(user string) (*sql.DB, error) {

	pwd := os.Getenv("CLIENT_PASSWORD")
	if len(pwd) == 0 {
		pwd = "password"
	}

	host := os.Getenv("DATABASE_DOCKER_CONTAINER_IP_ADDRESS")
	if len(host) == 0 {
		host = "localhost"
	}

	port := os.Getenv("DATABASE_DOCKER_CONTAINER_PORT")
	if len(port) == 0 {
		port = "5432"
	}

	psqlStr := fmt.Sprintf("postgres://%s:%s@%s:%s/community?sslmode=disable", user, pwd, host, port)
	var db, err = sql.Open("postgres", psqlStr)
	if err != nil {
		log.Fatalf("failed to open postgres db. error: +%v", err)
		return nil, err
	}

	// if the user is unable to ping the db, we don't want to submit the request
	err = db.Ping()
	if err != nil {
		fmt.Printf("unable to ping. error: +%v", err)
		return nil, err
	}
	return db, nil
}
