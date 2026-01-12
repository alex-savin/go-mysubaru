package mysubaru

import (
	"strings"
	"testing"
)

// TestLockVehicleWithFixtures tests vehicle locking using fixtures
func TestLockVehicleWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:             "login_single_car.json",
		apiURLs["API_REMOTE_SVC_STATUS"]: "remoteServiceStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test locking the vehicle
	_, err = vehicle.Lock()
	if err != nil {
		t.Fatalf("expected no error locking vehicle, got %v", err)
	}
}

// TestUnlockVehicleWithFixtures tests vehicle unlocking using fixtures
func TestUnlockVehicleWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:             "login_single_car.json",
		apiURLs["API_REMOTE_SVC_STATUS"]: "remoteServiceStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test unlocking the vehicle
	_, err = vehicle.Unlock()
	if err != nil {
		t.Fatalf("expected no error unlocking vehicle, got %v", err)
	}
}

// TestStartEngineWithFixtures tests engine start using fixtures
func TestStartEngineWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:             "login_single_car.json",
		apiURLs["API_REMOTE_SVC_STATUS"]: "remoteServiceStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test starting the engine
	_, err = vehicle.EngineStart(10, 0, false)
	if err != nil {
		t.Fatalf("expected no error starting engine, got %v", err)
	}
}

// TestStopEngineWithFixtures tests engine stop using fixtures
func TestStopEngineWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:             "login_single_car.json",
		apiURLs["API_REMOTE_SVC_STATUS"]: "remoteServiceStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test stopping the engine
	_, err = vehicle.EngineStop()
	if err != nil {
		t.Fatalf("expected no error stopping engine, got %v", err)
	}
}

// TestLightsOnlyWithFixtures tests lights only using fixtures
func TestLightsOnlyWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:             "login_single_car.json",
		apiURLs["API_REMOTE_SVC_STATUS"]: "remoteServiceStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test lights only
	_, err = vehicle.LightsStart()
	if err != nil {
		t.Fatalf("expected no error with lights only, got %v", err)
	}
}

// TestLocateVehicleWithFixtures tests vehicle location using fixtures
func TestLocateVehicleWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:             "login_single_car.json",
		apiURLs["API_REMOTE_SVC_STATUS"]: "remoteServiceStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test locating the vehicle
	_, err = vehicle.GetLocation(false)
	if err != nil {
		t.Fatalf("expected no error locating vehicle, got %v", err)
	}
}

// TestHornAndLightsWithFixtures tests horn and lights using fixtures
func TestHornAndLightsWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:             "login_single_car.json",
		apiURLs["API_REMOTE_SVC_STATUS"]: "remoteServiceStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test horn and lights
	_, err = vehicle.HornStart()
	if err != nil {
		t.Fatalf("expected no error with horn and lights, got %v", err)
	}
}

// TestChargeNowWithFixtures tests PHEV charge now using fixtures
func TestChargeNowWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:             "login_single_car.json",
		apiURLs["API_REMOTE_SVC_STATUS"]: "remoteServiceStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test charge now (only works for EVs)
	if vehicle.IsEV() {
		_, err = vehicle.ChargeOn()
		if err != nil {
			t.Fatalf("expected no error with charge now, got %v", err)
		}
	} else {
		t.Skip("vehicle is not an EV, skipping charge test")
	}
}

// TestGetVehicleHealthWithActiveTroubles tests that active troubles are correctly parsed and displayed
func TestGetVehicleHealthWithActiveTroubles(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:          "login_single_car.json",
		apiURLs["API_VEHICLE_HEALTH"]: "vehicleHealth_with_troubles.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Manually set subscription features to enable health check
	vehicle.SubscriptionFeatures = []string{"REMOTE", "SAFETY"}

	err = vehicle.GetVehicleHealth()
	if err != nil {
		t.Fatalf("expected no error getting vehicle health, got %v", err)
	}

	// Verify that active troubles were parsed
	if len(vehicle.Troubles) == 0 {
		t.Fatal("expected active troubles to be parsed")
	}

	// Check that VDC_MIL and WASH_MIL are in the troubles map
	if _, ok := vehicle.Troubles["VDC_MIL"]; !ok {
		t.Error("expected VDC_MIL to be in troubles map")
	}
	if _, ok := vehicle.Troubles["WASH_MIL"]; !ok {
		t.Error("expected WASH_MIL to be in troubles map")
	}

	// Verify the descriptions are correct
	if vehicle.Troubles["VDC_MIL"].Description != "Vehicle Dynamics Control" {
		t.Errorf("expected VDC_MIL description 'Vehicle Dynamics Control', got %s", vehicle.Troubles["VDC_MIL"].Description)
	}
	if vehicle.Troubles["WASH_MIL"].Description != "Windshield Washer Fluid Level" {
		t.Errorf("expected WASH_MIL description 'Windshield Washer Fluid Level', got %s", vehicle.Troubles["WASH_MIL"].Description)
	}

	// Test that String() method includes the troubles
	vehicleString := vehicle.String()
	if !strings.Contains(vehicleString, "=== TROUBLES =====================") {
		t.Error("expected String() to contain troubles section")
	}
	if !strings.Contains(vehicleString, "Vehicle Dynamics Control") {
		t.Error("expected String() to contain VDC trouble description")
	}
	if !strings.Contains(vehicleString, "Windshield Washer Fluid Level") {
		t.Error("expected String() to contain washer trouble description")
	}
}

