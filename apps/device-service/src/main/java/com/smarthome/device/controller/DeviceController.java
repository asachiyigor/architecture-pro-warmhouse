package com.smarthome.device.controller;

import com.smarthome.device.dto.DeviceDto;
import com.smarthome.device.model.DeviceStatus;
import com.smarthome.device.model.DeviceType;
import com.smarthome.device.service.DeviceService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;
import java.util.Optional;

@RestController
@RequestMapping("/api/v1/devices")
@Tag(name = "Device Management", description = "API for managing smart home devices")
public class DeviceController {

    @Autowired
    private DeviceService deviceService;

    @GetMapping
    @Operation(summary = "Get all devices", description = "Retrieve all registered devices")
    @ApiResponse(responseCode = "200", description = "Successfully retrieved devices")
    public ResponseEntity<List<DeviceDto>> getAllDevices(
            @RequestParam(required = false) @Parameter(description = "Filter by user ID") Long userId,
            @RequestParam(required = false) @Parameter(description = "Filter by device type") DeviceType type,
            @RequestParam(required = false) @Parameter(description = "Filter by location") String location) {
        
        List<DeviceDto> devices;
        
        if (userId != null) {
            devices = deviceService.getDevicesByUserId(userId);
        } else if (type != null) {
            devices = deviceService.getDevicesByType(type);
        } else if (location != null) {
            devices = deviceService.getDevicesByLocation(location);
        } else {
            devices = deviceService.getAllDevices();
        }
        
        return ResponseEntity.ok(devices);
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get device by ID", description = "Retrieve a specific device by its ID")
    @ApiResponse(responseCode = "200", description = "Device found")
    @ApiResponse(responseCode = "404", description = "Device not found")
    public ResponseEntity<DeviceDto> getDeviceById(
            @PathVariable @Parameter(description = "Device ID", required = true) Long id) {
        
        Optional<DeviceDto> device = deviceService.getDeviceById(id);
        return device.map(ResponseEntity::ok)
                    .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    @Operation(summary = "Create new device", description = "Register a new device in the system")
    @ApiResponse(responseCode = "201", description = "Device created successfully")
    @ApiResponse(responseCode = "400", description = "Invalid device data")
    public ResponseEntity<DeviceDto> createDevice(
            @Valid @RequestBody @Parameter(description = "Device data", required = true) DeviceDto deviceDto) {
        
        DeviceDto createdDevice = deviceService.createDevice(deviceDto);
        return ResponseEntity.status(HttpStatus.CREATED).body(createdDevice);
    }

    @PutMapping("/{id}")
    @Operation(summary = "Update device", description = "Update an existing device")
    @ApiResponse(responseCode = "200", description = "Device updated successfully")
    @ApiResponse(responseCode = "404", description = "Device not found")
    @ApiResponse(responseCode = "400", description = "Invalid device data")
    public ResponseEntity<DeviceDto> updateDevice(
            @PathVariable @Parameter(description = "Device ID", required = true) Long id,
            @Valid @RequestBody @Parameter(description = "Updated device data", required = true) DeviceDto deviceDto) {
        
        Optional<DeviceDto> updatedDevice = deviceService.updateDevice(id, deviceDto);
        return updatedDevice.map(ResponseEntity::ok)
                           .orElse(ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{id}")
    @Operation(summary = "Delete device", description = "Remove a device from the system")
    @ApiResponse(responseCode = "204", description = "Device deleted successfully")
    @ApiResponse(responseCode = "404", description = "Device not found")
    public ResponseEntity<Void> deleteDevice(
            @PathVariable @Parameter(description = "Device ID", required = true) Long id) {
        
        boolean deleted = deviceService.deleteDevice(id);
        return deleted ? ResponseEntity.noContent().build() 
                      : ResponseEntity.notFound().build();
    }

    @PatchMapping("/{id}/status")
    @Operation(summary = "Update device status", description = "Update the operational status of a device")
    @ApiResponse(responseCode = "200", description = "Device status updated successfully")
    @ApiResponse(responseCode = "404", description = "Device not found")
    @ApiResponse(responseCode = "400", description = "Invalid status")
    public ResponseEntity<DeviceDto> updateDeviceStatus(
            @PathVariable @Parameter(description = "Device ID", required = true) Long id,
            @RequestBody @Parameter(description = "Status update", required = true) Map<String, DeviceStatus> statusUpdate) {
        
        DeviceStatus status = statusUpdate.get("status");
        if (status == null) {
            return ResponseEntity.badRequest().build();
        }
        
        Optional<DeviceDto> updatedDevice = deviceService.updateDeviceStatus(id, status);
        return updatedDevice.map(ResponseEntity::ok)
                           .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping("/{id}/heartbeat")
    @Operation(summary = "Device heartbeat", description = "Register device heartbeat to indicate it's online")
    @ApiResponse(responseCode = "200", description = "Heartbeat registered successfully")
    @ApiResponse(responseCode = "404", description = "Device not found")
    public ResponseEntity<Map<String, String>> deviceHeartbeat(
            @PathVariable @Parameter(description = "Device ID", required = true) Long id,
            @RequestBody(required = false) @Parameter(description = "Heartbeat data") Map<String, String> heartbeatData) {
        
        String ipAddress = heartbeatData != null ? heartbeatData.get("ipAddress") : null;
        deviceService.heartbeat(id, ipAddress);
        
        return ResponseEntity.ok(Map.of("message", "Heartbeat registered successfully"));
    }

    @GetMapping("/stats")
    @Operation(summary = "Get device statistics", description = "Retrieve device statistics and counts")
    @ApiResponse(responseCode = "200", description = "Statistics retrieved successfully")
    public ResponseEntity<DeviceService.DeviceStats> getDeviceStats() {
        DeviceService.DeviceStats stats = deviceService.getDeviceStats();
        return ResponseEntity.ok(stats);
    }

    @GetMapping("/types")
    @Operation(summary = "Get device types", description = "Retrieve all available device types")
    @ApiResponse(responseCode = "200", description = "Device types retrieved successfully")
    public ResponseEntity<DeviceType[]> getDeviceTypes() {
        return ResponseEntity.ok(DeviceType.values());
    }

    @GetMapping("/statuses")
    @Operation(summary = "Get device statuses", description = "Retrieve all available device statuses")
    @ApiResponse(responseCode = "200", description = "Device statuses retrieved successfully")
    public ResponseEntity<DeviceStatus[]> getDeviceStatuses() {
        return ResponseEntity.ok(DeviceStatus.values());
    }
}