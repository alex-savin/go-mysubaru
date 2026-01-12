package mysubaru

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// TestGetClimateQuickPresets_Success tests the retrieval of quick climate presets.
func TestGetClimatePresets_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Handle API_LOGIN endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_LOGIN"] && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":"BIOMETRICS_DISABLED","dataName":"sessionData","data":{"sessionChanged":false,"vehicleInactivated":false,"account":{"marketId":1,"createdDate":1476984644000,"firstName":"Tatiana","lastName":"Savin","zipCode":"07974","accountKey":765268,"lastLoginDate":1751738613000,"zipCode5":"07974"},"resetPassword":false,"deviceId":"JddMBQXvAkgutSmEP6uFsThbq4QgEBBQ","sessionId":"9D7FCDF274794346689D3FA0D693CBBF","deviceRegistered":true,"passwordToken":null,"vehicles":[{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":null,"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":null,"authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":null,"subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":null,"sunsetUpgraded":true,"intDescrip":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}],"rightToRepairEnabled":true,"rightToRepairStartYear":2022,"rightToRepairStates":"MA","enableXtime":true,"termsAndConditionsAccepted":true,"digitalGlobeConnectId":"0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeImageTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/DigitalGlobe:ImageryTileService@EPSG:3857@png/{z}/{x}/{y}.png?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeTransparentTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/Digitalglobe:OSMTransparentTMSTileService@EPSG:3857@png/{z}/{x}/{-y}.png/?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","tomtomKey":"DHH9SwEQ4MW55Hj2TfqMeldbsDjTdgAs","currentVehicleIndex":0,"handoffToken":"$2a$08$rOb/uqhm8I3QtSel2phOCOxNM51w43eqXDDksMkJ.1a5KsaQuLvEu$1751745334477","satelliteViewEnabled":true,"registeredDevicePermanent":true}}`)
		}
		// Handle API_VALIDATE_SESSION endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_VALIDATE_SESSION"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":null}`)
		}
		// Handle SELECT_VEHICLE endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_SELECT_VEHICLE"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"vehicle","data":{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":"2023","modelCode":"PDL","engineSize":2.4,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"8KV8","licensePlateState":"NJ","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":["REMOTE","SAFETY","Retail3"],"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":["ABS_MIL","ACCS","AHBL_MIL","ATF_MIL","AWD_MIL","BSD","BSDRCT_MIL","CEL_MIL","CP1_5HHU","EBD_MIL","EOL_MIL","EPAS_MIL","EPB_MIL","ESS_MIL","EYESIGHT","ISS_MIL","MOONSTAT","OPL_MIL","PANPM-TUIRWAOC","PWAAADWWAP","RAB_MIL","RCC","REARBRK","RES","RESCC","RES_HVAC_HFS","RES_HVAC_VFS","RHSF","RPOI","RPOIA","RTGU","RVFS","SRH_MIL","SRS_MIL","SXM360L","T23DCM","TEL_MIL","TIF_35","TIR_33","TLD","TPMS_MIL","VALET","VDC_MIL","WASH_MIL","WDWSTAT","g3"],"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":"ACTIVE","authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":"Outback","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":"Cosmic Blue Pearl","sunsetUpgraded":true,"intDescrip":"Black","transCode":"CVT","provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}}`)
		}
		// Handle API_G2_FETCH_RES_SUBARU_PRESETS endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_G2_FETCH_RES_SUBARU_PRESETS"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":["{\"name\": \"Auto\", \"runTimeMinutes\": \"10\", \"climateZoneFrontTemp\": \"74\", \"climateZoneFrontAirMode\": \"AUTO\", \"climateZoneFrontAirVolume\": \"AUTO\", \"outerAirCirculation\": \"auto\", \"heatedRearWindowActive\": \"false\", \"airConditionOn\": \"false\", \"heatedSeatFrontLeft\": \"off\", \"heatedSeatFrontRight\": \"off\", \"startConfiguration\": \"START_ENGINE_ALLOW_KEY_IN_IGNITION\", \"canEdit\": \"true\", \"disabled\": \"false\", \"vehicleType\": \"gas\", \"presetType\": \"subaruPreset\" }","{\"name\":\"Full Cool\",\"runTimeMinutes\":\"10\",\"climateZoneFrontTemp\":\"60\",\"climateZoneFrontAirMode\":\"feet_face_balanced\",\"climateZoneFrontAirVolume\":\"7\",\"airConditionOn\":\"true\",\"heatedSeatFrontLeft\":\"high_cool\",\"heatedSeatFrontRight\":\"high_cool\",\"heatedRearWindowActive\":\"false\",\"outerAirCirculation\":\"outsideAir\",\"startConfiguration\":\"START_ENGINE_ALLOW_KEY_IN_IGNITION\",\"canEdit\":\"true\",\"disabled\":\"true\",\"vehicleType\":\"gas\",\"presetType\":\"subaruPreset\"}","{\"name\": \"Full Heat\", \"runTimeMinutes\": \"10\", \"climateZoneFrontTemp\": \"85\", \"climateZoneFrontAirMode\": \"feet_window\", \"climateZoneFrontAirVolume\": \"7\", \"airConditionOn\": \"false\", \"heatedSeatFrontLeft\": \"high_heat\", \"heatedSeatFrontRight\": \"high_heat\", \"heatedRearWindowActive\": \"true\", \"outerAirCirculation\": \"outsideAir\", \"startConfiguration\": \"START_ENGINE_ALLOW_KEY_IN_IGNITION\", \"canEdit\": \"true\", \"disabled\": \"true\", \"vehicleType\": \"gas\", \"presetType\": \"subaruPreset\" }","{\"name\": \"Full Cool\", \"runTimeMinutes\": \"10\", \"climateZoneFrontTemp\": \"60\", \"climateZoneFrontAirMode\": \"feet_face_balanced\", \"climateZoneFrontAirVolume\": \"7\", \"airConditionOn\": \"true\", \"heatedSeatFrontLeft\": \"OFF\", \"heatedSeatFrontRight\": \"OFF\", \"heatedRearWindowActive\": \"false\", \"outerAirCirculation\": \"outsideAir\", \"startConfiguration\": \"START_CLIMATE_CONTROL_ONLY_ALLOW_KEY_IN_IGNITION\", \"canEdit\": \"true\", \"disabled\": \"true\", \"vehicleType\": \"phev\", \"presetType\": \"subaruPreset\" }","{\"name\": \"Full Heat\", \"runTimeMinutes\": \"10\", \"climateZoneFrontTemp\": \"85\", \"climateZoneFrontAirMode\": \"feet_window\", \"climateZoneFrontAirVolume\": \"7\", \"airConditionOn\": \"false\", \"heatedSeatFrontLeft\": \"high_heat\", \"heatedSeatFrontRight\": \"high_heat\", \"heatedRearWindowActive\": \"true\", \"outerAirCirculation\": \"outsideAir\", \"startConfiguration\": \"START_CLIMATE_CONTROL_ONLY_ALLOW_KEY_IN_IGNITION\", \"canEdit\": \"true\", \"disabled\": \"true\", \"vehicleType\": \"phev\", \"presetType\": \"subaruPreset\" }"]}`)
		}
	}

	ts := mockMySubaruApi(t, handler)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msc == nil {
		t.Fatalf("expected MySubaru API client, got nil")
	}

	_, authErr, _ := msc.Authenticate()
	if authErr != nil {
		t.Fatalf("expected no error, got %v", authErr)
	}
	if !msc.isAuthenticated || !msc.isRegistered {
		t.Errorf("expected authenticated and registered true, got %v %v", msc.isAuthenticated, msc.isRegistered)
	}
	if msc.currentVin != "1HGCM82633A004352" {
		t.Errorf("expected currentVin 1HGCM82633A004352, got %v", msc.currentVin)
	}
}