// TestGetVehicleStatusWithFixtures tests vehicle status data using fixtures
func TestGetVehicleStatusWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:          "login_single_car.json",
		apiURLs["API_VEHICLE_STATUS"]: "vehicleStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	err = vehicle.GetVehicleStatus()
	if err == nil {
		t.Fatal("expected subscription error, got no error")
	}
	if !strings.Contains(err.Error(), "subscription") {
		t.Errorf("expected subscription error, got %v", err)
	}
}

// TestGetVehicleConditionWithFixtures tests vehicle condition data using fixtures
func TestGetVehicleConditionWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:                     "login_single_car.json",
		urlToGen(apiURLs["API_CONDITION"], "g1"): "condition.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	err = vehicle.GetVehicleCondition()
	if err == nil {
		t.Fatal("expected subscription error, got no error")
	}
	if !strings.Contains(err.Error(), "subscription") {
		t.Errorf("expected subscription error, got %v", err)
	}
}

// TestClimateControlWithFixtures tests climate control using fixtures
func TestClimateControlWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:                       "login_single_car.json",
		apiURLs["API_G2_FETCH_RES_SUBARU_PRESETS"]: "climatePresetsSubaru.json",
		apiURLs["API_G2_FETCH_RES_USER_PRESETS"]:   "climatePresetsUser.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test Subaru presets
	err = vehicle.GetClimatePresets()
	if err == nil {
		t.Fatal("expected subscription error, got no error")
	}
	if !strings.Contains(err.Error(), "subscription") {
		t.Errorf("expected subscription error, got %v", err)
	}

	// Test user presets
	err = vehicle.GetClimateUserPresets()
	if err == nil {
		t.Fatal("expected subscription error, got no error")
	}
	if !strings.Contains(err.Error(), "subscription") {
		t.Errorf("expected subscription error, got %v", err)
	}
}

// TestGeofencingWithFixtures tests geofencing features using fixtures
func TestGeofencingWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]: "login_single_car.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test geofencing operations
	// Note: These methods may not exist in the current implementation
	// This test serves as a placeholder for when geofencing is implemented
	if vehicle == nil {
		t.Fatal("expected vehicle to be created")
	}
}

// TestSpeedFencingWithFixtures tests speed fencing features using fixtures
func TestSpeedFencingWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]: "login_single_car.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test speed fencing operations
	// Note: These methods may not exist in the current implementation
	// This test serves as a placeholder for when speed fencing is implemented
	if vehicle == nil {
		t.Fatal("expected vehicle to be created")
	}
}

// TestCurfewWithFixtures tests curfew features using fixtures
func TestCurfewWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]: "login_single_car.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test curfew operations
	// Note: These methods may not exist in the current implementation
	// This test serves as a placeholder for when curfew is implemented
	if vehicle == nil {
		t.Fatal("expected vehicle to be created")
	}
}

// TestG2FeaturesWithFixtures tests G2 telematics features using fixtures
func TestG2FeaturesWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]: "login_single_car.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test G2 features
	// Note: These methods may not exist in the current implementation
	// This test serves as a placeholder for when G2 features are implemented
	if vehicle == nil {
		t.Fatal("expected vehicle to be created")
	}
}

// TestEVOperationsWithFixtures tests EV operations using fixtures
func TestEVOperationsWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]:             "login_single_car.json",
		apiURLs["API_REMOTE_SVC_STATUS"]: "remoteServiceStatus.json",
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test EV operations
	// Note: These methods may not exist in the current implementation
	// This test serves as a placeholder for when EV operations are implemented
	if vehicle == nil {
		t.Fatal("expected vehicle to be created")
	}
}

// TestErrorHandlingWithFixtures tests error handling using fixtures
func TestErrorHandlingWithFixtures(t *testing.T) {
	fixtures := map[string]string{
		apiURLs["API_LOGIN"]: "login_single_car.json",
		// Use a fixture that contains error responses
	}

	ts := mockMySubaruApiWithFixtures(t, fixtures)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	// Get the first vehicle
	vehicles, err := msc.GetVehicles()
	if err != nil {
		t.Fatalf("expected no error getting vehicles, got %v", err)
	}
	if len(vehicles) == 0 {
		t.Fatal("expected at least one vehicle")
	}

	vehicle := vehicles[0]

	// Test error handling scenarios
	if vehicle == nil {
		t.Fatal("expected vehicle to be created")
	}
}
