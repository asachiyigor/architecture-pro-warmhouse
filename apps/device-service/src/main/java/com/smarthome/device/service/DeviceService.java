package com.smarthome.device.service;

import com.smarthome.device.dto.DeviceDto;
import com.smarthome.device.model.Device;
import com.smarthome.device.model.DeviceStatus;
import com.smarthome.device.model.DeviceType;
import com.smarthome.device.repository.DeviceRepository;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

@Service
public class DeviceService {
    
    @Autowired
    private DeviceRepository deviceRepository;
    
    @Autowired
    private KafkaTemplate<String, Object> kafkaTemplate;
    
    private static final String DEVICE_EVENTS_TOPIC = "device.events";

    /**
     * Get all devices
     */
    public List<DeviceDto> getAllDevices() {
        return deviceRepository.findAll().stream()
                .map(this::convertToDto)
                .collect(Collectors.toList());
    }

    /**
     * Get device by ID
     */
    public Optional<DeviceDto> getDeviceById(Long id) {
        return deviceRepository.findById(id)
                .map(this::convertToDto);
    }

    /**
     * Get devices by user ID
     */
    public List<DeviceDto> getDevicesByUserId(Long userId) {
        return deviceRepository.findByUserId(userId).stream()
                .map(this::convertToDto)
                .collect(Collectors.toList());
    }

    /**
     * Get devices by type
     */
    public List<DeviceDto> getDevicesByType(DeviceType type) {
        return deviceRepository.findByType(type).stream()
                .map(this::convertToDto)
                .collect(Collectors.toList());
    }

    /**
     * Get devices by location
     */
    public List<DeviceDto> getDevicesByLocation(String location) {
        return deviceRepository.findByLocationContainingIgnoreCase(location).stream()
                .map(this::convertToDto)
                .collect(Collectors.toList());
    }

    /**
     * Create new device
     */
    public DeviceDto createDevice(DeviceDto deviceDto) {
        Device device = convertToEntity(deviceDto);
        device.setCreatedAt(LocalDateTime.now());
        device.setStatus(DeviceStatus.OFFLINE); // New devices start offline
        
        Device savedDevice = deviceRepository.save(device);
        
        // Publish device created event
        publishDeviceEvent("device.created", savedDevice);
        
        return convertToDto(savedDevice);
    }

    /**
     * Update device
     */
    public Optional<DeviceDto> updateDevice(Long id, DeviceDto deviceDto) {
        Optional<Device> existingDevice = deviceRepository.findById(id);
        
        if (existingDevice.isPresent()) {
            Device device = existingDevice.get();
            updateDeviceFields(device, deviceDto);
            device.setUpdatedAt(LocalDateTime.now());
            
            Device savedDevice = deviceRepository.save(device);
            
            // Publish device updated event
            publishDeviceEvent("device.updated", savedDevice);
            
            return Optional.of(convertToDto(savedDevice));
        }
        
        return Optional.empty();
    }

    /**
     * Delete device
     */
    public boolean deleteDevice(Long id) {
        Optional<Device> device = deviceRepository.findById(id);
        
        if (device.isPresent()) {
            deviceRepository.deleteById(id);
            
            // Publish device deleted event
            publishDeviceEvent("device.deleted", device.get());
            
            return true;
        }
        
        return false;
    }

    /**
     * Update device status
     */
    public Optional<DeviceDto> updateDeviceStatus(Long id, DeviceStatus status) {
        Optional<Device> existingDevice = deviceRepository.findById(id);
        
        if (existingDevice.isPresent()) {
            Device device = existingDevice.get();
            DeviceStatus oldStatus = device.getStatus();
            
            device.setStatus(status);
            device.setLastSeen(LocalDateTime.now());
            device.setUpdatedAt(LocalDateTime.now());
            
            Device savedDevice = deviceRepository.save(device);
            
            // Publish status change event if status actually changed
            if (!oldStatus.equals(status)) {
                publishDeviceStatusChangeEvent(savedDevice, oldStatus, status);
            }
            
            return Optional.of(convertToDto(savedDevice));
        }
        
        return Optional.empty();
    }

    /**
     * Register device heartbeat
     */
    public void heartbeat(Long deviceId, String ipAddress) {
        Optional<Device> existingDevice = deviceRepository.findById(deviceId);
        
        if (existingDevice.isPresent()) {
            Device device = existingDevice.get();
            device.setLastSeen(LocalDateTime.now());
            device.setIpAddress(ipAddress);
            
            // If device was offline, mark as online
            if (device.getStatus() == DeviceStatus.OFFLINE) {
                device.setStatus(DeviceStatus.ONLINE);
                publishDeviceStatusChangeEvent(device, DeviceStatus.OFFLINE, DeviceStatus.ONLINE);
            }
            
            deviceRepository.save(device);
        }
    }

    /**
     * Get device statistics
     */
    public DeviceStats getDeviceStats() {
        long totalDevices = deviceRepository.count();
        long onlineDevices = deviceRepository.countByStatus(DeviceStatus.ONLINE);
        long offlineDevices = deviceRepository.countByStatus(DeviceStatus.OFFLINE);
        long errorDevices = deviceRepository.countByStatus(DeviceStatus.ERROR);
        
        return new DeviceStats(totalDevices, onlineDevices, offlineDevices, errorDevices);
    }