// TestGetClimateQuickPresets_Success tests the retrieval of quick climate presets.
func TestGetClimateQuickPresets_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Handle API_LOGIN endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_LOGIN"] && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":"BIOMETRICS_DISABLED","dataName":"sessionData","data":{"sessionChanged":false,"vehicleInactivated":false,"account":{"marketId":1,"createdDate":1476984644000,"firstName":"Tatiana","lastName":"Savin","zipCode":"07974","accountKey":765268,"lastLoginDate":1751738613000,"zipCode5":"07974"},"resetPassword":false,"deviceId":"JddMBQXvAkgutSmEP6uFsThbq4QgEBBQ","sessionId":"9D7FCDF274794346689D3FA0D693CBBF","deviceRegistered":true,"passwordToken":null,"vehicles":[{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":null,"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":null,"authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":null,"subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":null,"sunsetUpgraded":true,"intDescrip":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}],"rightToRepairEnabled":true,"rightToRepairStartYear":2022,"rightToRepairStates":"MA","enableXtime":true,"termsAndConditionsAccepted":true,"digitalGlobeConnectId":"0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeImageTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/DigitalGlobe:ImageryTileService@EPSG:3857@png/{z}/{x}/{y}.png?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeTransparentTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/Digitalglobe:OSMTransparentTMSTileService@EPSG:3857@png/{z}/{x}/{-y}.png/?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","tomtomKey":"DHH9SwEQ4MW55Hj2TfqMeldbsDjTdgAs","currentVehicleIndex":0,"handoffToken":"$2a$08$rOb/uqhm8I3QtSel2phOCOxNM51w43eqXDDksMkJ.1a5KsaQuLvEu$1751745334477","satelliteViewEnabled":true,"registeredDevicePermanent":true}}`)
		}
		// Handle API_VALIDATE_SESSION endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_VALIDATE_SESSION"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":null}`)
		}
		// Handle SELECT_VEHICLE endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_SELECT_VEHICLE"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"vehicle","data":{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":"2023","modelCode":"PDL","engineSize":2.4,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"8KV8","licensePlateState":"NJ","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":["REMOTE","SAFETY","Retail3"],"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":["ABS_MIL","ACCS","AHBL_MIL","ATF_MIL","AWD_MIL","BSD","BSDRCT_MIL","CEL_MIL","CP1_5HHU","EBD_MIL","EOL_MIL","EPAS_MIL","EPB_MIL","ESS_MIL","EYESIGHT","ISS_MIL","MOONSTAT","OPL_MIL","PANPM-TUIRWAOC","PWAAADWWAP","RAB_MIL","RCC","REARBRK","RES","RESCC","RES_HVAC_HFS","RES_HVAC_VFS","RHSF","RPOI","RPOIA","RTGU","RVFS","SRH_MIL","SRS_MIL","SXM360L","T23DCM","TEL_MIL","TIF_35","TIR_33","TLD","TPMS_MIL","VALET","VDC_MIL","WASH_MIL","WDWSTAT","g3"],"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":"ACTIVE","authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":"Outback","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":"Cosmic Blue Pearl","sunsetUpgraded":true,"intDescrip":"Black","transCode":"CVT","provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}}`)
		}
		// Handle API_G2_FETCH_RES_QUICK_START_SETTINGS endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_G2_FETCH_RES_QUICK_START_SETTINGS"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":"{\"name\":\"Cooling\",\"runTimeMinutes\":\"10\",\"climateZoneFrontTemp\":\"65\",\"climateZoneFrontAirMode\":\"FEET_FACE_BALANCED\",\"climateZoneFrontAirVolume\":\"7\",\"outerAirCirculation\":\"outsideAir\",\"heatedRearWindowActive\":\"false\",\"heatedSeatFrontLeft\":\"HIGH_COOL\",\"airConditionOn\":\"false\",\"canEdit\":\"true\",\"disabled\":\"false\",\"presetType\":\"userPreset\",\"startConfiguration\":\"START_ENGINE_ALLOW_KEY_IN_IGNITION\"}"}`)
		}
	}

	ts := mockMySubaruApi(t, handler)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msc == nil {
		t.Fatalf("expected MySubaru API client, got nil")
	}

	_, authErr, _ := msc.Authenticate()
	if authErr != nil {
		t.Fatalf("expected no error, got %v", authErr)
	}
	if !msc.isAuthenticated || !msc.isRegistered {
		t.Errorf("expected authenticated and registered true, got %v %v", msc.isAuthenticated, msc.isRegistered)
	}
	if msc.currentVin != "1HGCM82633A004352" {
		t.Errorf("expected currentVin 1HGCM82633A004352, got %v", msc.currentVin)
	}
}

