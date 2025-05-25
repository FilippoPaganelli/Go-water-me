/*
  arduino.ino

  This program uses the ArduinoBLE library to set-up an Arduino Nano 33 BLE 
  as a peripheral device and specifies a service and a characteristic. ???

  The circuit:
  - Arduino Nano 33 BLE.
  - ???
*/

#include <ArduinoBLE.h>

const char* deviceServiceUuid = "e969c779-776f-4979-8eb4-d6250e8ea79b";
const char* deviceServiceCharacteristicUuid = "4f6b5586-709d-4b06-94fd-8cbea7c32c28";
BLEService soilMoistureService(deviceServiceUuid);
BLEByteCharacteristic soilMoistureCharacteristic(deviceServiceCharacteristicUuid, BLERead | BLENotify);
byte currentValue = 0;

void setup() {
  Serial.begin(9600);
  while (!Serial)
    ;

  if (!BLE.begin()) {
    Serial.println("ERROR: Starting BLE module failed!");
    while (1)
      ;
  }

  BLE.setLocalName("Go-water-me (peripheral)");
  BLE.setAdvertisedService(soilMoistureService);
  soilMoistureService.addCharacteristic(soilMoistureCharacteristic);
  BLE.addService(soilMoistureService);
  soilMoistureCharacteristic.writeValue(123);
  BLE.advertise();

  Serial.println("INFO: Starting \"Go-water-me (peripheral)\"\n");
}

void loop() {
  BLEDevice central = BLE.central();
  Serial.println("INFO: Discovering central device...");
  delay(500);

  if (central) {
    Serial.println("INFO: Connected to central device!");
    Serial.print("INFO: Device MAC address: ");
    Serial.println(central.address());
    Serial.println("---");

    while (central.connected()) {
      if (soilMoistureCharacteristic.subscribed()) {
        byte newValue = random(0, 256);
        Serial.print("INFO: Sending new value: ");
        Serial.println(newValue);

        soilMoistureCharacteristic.writeValue(newValue);

        delay(2000);
      }
    }

    Serial.println("INFO: Disconnected from central device!");
  }
}
