package com.reski.payment.repository;

import com.reski.payment.entity.ProcessedEvent;
import org.springframework.data.jpa.repository.JpaRepository;

public interface ProcessedEventRepository
        extends JpaRepository<ProcessedEvent, Long> {

    boolean existsByEventIdAndConsumerGroup(
            String eventId,
            String consumerGroup
    );
}