// TestGetClimateUserPresets_Success tests the retrieval of user-defined climate presets.
func TestGetClimateUserPresets_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Handle API_LOGIN endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_LOGIN"] && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":"BIOMETRICS_DISABLED","dataName":"sessionData","data":{"sessionChanged":false,"vehicleInactivated":false,"account":{"marketId":1,"createdDate":1476984644000,"firstName":"Tatiana","lastName":"Savin","zipCode":"07974","accountKey":765268,"lastLoginDate":1751738613000,"zipCode5":"07974"},"resetPassword":false,"deviceId":"JddMBQXvAkgutSmEP6uFsThbq4QgEBBQ","sessionId":"9D7FCDF274794346689D3FA0D693CBBF","deviceRegistered":true,"passwordToken":null,"vehicles":[{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":null,"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":null,"authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":null,"subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":null,"sunsetUpgraded":true,"intDescrip":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}],"rightToRepairEnabled":true,"rightToRepairStartYear":2022,"rightToRepairStates":"MA","enableXtime":true,"termsAndConditionsAccepted":true,"digitalGlobeConnectId":"0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeImageTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/DigitalGlobe:ImageryTileService@EPSG:3857@png/{z}/{x}/{y}.png?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeTransparentTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/Digitalglobe:OSMTransparentTMSTileService@EPSG:3857@png/{z}/{x}/{-y}.png/?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","tomtomKey":"DHH9SwEQ4MW55Hj2TfqMeldbsDjTdgAs","currentVehicleIndex":0,"handoffToken":"$2a$08$rOb/uqhm8I3QtSel2phOCOxNM51w43eqXDDksMkJ.1a5KsaQuLvEu$1751745334477","satelliteViewEnabled":true,"registeredDevicePermanent":true}}`)
		}
		// Handle API_VALIDATE_SESSION endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_VALIDATE_SESSION"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":null}`)
		}
		// Handle SELECT_VEHICLE endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_SELECT_VEHICLE"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"vehicle","data":{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":"2023","modelCode":"PDL","engineSize":2.4,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"8KV8","licensePlateState":"NJ","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":["REMOTE","SAFETY","Retail3"],"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":["ABS_MIL","ACCS","AHBL_MIL","ATF_MIL","AWD_MIL","BSD","BSDRCT_MIL","CEL_MIL","CP1_5HHU","EBD_MIL","EOL_MIL","EPAS_MIL","EPB_MIL","ESS_MIL","EYESIGHT","ISS_MIL","MOONSTAT","OPL_MIL","PANPM-TUIRWAOC","PWAAADWWAP","RAB_MIL","RCC","REARBRK","RES","RESCC","RES_HVAC_HFS","RES_HVAC_VFS","RHSF","RPOI","RPOIA","RTGU","RVFS","SRH_MIL","SRS_MIL","SXM360L","T23DCM","TEL_MIL","TIF_35","TIR_33","TLD","TPMS_MIL","VALET","VDC_MIL","WASH_MIL","WDWSTAT","g3"],"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":"ACTIVE","authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":"Outback","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":"Cosmic Blue Pearl","sunsetUpgraded":true,"intDescrip":"Black","transCode":"CVT","provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}}`)
		}
		// Handle API_G2_FETCH_RES_USER_PRESETS endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_G2_FETCH_RES_USER_PRESETS"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":"[{\"name\":\"Cooling\",\"runTimeMinutes\":\"10\",\"climateZoneFrontTemp\":\"65\",\"climateZoneFrontAirMode\":\"FEET_FACE_BALANCED\",\"climateZoneFrontAirVolume\":\"7\",\"outerAirCirculation\":\"outsideAir\",\"heatedRearWindowActive\":\"false\",\"heatedSeatFrontLeft\":\"HIGH_COOL\",\"heatedSeatFrontRight\":\"HIGH_COOL\",\"airConditionOn\":\"false\",\"canEdit\":\"true\",\"disabled\":\"false\",\"presetType\":\"userPreset\",\"startConfiguration\":\"START_ENGINE_ALLOW_KEY_IN_IGNITION\"}]"}`)
		}
	}

	ts := mockMySubaruApi(t, handler)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msc == nil {
		t.Fatalf("expected MySubaru API client, got nil")
	}

	_, authErr, _ := msc.Authenticate()
	if authErr != nil {
		t.Fatalf("expected no error, got %v", authErr)
	}
	if !msc.isAuthenticated || !msc.isRegistered {
		t.Errorf("expected authenticated and registered true, got %v %v", msc.isAuthenticated, msc.isRegistered)
	}
	if msc.currentVin != "1HGCM82633A004352" {
		t.Errorf("expected currentVin 1HGCM82633A004352, got %v", msc.currentVin)
	}
}

