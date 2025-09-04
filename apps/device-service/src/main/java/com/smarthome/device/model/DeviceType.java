package com.smarthome.device.model;

/**
 * Enum for different types of smart home devices
 */
public enum DeviceType {
    TEMPERATURE_SENSOR("Temperature Sensor"),
    HUMIDITY_SENSOR("Humidity Sensor"),
    MOTION_SENSOR("Motion Sensor"),
    DOOR_SENSOR("Door/Window Sensor"),
    SMOKE_DETECTOR("Smoke Detector"),
    SMART_LIGHT("Smart Light"),
    SMART_PLUG("Smart Plug"),
    THERMOSTAT("Thermostat"),
    SECURITY_CAMERA("Security Camera"),
    SMART_LOCK("Smart Lock"),
    AIR_CONDITIONER("Air Conditioner"),
    HEATER("Heater"),
    UNKNOWN("Unknown Device");

    private final String displayName;

    DeviceType(String displayName) {
        this.displayName = displayName;
    }

    public String getDisplayName() {
        return displayName;
    }
}