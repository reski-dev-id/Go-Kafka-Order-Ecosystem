package com.reski.payment.consumer;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.reski.payment.dto.OrderCreatedEvent;
import com.reski.payment.dto.PaymentCompletedEvent;
import com.reski.payment.entity.Payment;
import com.reski.payment.entity.ProcessedEvent;
import com.reski.payment.repository.PaymentRepository;
import com.reski.payment.repository.ProcessedEventRepository;

import lombok.RequiredArgsConstructor;

import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

import java.time.Instant;
import java.util.UUID;

@Component
@RequiredArgsConstructor
public class OrderCreatedConsumer {

    private final PaymentRepository paymentRepository;

    private final ProcessedEventRepository processedEventRepository;

    private final ObjectMapper objectMapper;

    private final KafkaTemplate<String, String> kafkaTemplate;

    @KafkaListener(
            topics = "order.created",
            groupId = "payment-group"
    )
    public void consume(String message) throws Exception {

        OrderCreatedEvent event =
                objectMapper.readValue(
                        message,
                        OrderCreatedEvent.class
                );

        boolean alreadyProcessed =
                processedEventRepository
                        .existsByEventIdAndConsumerGroup(
                                event.getEventId(),
                                "payment-service"
                        );

        if (alreadyProcessed) {

            System.out.println(
                    "duplicate event skipped: "
                            + event.getEventId()
            );

            return;
        }

        Payment payment = Payment.builder()
                .id(UUID.randomUUID())
                .orderId(UUID.fromString(event.getId()))
                .amount(event.getAmount())
                .status("PAID")
                .createdAt(Instant.now())
                .build();

        paymentRepository.save(payment);

        ProcessedEvent processedEvent =
                new ProcessedEvent();

        processedEvent.setEventId(
                event.getEventId()
        );

        processedEvent.setConsumerGroup(
                "payment-service"
        );

        processedEventRepository.save(
                processedEvent
        );

        PaymentCompletedEvent completedEvent =
                PaymentCompletedEvent.builder()
                        .orderId(event.getId())
                        .amount(event.getAmount())
                        .status("PAID")
                        .build();

        kafkaTemplate.send(
                "payment.completed",
                event.getId(),
                objectMapper.writeValueAsString(completedEvent)
        );

        System.out.println(
                "payment completed published for order: "
                        + event.getId()
        );
    }
}