// TestGetVehicleCondition_Success tests the GetVehicleCondition method
func TestGetVehicleStatus_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Handle API_LOGIN endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_LOGIN"] && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":"BIOMETRICS_DISABLED","dataName":"sessionData","data":{"sessionChanged":false,"vehicleInactivated":false,"account":{"marketId":1,"createdDate":1476984644000,"firstName":"Tatiana","lastName":"Savin","zipCode":"07974","accountKey":765268,"lastLoginDate":1751738613000,"zipCode5":"07974"},"resetPassword":false,"deviceId":"JddMBQXvAkgutSmEP6uFsThbq4QgEBBQ","sessionId":"9D7FCDF274794346689D3FA0D693CBBF","deviceRegistered":true,"passwordToken":null,"vehicles":[{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":null,"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":null,"authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":null,"subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":null,"sunsetUpgraded":true,"intDescrip":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}],"rightToRepairEnabled":true,"rightToRepairStartYear":2022,"rightToRepairStates":"MA","enableXtime":true,"termsAndConditionsAccepted":true,"digitalGlobeConnectId":"0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeImageTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/DigitalGlobe:ImageryTileService@EPSG:3857@png/{z}/{x}/{y}.png?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeTransparentTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/Digitalglobe:OSMTransparentTMSTileService@EPSG:3857@png/{z}/{x}/{-y}.png/?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","tomtomKey":"DHH9SwEQ4MW55Hj2TfqMeldbsDjTdgAs","currentVehicleIndex":0,"handoffToken":"$2a$08$rOb/uqhm8I3QtSel2phOCOxNM51w43eqXDDksMkJ.1a5KsaQuLvEu$1751745334477","satelliteViewEnabled":true,"registeredDevicePermanent":true}}`)
		}
		// Handle API_VALIDATE_SESSION endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_VALIDATE_SESSION"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":null}`)
		}
		// Handle SELECT_VEHICLE endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_SELECT_VEHICLE"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"vehicle","data":{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":"2023","modelCode":"PDL","engineSize":2.4,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"8KV8","licensePlateState":"NJ","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":["REMOTE","SAFETY","Retail3"],"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":["ABS_MIL","ACCS","AHBL_MIL","ATF_MIL","AWD_MIL","BSD","BSDRCT_MIL","CEL_MIL","CP1_5HHU","EBD_MIL","EOL_MIL","EPAS_MIL","EPB_MIL","ESS_MIL","EYESIGHT","ISS_MIL","MOONSTAT","OPL_MIL","PANPM-TUIRWAOC","PWAAADWWAP","RAB_MIL","RCC","REARBRK","RES","RESCC","RES_HVAC_HFS","RES_HVAC_VFS","RHSF","RPOI","RPOIA","RTGU","RVFS","SRH_MIL","SRS_MIL","SXM360L","T23DCM","TEL_MIL","TIF_35","TIR_33","TLD","TPMS_MIL","VALET","VDC_MIL","WASH_MIL","WDWSTAT","g3"],"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":"ACTIVE","authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":"Outback","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":"Cosmic Blue Pearl","sunsetUpgraded":true,"intDescrip":"Black","transCode":"CVT","provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}}`)
		}
		// Handle API_VEHICLE_STATUS endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_VEHICLE_STATUS"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":{"vhsId":14662115789,"odometerValue":31694,"odometerValueKilometers":50996,"eventDate":1751742945000,"eventDateStr":"2025-07-05T19:15+0000","eventDateCarUser":1751742945000,"eventDateStrCarUser":"2025-07-05T19:15+0000","latitude":40.700153,"longitude":-74.401405,"positionHeadingDegree":"154","tirePressureFrontLeft":"2482","tirePressureFrontRight":"2482","tirePressureRearLeft":"2413","tirePressureRearRight":"2482","tirePressureFrontLeftPsi":"36","tirePressureFrontRightPsi":"36","tirePressureRearLeftPsi":"35","tirePressureRearRightPsi":"36","doorBootPosition":"CLOSED","doorEngineHoodPosition":"CLOSED","doorFrontLeftPosition":"CLOSED","doorFrontRightPosition":"CLOSED","doorRearLeftPosition":"CLOSED","doorRearRightPosition":"CLOSED","doorBootLockStatus":"LOCKED","doorFrontLeftLockStatus":"LOCKED","doorFrontRightLockStatus":"LOCKED","doorRearLeftLockStatus":"LOCKED","doorRearRightLockStatus":"LOCKED","distanceToEmptyFuelMiles":259.73,"distanceToEmptyFuelKilometers":418,"avgFuelConsumptionMpg":102.2,"avgFuelConsumptionLitersPer100Kilometers":2.3,"evStateOfChargePercent":null,"evDistanceToEmptyMiles":null,"evDistanceToEmptyKilometers":null,"evDistanceToEmptyByStateMiles":null,"evDistanceToEmptyByStateKilometers":null,"vehicleStateType":"IGNITION_OFF","windowFrontLeftStatus":"CLOSE","windowFrontRightStatus":"CLOSE","windowRearLeftStatus":"CLOSE","windowRearRightStatus":"CLOSE","windowSunroofStatus":"CLOSE","tyreStatusFrontLeft":"UNKNOWN","tyreStatusFrontRight":"UNKNOWN","tyreStatusRearLeft":"UNKNOWN","tyreStatusRearRight":"UNKNOWN","remainingFuelPercent":90,"distanceToEmptyFuelMiles10s":260,"distanceToEmptyFuelKilometers10s":420}}`)
		}
	}

	ts := mockMySubaruApi(t, handler)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msc == nil {
		t.Fatalf("expected MySubaru API client, got nil")
	}

	_, authErr, _ := msc.Authenticate()
	if authErr != nil {
		t.Fatalf("expected no error, got %v", authErr)
	}
	if !msc.isAuthenticated || !msc.isRegistered {
		t.Errorf("expected authenticated and registered true, got %v %v", msc.isAuthenticated, msc.isRegistered)
	}
	if msc.currentVin != "1HGCM82633A004352" {
		t.Errorf("expected currentVin 1HGCM82633A004352, got %v", msc.currentVin)
	}
}