    // Helper methods

    private DeviceDto convertToDto(Device device) {
        DeviceDto dto = new DeviceDto();
        dto.setId(device.getId());
        dto.setName(device.getName());
        dto.setType(device.getType());
        dto.setLocation(device.getLocation());
        dto.setStatus(device.getStatus());
        dto.setMacAddress(device.getMacAddress());
        dto.setIpAddress(device.getIpAddress());
        dto.setFirmwareVersion(device.getFirmwareVersion());
        dto.setUserId(device.getUserId());
        dto.setCreatedAt(device.getCreatedAt());
        dto.setUpdatedAt(device.getUpdatedAt());
        dto.setLastSeen(device.getLastSeen());
        return dto;
    }

    private Device convertToEntity(DeviceDto dto) {
        Device device = new Device();
        device.setName(dto.getName());
        device.setType(dto.getType());
        device.setLocation(dto.getLocation());
        device.setMacAddress(dto.getMacAddress());
        device.setIpAddress(dto.getIpAddress());
        device.setFirmwareVersion(dto.getFirmwareVersion());
        device.setUserId(dto.getUserId());
        return device;
    }

    private void updateDeviceFields(Device device, DeviceDto dto) {
        if (dto.getName() != null) device.setName(dto.getName());
        if (dto.getType() != null) device.setType(dto.getType());
        if (dto.getLocation() != null) device.setLocation(dto.getLocation());
        if (dto.getMacAddress() != null) device.setMacAddress(dto.getMacAddress());
        if (dto.getIpAddress() != null) device.setIpAddress(dto.getIpAddress());
        if (dto.getFirmwareVersion() != null) device.setFirmwareVersion(dto.getFirmwareVersion());
        if (dto.getUserId() != null) device.setUserId(dto.getUserId());
    }

    private void publishDeviceEvent(String eventType, Device device) {
        DeviceEvent event = new DeviceEvent(eventType, device);
        kafkaTemplate.send(DEVICE_EVENTS_TOPIC, eventType, event);
    }

    private void publishDeviceStatusChangeEvent(Device device, DeviceStatus oldStatus, DeviceStatus newStatus) {
        DeviceStatusChangeEvent event = new DeviceStatusChangeEvent(
            device.getId(), device.getName(), device.getLocation(),
            oldStatus, newStatus, LocalDateTime.now()
        );
        kafkaTemplate.send(DEVICE_EVENTS_TOPIC, "device.status.changed", event);
    }

    // Inner classes for events and stats
    public static class DeviceEvent {
        private String eventType;
        private Long deviceId;
        private String deviceName;
        private DeviceType deviceType;
        private String location;
        private DeviceStatus status;
        private LocalDateTime timestamp;

        public DeviceEvent(String eventType, Device device) {
            this.eventType = eventType;
            this.deviceId = device.getId();
            this.deviceName = device.getName();
            this.deviceType = device.getType();
            this.location = device.getLocation();
            this.status = device.getStatus();
            this.timestamp = LocalDateTime.now();
        }

        // Getters
        public String getEventType() { return eventType; }
        public Long getDeviceId() { return deviceId; }
        public String getDeviceName() { return deviceName; }
        public DeviceType getDeviceType() { return deviceType; }
        public String getLocation() { return location; }
        public DeviceStatus getStatus() { return status; }
        public LocalDateTime getTimestamp() { return timestamp; }
    }

    public static class DeviceStatusChangeEvent {
        private Long deviceId;
        private String deviceName;
        private String location;
        private DeviceStatus oldStatus;
        private DeviceStatus newStatus;
        private LocalDateTime timestamp;

        public DeviceStatusChangeEvent(Long deviceId, String deviceName, String location,
                                     DeviceStatus oldStatus, DeviceStatus newStatus, LocalDateTime timestamp) {
            this.deviceId = deviceId;
            this.deviceName = deviceName;
            this.location = location;
            this.oldStatus = oldStatus;
            this.newStatus = newStatus;
            this.timestamp = timestamp;
        }

        // Getters
        public Long getDeviceId() { return deviceId; }
        public String getDeviceName() { return deviceName; }
        public String getLocation() { return location; }
        public DeviceStatus getOldStatus() { return oldStatus; }
        public DeviceStatus getNewStatus() { return newStatus; }
        public LocalDateTime getTimestamp() { return timestamp; }
    }

    public static class DeviceStats {
        private long totalDevices;
        private long onlineDevices;
        private long offlineDevices;
        private long errorDevices;

        public DeviceStats(long totalDevices, long onlineDevices, long offlineDevices, long errorDevices) {
            this.totalDevices = totalDevices;
            this.onlineDevices = onlineDevices;
            this.offlineDevices = offlineDevices;
            this.errorDevices = errorDevices;
        }

        // Getters
        public long getTotalDevices() { return totalDevices; }
        public long getOnlineDevices() { return onlineDevices; }
        public long getOfflineDevices() { return offlineDevices; }
        public long getErrorDevices() { return errorDevices; }
    }
}