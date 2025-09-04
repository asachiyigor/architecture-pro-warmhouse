package com.smarthome.device.repository;

import com.smarthome.device.model.Device;
import com.smarthome.device.model.DeviceStatus;
import com.smarthome.device.model.DeviceType;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Repository
public interface DeviceRepository extends JpaRepository<Device, Long> {
    
    /**
     * Find devices by user ID
     */
    List<Device> findByUserId(Long userId);
    
    /**
     * Find devices by location
     */
    List<Device> findByLocationContainingIgnoreCase(String location);
    
    /**
     * Find devices by type
     */
    List<Device> findByType(DeviceType type);
    
    /**
     * Find devices by status
     */
    List<Device> findByStatus(DeviceStatus status);
    
    /**
     * Find device by MAC address
     */
    Optional<Device> findByMacAddress(String macAddress);
    
    /**
     * Find devices by user and type
     */
    List<Device> findByUserIdAndType(Long userId, DeviceType type);
    
    /**
     * Find devices that haven't been seen since given time
     */
    @Query("SELECT d FROM Device d WHERE d.lastSeen < :cutoffTime OR d.lastSeen IS NULL")
    List<Device> findStaleDevices(@Param("cutoffTime") LocalDateTime cutoffTime);
    
    /**
     * Count devices by status
     */
    long countByStatus(DeviceStatus status);
    
    /**
     * Count devices by user
     */
    long countByUserId(Long userId);
}