// TestGetVehicleCondition_Success tests the GetVehicleCondition method
func TestGetVehicleCondition_Success(t *testing.T) {
	routes := []endpointRoute{
		{Method: http.MethodPost, Path: apiURLs["API_LOGIN"], Response: testLoginResponse},
		{Method: http.MethodGet, Path: apiURLs["API_VALIDATE_SESSION"], Response: testValidateSessionResponse},
		{Method: http.MethodGet, Path: apiURLs["API_SELECT_VEHICLE"], Response: testSelectVehicleResponse},
		{Method: http.MethodGet, Path: apiURLs["API_CONDITION"], Response: testConditionResponse},
	}

	ts := mockServerWithRoutes(t, routes)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msc == nil {
		t.Fatalf("expected MySubaru API client, got nil")
	}

	_, authErr, _ := msc.Authenticate()
	if authErr != nil {
		t.Fatalf("expected no error, got %v", authErr)
	}
	if !msc.isAuthenticated || !msc.isRegistered {
		t.Errorf("expected authenticated and registered true, got %v %v", msc.isAuthenticated, msc.isRegistered)
	}
	if msc.currentVin != "1HGCM82633A004352" {
		t.Errorf("expected currentVin 1HGCM82633A004352, got %v", msc.currentVin)
	}
}

