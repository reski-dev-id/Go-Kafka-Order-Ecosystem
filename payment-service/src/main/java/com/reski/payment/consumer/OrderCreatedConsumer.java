package com.reski.payment.consumer;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.reski.payment.dto.OrderCreatedEvent;
import com.reski.payment.dto.PaymentCompletedEvent;
import com.reski.payment.entity.Payment;
import com.reski.payment.entity.ProcessedEvent;
import com.reski.payment.repository.PaymentRepository;
import com.reski.payment.repository.ProcessedEventRepository;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.common.header.Header;

import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.KafkaHeaders;

import org.springframework.messaging.support.MessageBuilder;

import org.springframework.stereotype.Component;

import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.UUID;

@Slf4j
@Component
@RequiredArgsConstructor
public class OrderCreatedConsumer {

    private final PaymentRepository paymentRepository;

    private final ProcessedEventRepository processedEventRepository;

    private final ObjectMapper objectMapper;

    private final KafkaTemplate<String, String> kafkaTemplate;

    @Transactional
    @KafkaListener(
            topics = "order.created",
            groupId = "payment-group"
    )
    public void consume(
            ConsumerRecord<String, String> record
    ) {

        String message = record.value();

        try {

            log.info(
                    "received kafka message: {}",
                    message
            );

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

                log.warn(
                        "duplicate event skipped: {}",
                        event.getEventId()
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

            log.info(
                    "payment completed published for order: {}",
                    event.getId()
            );

        } catch (Exception e) {

            log.error(
                    "payment processing failed: {}",
                    e.getMessage()
            );

            Header retryHeader =
                    record.headers().lastHeader(
                            "retry-count"
                    );

            int retryCount = 0;

            if (retryHeader != null) {

                retryCount = Integer.parseInt(
                        new String(retryHeader.value())
                );
            }

            retryCount++;

            log.warn(
                    "retry count: {}",
                    retryCount
            );

            if (retryCount < 3) {

                kafkaTemplate.send(
                        MessageBuilder
                                .withPayload(message)
                                .setHeader(
                                        KafkaHeaders.TOPIC,
                                        "order.created"
                                )
                                .setHeader(
                                        "retry-count",
                                        String.valueOf(retryCount)
                                )
                                .build()
                );

                log.warn(
                        "message retried to kafka topic"
                );

                return;
            }

            kafkaTemplate.send(
                    "order.created.dlq",
                    message
            );

            log.error(
                    "message sent to dlq after max retries"
            );
        }
    }
}