package com.smarthome.device.model;

/**
 * Enum for device operational status
 */
public enum DeviceStatus {
    ONLINE("Online"),
    OFFLINE("Offline"),
    ERROR("Error"),
    MAINTENANCE("Maintenance"),
    CONFIGURING("Configuring"),
    LOW_BATTERY("Low Battery");

    private final String displayName;

    DeviceStatus(String displayName) {
        this.displayName = displayName;
    }

    public String getDisplayName() {
        return displayName;
    }
}