// TestGetVehicleHealth_Success tests the successful retrieval of vehicle health data.
func TestGetVehicleHealth_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Handle API_LOGIN endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_LOGIN"] && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":"BIOMETRICS_DISABLED","dataName":"sessionData","data":{"sessionChanged":false,"vehicleInactivated":false,"account":{"marketId":1,"createdDate":1476984644000,"firstName":"Tatiana","lastName":"Savin","zipCode":"07974","accountKey":765268,"lastLoginDate":1751738613000,"zipCode5":"07974"},"resetPassword":false,"deviceId":"JddMBQXvAkgutSmEP6uFsThbq4QgEBBQ","sessionId":"9D7FCDF274794346689D3FA0D693CBBF","deviceRegistered":true,"passwordToken":null,"vehicles":[{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":null,"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":null,"authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":null,"subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":null,"sunsetUpgraded":true,"intDescrip":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}],"rightToRepairEnabled":true,"rightToRepairStartYear":2022,"rightToRepairStates":"MA","enableXtime":true,"termsAndConditionsAccepted":true,"digitalGlobeConnectId":"0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeImageTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/DigitalGlobe:ImageryTileService@EPSG:3857@png/{z}/{x}/{y}.png?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeTransparentTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/Digitalglobe:OSMTransparentTMSTileService@EPSG:3857@png/{z}/{x}/{-y}.png/?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","tomtomKey":"DHH9SwEQ4MW55Hj2TfqMeldbsDjTdgAs","currentVehicleIndex":0,"handoffToken":"$2a$08$rOb/uqhm8I3QtSel2phOCOxNM51w43eqXDDksMkJ.1a5KsaQuLvEu$1751745334477","satelliteViewEnabled":true,"registeredDevicePermanent":true}}`)
		}
		// Handle API_VALIDATE_SESSION endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_VALIDATE_SESSION"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":null}`)
		}
		// Handle SELECT_VEHICLE endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_SELECT_VEHICLE"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"vehicle","data":{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":"2023","modelCode":"PDL","engineSize":2.4,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"8KV8","licensePlateState":"NJ","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":["REMOTE","SAFETY","Retail3"],"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":["ABS_MIL","ACCS","AHBL_MIL","ATF_MIL","AWD_MIL","BSD","BSDRCT_MIL","CEL_MIL","CP1_5HHU","EBD_MIL","EOL_MIL","EPAS_MIL","EPB_MIL","ESS_MIL","EYESIGHT","ISS_MIL","MOONSTAT","OPL_MIL","PANPM-TUIRWAOC","PWAAADWWAP","RAB_MIL","RCC","REARBRK","RES","RESCC","RES_HVAC_HFS","RES_HVAC_VFS","RHSF","RPOI","RPOIA","RTGU","RVFS","SRH_MIL","SRS_MIL","SXM360L","T23DCM","TEL_MIL","TIF_35","TIR_33","TLD","TPMS_MIL","VALET","VDC_MIL","WASH_MIL","WDWSTAT","g3"],"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":"ACTIVE","authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":"Outback","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":"Cosmic Blue Pearl","sunsetUpgraded":true,"intDescrip":"Black","transCode":"CVT","provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}}`)
		}
		// Handle API_VEHICLE_HEALTH endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_VEHICLE_HEALTH"] && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":{"lastUpdatedDate":1751742945000,"vehicleHealthItems":[{"warningCode":10,"b2cCode":"airbag","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"SRS_MIL"},{"warningCode":4,"b2cCode":"oilTemp","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"ATF_MIL"},{"warningCode":39,"b2cCode":"blindspot","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"BSDRCT_MIL"},{"warningCode":2,"b2cCode":"engineFail","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"CEL_MIL"},{"warningCode":44,"b2cCode":"pkgBrake","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"EPB_MIL"},{"warningCode":8,"b2cCode":"ebd","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"EBD_MIL"},{"warningCode":3,"b2cCode":"oilWarning","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"EOL_MIL"},{"warningCode":1,"b2cCode":"washer","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"WASH_MIL"},{"warningCode":50,"b2cCode":"iss","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"ISS_MIL"},{"warningCode":53,"b2cCode":"oilPres","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"OPL_MIL"},{"warningCode":11,"b2cCode":"epas","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"EPAS_MIL"},{"warningCode":69,"b2cCode":"revBrake","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"RAB_MIL"},{"warningCode":14,"b2cCode":"telematics","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"TEL_MIL"},{"warningCode":9,"b2cCode":"tpms","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"TPMS_MIL"},{"warningCode":7,"b2cCode":"vdc","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"VDC_MIL"},{"warningCode":6,"b2cCode":"abs","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"ABS_MIL"},{"warningCode":5,"b2cCode":"awd","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"AWD_MIL"},{"warningCode":12,"b2cCode":"eyesight","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"ESS_MIL"},{"warningCode":30,"b2cCode":"ahbl","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"AHBL_MIL"},{"warningCode":31,"b2cCode":"srh","isTrouble":false,"onDates":[],"onDaiId":0,"featureCode":"SRH_MIL"}]}}`)
		}
	}

	ts := mockMySubaruApi(t, handler)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msc == nil {
		t.Fatalf("expected MySubaru API client, got nil")
	}

	_, authErr, _ := msc.Authenticate()
	if authErr != nil {
		t.Fatalf("expected no error, got %v", authErr)
	}
	if !msc.isAuthenticated || !msc.isRegistered {
		t.Errorf("expected authenticated and registered true, got %v %v", msc.isAuthenticated, msc.isRegistered)
	}
	if msc.currentVin != "1HGCM82633A004352" {
		t.Errorf("expected currentVin 1HGCM82633A004352, got %v", msc.currentVin)
	}
}

