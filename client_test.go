package mysubaru

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex-savin/go-mysubaru/config"
)

func mockConfig(t *testing.T) *config.Config {
	return &config.Config{
		MySubaru: config.MySubaru{
			Credentials: config.Credentials{
				Username:   "user",
				Password:   "pass",
				DeviceID:   "dev123",
				DeviceName: "devname",
			},
			Region:        "TEST",
			AutoReconnect: true,
		},
		TimeZone: "America/New_York",
		// Logger:   slogt.New(t),
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

// mockMySubaruApi creates a mock MySubaru API server for testing.
func mockMySubaruApi(t *testing.T, handler http.HandlerFunc) *httptest.Server {

	// Create a listener with the desired port
	l, err := net.Listen("tcp", "127.0.0.1:56765")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}

	ts := httptest.NewUnstartedServer(handler)
	// http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// Handle API_LOGIN endpoint
	// 	if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_LOGIN"] && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"success":true,"errorCode":"BIOMETRICS_DISABLED","dataName":"sessionData","data":{"sessionChanged":false,"vehicleInactivated":false,"account":{"marketId":1,"createdDate":1476984644000,"firstName":"Tatiana","lastName":"Savin","zipCode":"07974","accountKey":765268,"lastLoginDate":1751738613000,"zipCode5":"07974"},"resetPassword":false,"deviceId":"JddMBQXvAkgutSmEP6uFsThbq4QgEBBQ","sessionId":"9D7FCDF274794346689D3FA0D693CBBF","deviceRegistered":true,"passwordToken":null,"vehicles":[{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":null,"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":null,"authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":null,"subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":null,"sunsetUpgraded":true,"intDescrip":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}],"rightToRepairEnabled":true,"rightToRepairStartYear":2022,"rightToRepairStates":"MA","enableXtime":true,"termsAndConditionsAccepted":true,"digitalGlobeConnectId":"0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeImageTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/DigitalGlobe:ImageryTileService@EPSG:3857@png/{z}/{x}/{y}.png?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeTransparentTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/Digitalglobe:OSMTransparentTMSTileService@EPSG:3857@png/{z}/{x}/{-y}.png/?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","tomtomKey":"DHH9SwEQ4MW55Hj2TfqMeldbsDjTdgAs","currentVehicleIndex":0,"handoffToken":"$2a$08$rOb/uqhm8I3QtSel2phOCOxNM51w43eqXDDksMkJ.1a5KsaQuLvEu$1751745334477","satelliteViewEnabled":true,"registeredDevicePermanent":true}}`)
	// 	}
	// 	// Handle API_VALIDATE_SESSION endpoint
	// 	if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_VALIDATE_SESSION"] && r.Method == http.MethodGet {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	//		fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":null,"data":null}`)
	// 	}
	// 	// Handle API_LOCATE endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LOCATE"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LOCATE"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"id": 1, "name": "John Doe"}`)
	// 	}
	// 	// Handle API_LOCK endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LOCK"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LOCK"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"remoteServiceStatus","data":{"serviceRequestId":"1HGCM82633A004352_1751747301812_20_@NGTP","success":false,"cancelled":false,"remoteServiceType":"lock","remoteServiceState":"started","subState":null,"errorCode":null,"result":null,"updateTime":null,"vin":"1HGCM82633A004352","errorDescription":null}}`)
	// 	}
	// 	// Handle API_LOCK_CANCEL endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LOCK_CANCEL"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LOCK_CANCEL"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"id": 1, "name": "John Doe"}`)
	// 	}
	// 	// Handle API_UNLOCK endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_UNLOCK"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_UNLOCK"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"remoteServiceStatus","data":{"serviceRequestId":"1HGCM82633A004352_1751747271262_19_@NGTP","success":false,"cancelled":false,"remoteServiceType":"unlock","remoteServiceState":"started","subState":null,"errorCode":null,"result":null,"updateTime":null,"vin":"1HGCM82633A004352","errorDescription":null}}`)
	// 	}
	// 	// Handle API_UNLOCK_CANCEL endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_UNLOCK_CANCEL"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_UNLOCK_CANCEL"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"id": 1, "name": "John Doe"}`)
	// 	}
	// 	// Handle API_HORN_LIGHTS endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_HORN_LIGHTS"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_HORN_LIGHTS"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"remoteServiceStatus","data":{"serviceRequestId":"1HGCM82633A004352_1751746568969_47_@NGTP","success":true,"cancelled":false,"remoteServiceType":"hornLights","remoteServiceState":"started","subState":null,"errorCode":null,"result":null,"updateTime":1751746569000,"vin":"1HGCM82633A004352","errorDescription":null}}`)
	// 	}
	// 	// Handle API_HORN_LIGHTS_CANCEL endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_HORN_LIGHTS_CANCEL"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_HORN_LIGHTS_CANCEL"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"remoteServiceStatus","data":{"serviceRequestId":"1HGCM82633A004352_1751746568969_47_@NGTP","success":true,"cancelled":false,"remoteServiceType":"hornLights","remoteServiceState":"cancelling","subState":null,"errorCode":null,"result":null,"updateTime":1751746569000,"vin":"1HGCM82633A004352","errorDescription":null}}`)
	// 	}
	// 	// Handle API_HORN_LIGHTS_STOP endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_HORN_LIGHTS_STOP"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_HORN_LIGHTS_STOP"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"remoteServiceStatus","data":{"serviceRequestId":"1HGCM82633A004352_1751746568969_47_@NGTP","success":true,"cancelled":false,"remoteServiceType":"hornLights","remoteServiceState":"stopping","subState":null,"errorCode":null,"result":null,"updateTime":1751746569000,"vin":"1HGCM82633A004352","errorDescription":null}}`)
	// 	}
	// 	// Handle API_LIGHTS endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LIGHTS"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LIGHTS"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"remoteServiceStatus","data":{"serviceRequestId":"1HGCM82633A004352_1751746568969_47_@NGTP","success":true,"cancelled":false,"remoteServiceType":"lightsOnly","remoteServiceState":"started","subState":null,"errorCode":null,"result":null,"updateTime":1751746569000,"vin":"1HGCM82633A004352","errorDescription":null}}`)
	// 	}
	// 	// Handle API_LIGHTS_CANCEL endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LIGHTS_CANCEL"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LIGHTS_CANCEL"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"remoteServiceStatus","data":{"serviceRequestId":"1HGCM82633A004352_1751746568969_47_@NGTP","success":true,"cancelled":false,"remoteServiceType":"lightsOnly","remoteServiceState":"stopping","subState":null,"errorCode":null,"result":null,"updateTime":1751746569000,"vin":"1HGCM82633A004352","errorDescription":null}}`)
	// 	}
	// 	// Handle API_LIGHTS_STOP endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LIGHTS_STOP"], "g1") || r.URL.Path == MOBILE_API_VERSION+urlToGen(apiURLs["API_LIGHTS_STOP"], "g2")) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"remoteServiceStatus","data":{"serviceRequestId":"1HGCM82633A004352_1751746568969_47_@NGTP","success":true,"cancelled":false,"remoteServiceType":"lightsOnly","remoteServiceState":"finished","subState":null,"errorCode":null,"result":null,"updateTime":1751746569000,"vin":"1HGCM82633A004352","errorDescription":null}}`)
	// 	}
	// 	// Handle API_G2_REMOTE_ENGINE_START endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+apiURLs["API_G2_REMOTE_ENGINE_START"]) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"id": 1, "name": "John Doe"}`)
	// 	}
	// 	// Handle API_G2_REMOTE_ENGINE_START_CANCEL endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+apiURLs["API_G2_REMOTE_ENGINE_START_CANCEL"]) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		fmt.Fprint(w, `{"id": 1, "name": "John Doe"}`)
	// 	}
	// 	// Handle API_G2_REMOTE_ENGINE_STOP endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+apiURLs["API_G2_REMOTE_ENGINE_STOP"]) && r.Method == http.MethodPost {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		json.NewEncoder(w).Encode(
	// 			Response{
	// 				Success:  true,
	// 				DataName: "remoteServiceStatus",
	// 				Data:     []byte(`{"serviceRequestId":"1HGCM82633A004352_1751745367457_47_@NGTP","success":true,"cancelled":false,"remoteServiceType":"lightsOnly","remoteServiceState":"started","subState":null,"errorCode":null,"result":null,"updateTime":1751745367000,"vin":"1HGCM82633A004352","errorDescription":null}`),
	// 			},
	// 		)
	// 	}
	// 	// Handle API_REMOTE_SVC_STATUS endpoint
	// 	if (r.URL.Path == MOBILE_API_VERSION+apiURLs["API_REMOTE_SVC_STATUS"]) && r.Method == http.MethodGet {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		json.NewEncoder(w).Encode(
	// 			Response{
	// 				Success:  true,
	// 				DataName: "remoteServiceStatus",
	// 				Data:     []byte(`{"remoteServiceState":"finished"}`),
	// 			},
	// 		)
	// 	}

	// Close the default listener created by NewUnstartedServer and replace it
	// with our custom listener.
	ts.Listener.Close()
	ts.Listener = l

	return ts
}

func TestNew_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Handle API_LOGIN endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_LOGIN"] && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":"BIOMETRICS_DISABLED","dataName":"sessionData","data":{"sessionChanged":false,"vehicleInactivated":false,"account":{"marketId":1,"createdDate":1476984644000,"firstName":"Tatiana","lastName":"Savin","zipCode":"07974","accountKey":765268,"lastLoginDate":1751738613000,"zipCode5":"07974"},"resetPassword":false,"deviceId":"JddMBQXvAkgutSmEP6uFsThbq4QgEBBQ","sessionId":"9D7FCDF274794346689D3FA0D693CBBF","deviceRegistered":true,"passwordToken":null,"vehicles":[{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"Subaru Outback TXT","stolenVehicle":false,"vin":"1HGCM82633A004352","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"Subaru Outback TXT","vehicleKey":8211380,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":-1,"zip":null,"oemCustId":"CRM-41PLM-5TYE","vehicleMileage":null,"phone":null,"timeZone":"America/New_York","features":null,"userOemCustId":"CRM-41PLM-5TYE","subscriptionStatus":null,"authorizedVehicle":false,"preferredDealer":null,"cachedStateCode":"NJ","modelName":null,"subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"extDescrip":null,"sunsetUpgraded":true,"intDescrip":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false}],"rightToRepairEnabled":true,"rightToRepairStartYear":2022,"rightToRepairStates":"MA","enableXtime":true,"termsAndConditionsAccepted":true,"digitalGlobeConnectId":"0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeImageTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/DigitalGlobe:ImageryTileService@EPSG:3857@png/{z}/{x}/{y}.png?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","digitalGlobeTransparentTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/Digitalglobe:OSMTransparentTMSTileService@EPSG:3857@png/{z}/{x}/{-y}.png/?connectId=0572e32b-2fcf-4bc8-abe0-1e3da8767132","tomtomKey":"DHH9SwEQ4MW55Hj2TfqMeldbsDjTdgAs","currentVehicleIndex":0,"handoffToken":"$2a$08$rOb/uqhm8I3QtSel2phOCOxNM51w43eqXDDksMkJ.1a5KsaQuLvEu$1751745334477","satelliteViewEnabled":true,"registeredDevicePermanent":true}}`)
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

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	if !msc.isAuthenticated || !msc.isRegistered {
		t.Errorf("expected authenticated and registered true, got %v %v", msc.isAuthenticated, msc.isRegistered)
	}
	if msc.currentVin != "1HGCM82633A004352" {
		t.Errorf("expected currentVin 1HGCM82633A004352, got %v", msc.currentVin)
	}
}

func TestNew_MultiCarSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Handle API_LOGIN endpoint
		if r.URL.Path == MOBILE_API_VERSION+apiURLs["API_LOGIN"] && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success":true,"errorCode":null,"dataName":"sessionData","data":{"sessionChanged":false,"vehicleInactivated":false,"account":{"createdDate":1612345678000,"marketId":1,"firstName":"Joe","lastName":"User","zipCode":"54321","accountKey":1234567,"lastLoginDate":1612345678000,"zipCode5":"54321"},"resetPassword":false,"deviceId":"1612345678","sessionId":"0123456789ABCDEF01234567890ABCDE","deviceRegistered":true,"passwordToken":null,"vehicles":[{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"TEST_SUBARU_1","stolenVehicle":false,"features":null,"vin":"JF2ABCDE6L0000001","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"TEST_SUBARU_1","vehicleKey":1000001,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":1,"zip":null,"oemCustId":"1-TESTOEM_1","vehicleMileage":null,"phone":null,"userOemCustId":"1-TESTOEM_1","subscriptionStatus":null,"authorizedVehicle":true,"preferredDealer":null,"cachedStateCode":"TX","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"sunsetUpgraded":true,"extDescrip":null,"intDescrip":null,"modelName":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false,"timeZone":null},{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"TEST_SUBARU_2","stolenVehicle":false,"features":null,"vin":"JF2ABCDE6L0000002","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"TEST_SUBARU_2","vehicleKey":1000002,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":1,"zip":null,"oemCustId":"1-TESTOEM_2","vehicleMileage":null,"phone":null,"userOemCustId":"1-TESTOEM_2","subscriptionStatus":null,"authorizedVehicle":true,"preferredDealer":null,"cachedStateCode":"TX","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"sunsetUpgraded":true,"extDescrip":null,"intDescrip":null,"modelName":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false,"timeZone":null},{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"TEST_SUBARU_3","stolenVehicle":false,"features":null,"vin":"JF2ABCDE6L0000003","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"TEST_SUBARU_3","vehicleKey":1000003,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":1,"zip":null,"oemCustId":"1-TESTOEM_3","vehicleMileage":null,"phone":null,"userOemCustId":"1-TESTOEM_3","subscriptionStatus":null,"authorizedVehicle":true,"preferredDealer":null,"cachedStateCode":"TX","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"sunsetUpgraded":true,"extDescrip":null,"intDescrip":null,"modelName":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false,"timeZone":null},{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"TEST_SUBARU_4","stolenVehicle":false,"features":null,"vin":"JF2ABCDE6L0000004","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"TEST_SUBARU_4","vehicleKey":1000004,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":1,"zip":null,"oemCustId":"1-TESTOEM_4","vehicleMileage":null,"phone":null,"userOemCustId":"1-TESTOEM_4","subscriptionStatus":null,"authorizedVehicle":true,"preferredDealer":null,"cachedStateCode":"TX","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"sunsetUpgraded":true,"extDescrip":null,"intDescrip":null,"modelName":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false,"timeZone":null},{"customer":{"sessionCustomer":null,"email":null,"firstName":null,"lastName":null,"zip":null,"oemCustId":null,"phone":null},"vehicleName":"TEST_SUBARU_5","stolenVehicle":false,"features":null,"vin":"JF2ABCDE6L0000005","modelYear":null,"modelCode":null,"engineSize":null,"nickname":"TEST_SUBARU_5","vehicleKey":1000005,"active":true,"licensePlate":"","licensePlateState":"","email":null,"firstName":null,"lastName":null,"subscriptionFeatures":null,"accessLevel":1,"zip":null,"oemCustId":"1-TESTOEM_5","vehicleMileage":null,"phone":null,"userOemCustId":"1-TESTOEM_5","subscriptionStatus":null,"authorizedVehicle":true,"preferredDealer":null,"cachedStateCode":"TX","subscriptionPlans":[],"crmRightToRepair":false,"needMileagePrompt":false,"phev":null,"sunsetUpgraded":true,"extDescrip":null,"intDescrip":null,"modelName":null,"transCode":null,"provisioned":true,"remoteServicePinExist":true,"needEmergencyContactPrompt":false,"vehicleGeoPosition":null,"show3gSunsetBanner":false,"timeZone":null}],"rightToRepairEnabled":true,"rightToRepairStartYear":2022,"rightToRepairStates":"MA","enableXtime":true,"termsAndConditionsAccepted":true,"digitalGlobeConnectId":"00000000-0000-0000-0000-000000000000","digitalGlobeImageTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/DigitalGlobe:ImageryTileService@EPSG:3857@png/{z}/{x}/{y}.png?connectId=00000000-0000-0000-0000-000000000000","digitalGlobeTransparentTileService":"https://earthwatch.digitalglobe.com/earthservice/tmsaccess/tms/1.0.0/Digitalglobe:OSMTransparentTMSTileService@EPSG:3857@png/{z}/{x}/{-y}.png/?connectId=00000000-0000-0000-0000-000000000000","tomtomKey":"0123456789ABCDEF01234567890ABCDE","currentVehicleIndex":0,"handoffToken":"test","satelliteViewEnabled":true,"registeredDevicePermanent":true}}`)
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

	// Authenticate the client
	ok, authErr, _ := msc.Authenticate()
	if !ok || authErr != nil {
		t.Fatalf("expected authentication to succeed, got ok=%v, err=%v", ok, authErr)
	}

	if !msc.isAuthenticated || !msc.isRegistered {
		t.Errorf("expected authenticated and registered true, got %v %v", msc.isAuthenticated, msc.isRegistered)
	}
	if msc.currentVin != "JF2ABCDE6L0000001" {
		t.Errorf("expected currentVin JF2ABCDE6L0000001, got %v", msc.currentVin)
	}
}

// func TestNew_Failure(t *testing.T) {
// 	ts := mockMySubaruApi(t)
// 	ts.Start()
// 	defer ts.Close() // Ensure the server is closed after the test

// 	cfg := makeConfig(t)
// 	cfg.MySubaru.Credentials.Username = "" // Invalid username

// 	client, err := New(cfg)
// 	if err == nil {
// 		t.Fatalf("expected error, got nil")
// 	}
// 	if client != nil {
// 		t.Fatalf("expected nil client, got %v", client)
// 	}
// 	if client.isAuthenticated || client.isRegistered {
// 		t.Errorf("expected authenticated and registered false, got %v %v", client.isAuthenticated, client.isRegistered)
// 	}
// 	if client.currentVin != "" {
// 		t.Errorf("expected currentVin empty, got %v", client.currentVin)
// 	}
// }

func TestSelectVehicle_Success(t *testing.T) {
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
	}
	ts := mockMySubaruApi(t, handler)
	ts.Start()
	defer ts.Close()

	cfg := mockConfig(t)

	msc, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	vehicle, err := msc.SelectVehicle("1HGCM82633A004352")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vehicle == nil {
		t.Fatalf("expected vehicle, got nil")
	}
	if vehicle.Vin != "1HGCM82633A004352" {
		t.Errorf("expected vehicle VIN 1HGCM82633A004352, got %v", vehicle.Vin)
	}
}

func TestSelectVehicle_InvalidVIN(t *testing.T) {
	cfg := mockConfig(t)
	c, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error creating client, got %v", err)
	}

	if _, err := c.SelectVehicle("INVALIDVIN"); err == nil {
		t.Fatalf("expected VIN validation error, got nil")
	}
}

func TestGetVehicleByVin_Success(t *testing.T) {
	routes := standardTestRoutes()

	ts := mockServerWithRoutes(t, routes)
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

	vehicle, err := msc.GetVehicleByVin("1HGCM82633A004352")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vehicle == nil {
		t.Fatalf("expected vehicle, got nil")
	}
	if vehicle.Vin != "1HGCM82633A004352" {
		t.Errorf("expected vehicle VIN 1HGCM82633A004352, got %v", vehicle.Vin)
	}
}
