// Comprehensive example demonstrating all MySubaru Go client features
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/alex-savin/go-mysubaru"
	"github.com/alex-savin/go-mysubaru/config"
	"gopkg.in/yaml.v3"
)

// ConfigFile represents the structure of config.yaml
type ConfigFile struct {
	MySubaru config.MySubaru `yaml:"mysubaru"`
	Timezone string          `yaml:"timezone"`
	Logging  struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
		Source bool   `yaml:"source"`
	} `yaml:"logging"`
}

// loadConfig loads configuration from config.yaml file
func loadConfig(filename string) (*config.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configFile ConfigFile
	err = yaml.Unmarshal(data, &configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg := &config.Config{
		MySubaru: configFile.MySubaru,
	}

	// Set default logger if not configured
	if cfg.Logger == nil {
		cfg.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return cfg, nil
}

func main() {
	fmt.Println("=== MySubaru Go Client - Comprehensive Example ===")

	// Load configuration from config.yaml
	configFile := "config.yaml"
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "--help" || arg == "-h" {
			fmt.Println("MySubaru Go Client Example")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  ./example-app [config-file]")
			fmt.Println()
			fmt.Println("Arguments:")
			fmt.Println("  config-file    Path to YAML configuration file (default: config.yaml)")
			fmt.Println()
			fmt.Println("Examples:")
			fmt.Println("  ./example-app                          # Use config.yaml")
			fmt.Println("  ./example-app my-config.yaml          # Use custom config file")
			fmt.Println("  ./example-app --help                  # Show this help")
			os.Exit(0)
		}
		configFile = arg
	}

	fmt.Printf("📄 Loading configuration from: %s\n", configFile)
	cfg, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\nPlease copy config.sample.yaml to config.yaml and update with your credentials", err)
	}

	fmt.Printf("🔧 Configuration loaded successfully\n")
	fmt.Printf("   Username: %s\n", cfg.MySubaru.Credentials.Username)
	fmt.Printf("   Region: %s\n", cfg.MySubaru.Region)
	fmt.Printf("   Device: %s\n", cfg.MySubaru.Credentials.DeviceName)
	fmt.Println()

	client, err := mysubaru.New(cfg)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fmt.Println("🔑 Authenticating...")
	ok, authErr, requires2FA := client.Authenticate()
	if requires2FA {
		log.Fatal("Device not registered; 2FA/verification required before proceeding")
	}
	if !ok || authErr != nil {
		log.Fatalf("Authentication failed: %v", authErr)
	}
	fmt.Println("✅ Authentication successful")

	// Get vehicles
	fmt.Println("🚗 Retrieving vehicles...")
	vehicles, err := client.GetVehicles()
	if err != nil {
		log.Fatal("Failed to get vehicles:", err)
	}

	if len(vehicles) == 0 {
		log.Fatal("No vehicles found")
	}

	vehicle := vehicles[0]
	fmt.Printf("✅ Found %d vehicle(s). Using: %s (%s)\n\n", len(vehicles), vehicle.CarName, vehicle.Vin)

	// === BASIC VEHICLE OPERATIONS ===
	fmt.Println("=== BASIC VEHICLE OPERATIONS ===")

	// Get vehicle status
	fmt.Println("📊 Getting vehicle status...")
	err = vehicle.GetVehicleStatus()
	if err != nil {
		log.Printf("Failed to get vehicle status: %v", err)
	} else {
		fmt.Printf("✅ Vehicle status: Odometer: %d miles, Engine: %s\n", vehicle.Odometer.Miles, vehicle.EngineState)
	}

	// Get vehicle condition
	fmt.Println("🔍 Getting vehicle condition...")
	err = vehicle.GetVehicleCondition()
	if err != nil {
		log.Printf("Failed to get vehicle condition: %v", err)
	} else {
		fmt.Printf("✅ Vehicle condition: %d tires, %d doors, %d windows\n", len(vehicle.Tires), len(vehicle.Doors), len(vehicle.Windows))
	}

	// Get vehicle health
	fmt.Println("🏥 Getting vehicle health...")
	err = vehicle.GetVehicleHealth()
	if err != nil {
		log.Printf("Failed to get vehicle health: %v", err)
	} else {
		fmt.Printf("✅ Vehicle health: %d trouble codes\n", len(vehicle.Troubles))
	}

	// === REMOTE COMMANDS ===
	// fmt.Println("\n=== REMOTE COMMANDS ===")

	// Remote lock
	// fmt.Println("🔒 Executing remote lock...")
	// ch, err := vehicle.Lock()
	// if err != nil {
	// 	log.Printf("Remote lock failed: %v", err)
	// } else {
	// 	status := <-ch
	// 	fmt.Printf("✅ Remote lock result: %s\n", status)
	// }

	// Remote unlock
	// fmt.Println("🔓 Executing remote unlock...")
	// ch, err = vehicle.Unlock()
	// if err != nil {
	// 	log.Printf("Remote unlock failed: %v", err)
	// } else {
	// 	status := <-ch
	// 	fmt.Printf("✅ Remote unlock result: %s\n", status)
	// }

	// Horn and lights
	// fmt.Println("🚨 Executing horn and lights...")
	// ch, err = vehicle.HornStart()
	// if err != nil {
	// 	log.Printf("Horn and lights failed: %v", err)
	// } else {
	// 	status := <-ch
	// 	fmt.Printf("✅ Horn and lights result: %s\n", status)
	// }

	// === CLIMATE CONTROL ===
	fmt.Println("\n=== CLIMATE CONTROL ===")

	// Get climate presets
	fmt.Println("🌡️ Getting climate presets...")
	err = vehicle.GetClimatePresets()
	if err != nil {
		log.Printf("Failed to get climate presets: %v", err)
	} else {
		fmt.Printf("✅ Found %d climate presets\n", len(vehicle.ClimateProfiles))
	}

	// Get user presets
	fmt.Println("👤 Getting user climate presets...")
	err = vehicle.GetClimateUserPresets()
	if err != nil {
		log.Printf("Failed to get user presets: %v", err)
	} else {
		fmt.Printf("✅ Found %d user presets\n", len(vehicle.ClimateProfiles))
	}

	// Get quick presets
	fmt.Println("⚡ Getting quick presets...")
	err = vehicle.GetClimateQuickPresets()
	if err != nil {
		log.Printf("Failed to get quick presets: %v", err)
	} else {
		fmt.Printf("✅ Found %d quick presets\n", len(vehicle.ClimateProfiles))
	}

	// === EV CHARGING (if applicable) ===
	if vehicle.IsEV() {
		fmt.Println("\n=== EV CHARGING FEATURES ===")

		// Get EV charge settings
		fmt.Println("🔋 Getting EV charge settings...")
		err = vehicle.GetEVChargeSettings()
		if err != nil {
			log.Printf("Failed to get EV settings: %v", err)
		} else {
			fmt.Printf("✅ EV settings retrieved\n")
		}

		// Start charging
		fmt.Println("⚡ Starting EV charging...")
		ch, err := vehicle.ChargeOn()
		if err != nil {
			log.Printf("EV charging failed: %v", err)
		} else {
			status := <-ch
			fmt.Printf("✅ EV charging result: %s\n", status)
		}
	} else {
		fmt.Println("\n=== SKIPPING EV FEATURES ===")
		fmt.Println("ℹ️  Vehicle is not EV-capable, skipping EV charging features")
	}

	// // === ADVANCED G2 FEATURES ===
	// fmt.Println("\n=== ADVANCED G2 FEATURES ===")

	// // Check if vehicle supports G2 features
	// hasG2 := false
	// for _, feature := range vehicle.Features {
	// 	if feature == mysubaru.FEATURE_G2_TELEMATICS {
	// 		hasG2 = true
	// 		break
	// 	}
	// }

	// hasSafety := false
	// for _, feature := range vehicle.SubscriptionFeatures {
	// 	if feature == mysubaru.FEATURE_SAFETY {
	// 		hasSafety = true
	// 		break
	// 	}
	// }

	// if hasG2 && hasSafety {
	// 	fmt.Println("🛡️  Vehicle supports advanced G2 features")

	// 	// Set up geofence
	// 	fmt.Println("🌐 Setting up geofence...")
	// 	homeLat := 40.7128 // Example: New York City
	// 	homeLng := -74.0060
	// 	homeRadius := 500 // 500 meters

	// 	ch, err := vehicle.SetGeoFence(homeLat, homeLng, homeRadius, "Home", true, true, true)
	// 	if err != nil {
	// 		log.Printf("Geofence setup failed: %v", err)
	// 	} else {
	// 		status := <-ch
	// 		fmt.Printf("✅ Geofence setup result: %s\n", status)
	// 	}

	// 	// Set up speed fence
	// 	fmt.Println("🚦 Setting up speed fence...")
	// 	ch, err = vehicle.SetSpeedFence(65, true, true) // 65 mph limit
	// 	if err != nil {
	// 		log.Printf("Speed fence setup failed: %v", err)
	// 	} else {
	// 		status := <-ch
	// 		fmt.Printf("✅ Speed fence setup result: %s\n", status)
	// 	}

	// 	// Set up curfew
	// 	fmt.Println("🌙 Setting up curfew...")
	// 	daysOfWeek := []int{1, 2, 3, 4, 5} // Monday to Friday
	// 	ch, err = vehicle.SetCurfew("22:00", "06:00", daysOfWeek, true)
	// 	if err != nil {
	// 		log.Printf("Curfew setup failed: %v", err)
	// 	} else {
	// 		status := <-ch
	// 		fmt.Printf("✅ Curfew setup result: %s\n", status)
	// 	}

	// 	// Wait a moment before cleanup
	// 	fmt.Println("⏳ Waiting 2 seconds before cleanup...")
	// 	time.Sleep(2 * time.Second)

	// 	// Example cleanup (commented out to avoid accidental deletion)
	// 	/*
	// 	   fmt.Println("🗑️  Cleaning up geofence...")
	// 	   ch, err = vehicle.DeleteGeoFence("fence-id-here")
	// 	   if err != nil {
	// 	       log.Printf("Geofence deletion failed: %v", err)
	// 	   } else {
	// 	       status := <-ch
	// 	       fmt.Printf("✅ Geofence deletion result: %s\n", status)
	// 	   }
	// 	*/

	// } else {
	// 	fmt.Println("❌ Vehicle does not support advanced G2 features")
	// 	fmt.Printf("   G2 Telematics: %v\n", hasG2)
	// 	fmt.Printf("   Safety Plus: %v\n", hasSafety)
	// }

	// === CLEANUP ===
	fmt.Println("\n=== CLEANUP ===")
	fmt.Println("✅ Client operations completed successfully")

	fmt.Println("\n🎉 Comprehensive example completed!")
	fmt.Println("📚 Check the code comments for more details on each feature")
}