// TestDetectModelFromCode tests the model detection functionality
func TestDetectModelFromCode(t *testing.T) {
	tests := []struct {
		modelCode     string
		expectedModel string
		expectedTrim  string
	}{
		{"PDL", "Outback", "Limited XT"},
		{"PFL", "Outback", "Limited"},
		{"PCL", "Outback", "Convenience"},
		{"PBL", "Outback", "Base"},
		{"PDH", "Outback", "Limited XT Hybrid"},
		{"PFH", "Outback", "Limited Hybrid"},
		{"PCH", "Outback", "Convenience Hybrid"},
		{"PBH", "Outback", "Base Hybrid"},
		{"SDL", "Forester", "Limited"},
		{"SFL", "Forester", "Premier"},
		{"SCL", "Forester", "Convenience"},
		{"SBL", "Forester", "Base"},
		{"SDH", "Forester", "Limited Hybrid"},
		{"SFH", "Forester", "Premier Hybrid"},
		{"SCH", "Forester", "Convenience Hybrid"},
		{"SBH", "Forester", "Base Hybrid"},
		{"LDL", "Legacy", "Limited"},
		{"LFL", "Legacy", "Limited XT"},
		{"LCL", "Legacy", "Convenience"},
		{"LBL", "Legacy", "Base"},
		{"CDL", "Crosstrek", "Limited"},
		{"CFL", "Crosstrek", "Premier"},
		{"CCL", "Crosstrek", "Convenience"},
		{"CBL", "Crosstrek", "Base"},
		{"CDH", "Crosstrek", "Limited Hybrid"},
		{"CFH", "Crosstrek", "Premier Hybrid"},
		{"CCH", "Crosstrek", "Convenience Hybrid"},
		{"CBH", "Crosstrek", "Base Hybrid"},
		{"ADL", "Ascent", "Limited"},
		{"AFL", "Ascent", "Premier"},
		{"ACL", "Ascent", "Convenience"},
		{"ABL", "Ascent", "Base"},
		{"WDL", "WRX", "Limited"},
		{"WFL", "WRX", "STI"},
		{"WCL", "WRX", "Convenience"},
		{"WBL", "WRX", "Base"},
		{"IDL", "Impreza", "Limited"},
		{"IFL", "Impreza", "Premier"},
		{"ICL", "Impreza", "Convenience"},
		{"IBL", "Impreza", "Base"},
		{"BDL", "BRZ", "Limited"},
		{"BFL", "BRZ", "STI"},
		{"BCL", "BRZ", "Convenience"},
		{"BBL", "BRZ", "Base"},
		{"TDL", "Solterra", "Limited"},
		{"TFL", "Solterra", "Premier"},
		{"TCL", "Solterra", "Convenience"},
		{"TBL", "Solterra", "Base"},
		{"KRH", "Crosstrek", "Convenience"}, // Older model code for 2017 Crosstrek
		{"XXX", "Unknown", "Unknown"},       // Invalid code
	}

	for _, test := range tests {
		model, trim := DetectModelFromCode(test.modelCode)
		if model != test.expectedModel || trim != test.expectedTrim {
			t.Errorf("DetectModelFromCode(%s) = (%s, %s), expected (%s, %s)", test.modelCode, model, trim, test.expectedModel, test.expectedTrim)
		}
	}
}

func TestGetVehicleStatus_SkipsInvalidFuelPercent(t *testing.T) {
	// Use routes with invalid fuel percent (101%) in responses
	routes := []endpointRoute{
		{Method: http.MethodPost, Path: apiURLs["API_LOGIN"], Response: testLoginResponse},
		{Method: http.MethodGet, Path: apiURLs["API_VALIDATE_SESSION"], Response: testValidateSessionResponse},
		{Method: http.MethodGet, Path: apiURLs["API_SELECT_VEHICLE"], Response: testSelectVehicleResponse},
		{Method: http.MethodGet, Path: apiURLs["API_VEHICLE_HEALTH"], Response: testVehicleHealthResponse},
		{Method: http.MethodGet, Path: apiURLs["API_VEHICLE_STATUS"], Response: testVehicleStatusInvalidFuelResponse},
		{Method: http.MethodGet, Path: apiURLs["API_CONDITION"], Response: testConditionInvalidFuelResponse},
		{Method: http.MethodGet, Path: apiURLs["API_G2_FETCH_RES_SUBARU_PRESETS"], Response: testClimatePresetsSubaruResponse},
		{Method: http.MethodGet, Path: apiURLs["API_G2_FETCH_RES_USER_PRESETS"], Response: testClimatePresetsUserResponse},
		{Method: http.MethodGet, Path: apiURLs["API_G2_FETCH_RES_QUICK_START_SETTINGS"], Response: testClimateQuickStartResponse},
	}

	ts := mockServerWithRoutes(t, routes)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	vehicle, err := msc.GetVehicleByVin("1HGCM82633A004352")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vehicle == nil {
		t.Fatalf("expected vehicle, got nil")
	}
	if vehicle.DistanceToEmpty.Percentage != 0 {
		t.Fatalf("expected fuel percent to remain unset when invalid, got %d", vehicle.DistanceToEmpty.Percentage)
	}
}

// TestVehicleMarshalJSON_EVField verifies that the MarshalJSON method includes
// the computed EV field for backwards compatibility with clients.
func TestVehicleMarshalJSON_EVField(t *testing.T) {
	tests := []struct {
		name     string
		features []string
		expectEV bool
	}{
		{
			name:     "PHEV vehicle should have EV:true",
			features: []string{"PHEV", "g2", "RES"},
			expectEV: true,
		},
		{
			name:     "non-EV vehicle should have EV:false",
			features: []string{"g2", "RES", "ACCS"},
			expectEV: false,
		},
		{
			name:     "empty features should have EV:false",
			features: []string{},
			expectEV: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vehicle{
				Vin:      "TEST123456789",
				Features: tt.features,
			}

			data, err := v.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}

			// Parse the JSON to check the EV field
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			ev, ok := result["EV"]
			if !ok {
				t.Fatal("EV field not found in JSON output")
			}

			evBool, ok := ev.(bool)
			if !ok {
				t.Fatalf("EV field is not a boolean, got %T", ev)
			}

			if evBool != tt.expectEV {
				t.Errorf("expected EV=%v, got EV=%v", tt.expectEV, evBool)
			}
		})
	}
}
